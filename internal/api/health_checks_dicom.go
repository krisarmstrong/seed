package api

// health_checks_dicom.go contains DICOM C-ECHO (A-ASSOCIATE) reachability
// tests for the health-check pipeline (Issue #777).

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

// DICOM protocol constants.
const (
	dicomDefaultPort    = 104
	dicomTestTimeoutSec = 10
	// DICOM PDU types.
	dicomPDUAssocRQ = 0x01 // A-ASSOCIATE-RQ.
	dicomPDUAssocAC = 0x02 // A-ASSOCIATE-AC.
	dicomPDUAssocRJ = 0x03 // A-ASSOCIATE-RJ.
	// DICOM item types.
	dicomItemAppContext      = 0x10 // Application Context Item.
	dicomItemPresentationCtx = 0x20 // Presentation Context Item.
	dicomItemAbstractSyntax  = 0x30 // Abstract Syntax Sub-Item.
	dicomItemTransferSyntax  = 0x40 // Transfer Syntax Sub-Item.
	dicomItemUserInfo        = 0x50 // User Information Item.
	dicomItemMaxPDULength    = 0x51 // Maximum Length Sub-Item.
	// DICOM byte values.
	dicomReservedByte          = 0x00 // Reserved byte value.
	dicomProtocolVersionMSB    = 0x00 // Protocol version (1) MSB.
	dicomProtocolVersionLSB    = 0x01 // Protocol version (1) LSB.
	dicomPresentationContextID = 0x01 // Presentation context ID.
	dicomMaxPDULengthByte      = 0x04 // Length of max PDU sub-item (4 bytes).
	dicomMaxPDUValue           = 0x40 // High byte of 16384 (0x00004000).
	// DICOM size constants.
	dicomAETitleLength       = 16 // AE Title length.
	dicomReservedBlockLength = 32 // Reserved block length in A-ASSOCIATE-RQ.
	dicomMaxPDUItemSize      = 8  // Max PDU Length Sub-Item size.
	dicomPDUHeaderSize       = 6  // PDU header size (type + reserved + 4-byte length).
	// Bit shift constants.
	dicomShiftByteHigh    = 24 // Shift for high byte in 32-bit value.
	dicomShiftByteMidHigh = 16 // Shift for mid-high byte.
	dicomShiftByteMidLow  = 8  // Shift for mid-low byte.
)

// runDICOMTests runs all configured DICOM server tests and returns results.
func (s *Server) runDICOMTests(ctx context.Context) []CustomTestResult {
	results := make([]CustomTestResult, 0, len(s.config.HealthChecks.DICOMEndpoints))
	threshold := s.config.Thresholds.CustomTests.TCP // Use TCP thresholds for DICOM

	for _, endpoint := range s.config.HealthChecks.DICOMEndpoints {
		if !endpoint.Enabled {
			continue
		}

		name := endpoint.Name
		if name == "" {
			name = net.JoinHostPort(endpoint.Host, strconv.Itoa(endpoint.Port))
		}

		testResult := CustomTestResult{Name: name, Host: endpoint.Host, Port: endpoint.Port}
		latency, err := runDICOMTest(ctx, endpoint.Host, endpoint.Port, endpoint.CalledAE, endpoint.CallingAE)

		if err != nil {
			testResult.Success = false
			testResult.Error = "DICOM test failed: " + err.Error()
			testResult.TestStatus = statusError
		} else {
			testResult.Success = true
			testResult.Latency = latency
			warningMs := threshold.Warning.Milliseconds()
			criticalMs := threshold.Critical.Milliseconds()
			testResult.TestStatus = getTestStatus(latency, warningMs, criticalMs)
		}
		results = append(results, testResult)
	}
	return results
}

// runDICOMTest runs a DICOM C-ECHO (association request) and returns latency in ms.
// This tests the DICOM server's ability to accept associations (like a ping).
func runDICOMTest(ctx context.Context, host string, port int, calledAE, callingAE string) (float64, error) {
	if port == 0 {
		port = dicomDefaultPort
	}
	if calledAE == "" {
		calledAE = "ANY-SCP"
	}
	if callingAE == "" {
		callingAE = "SEED-SCU"
	}

	addr := net.JoinHostPort(host, strconv.Itoa(port))

	ctx, cancel := context.WithTimeout(ctx, dicomTestTimeoutSec*time.Second)
	defer cancel()

	start := time.Now()

	// Connect via TCP
	dialer := net.Dialer{Timeout: dicomTestTimeoutSec * time.Second}
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return 0, fmt.Errorf("connection failed: %w", err)
	}
	defer func() { _ = conn.Close() }()

	// Set deadline for the entire exchange
	if deadlineErr := conn.SetDeadline(time.Now().Add(dicomTestTimeoutSec * time.Second)); deadlineErr != nil {
		return 0, deadlineErr
	}

	// Build A-ASSOCIATE-RQ PDU
	pdu := buildDICOMAssociateRQ(calledAE, callingAE)

	// Send A-ASSOCIATE-RQ
	_, err = conn.Write(pdu)
	if err != nil {
		return 0, fmt.Errorf("failed to send A-ASSOCIATE-RQ: %w", err)
	}

	// Read response PDU header
	header := make([]byte, dicomPDUHeaderSize)
	_, err = conn.Read(header)
	if err != nil {
		return 0, fmt.Errorf("failed to read response: %w", err)
	}

	latency := time.Since(start).Seconds() * millisecondsPerSecond

	// Check PDU type
	pduType := header[0]
	switch pduType {
	case dicomPDUAssocAC:
		// Association accepted - server is healthy
		return latency, nil
	case dicomPDUAssocRJ:
		// Association rejected - but server responded, so it's reachable
		// This still means the DICOM server is running
		return latency, nil
	default:
		return latency, fmt.Errorf("unexpected PDU type: 0x%02x", pduType)
	}
}

// buildDICOMAssociateRQ builds a minimal DICOM A-ASSOCIATE-RQ PDU.
// This is a simplified version that requests Verification SOP Class (C-ECHO).
func buildDICOMAssociateRQ(calledAE, callingAE string) []byte {
	// Pad AE titles to 16 chars
	calledAE = padAETitle(calledAE)
	callingAE = padAETitle(callingAE)

	// Verification SOP Class UID (1.2.840.10008.1.1)
	verificationUID := "1.2.840.10008.1.1"
	// Implicit VR Little Endian Transfer Syntax UID (1.2.840.10008.1.2)
	implicitVRUID := "1.2.840.10008.1.2"

	// Build Application Context Item
	appContextUID := "1.2.840.10008.3.1.1.1"
	appContextItem := buildDICOMItem(dicomItemAppContext, []byte(appContextUID))

	// Build Abstract Syntax Sub-Item
	abstractSyntaxItem := buildDICOMItem(dicomItemAbstractSyntax, []byte(verificationUID))

	// Build Transfer Syntax Sub-Item
	transferSyntaxItem := buildDICOMItem(dicomItemTransferSyntax, []byte(implicitVRUID))

	// Build Presentation Context Item
	pcContent := make([]byte, 0)
	pcContent = append(pcContent, dicomPresentationContextID) // Presentation Context ID
	pcContent = append(pcContent, dicomReservedByte)          // Reserved
	pcContent = append(pcContent, dicomReservedByte)          // Reserved
	pcContent = append(pcContent, dicomReservedByte)          // Reserved
	pcContent = append(pcContent, abstractSyntaxItem...)
	pcContent = append(pcContent, transferSyntaxItem...)
	presentationContextItem := buildDICOMItem(dicomItemPresentationCtx, pcContent)

	// Build User Information Item with Max PDU Length Sub-Item
	maxPDULength := make([]byte, dicomMaxPDUItemSize)
	maxPDULength[0] = dicomItemMaxPDULength // Item type
	maxPDULength[1] = dicomReservedByte     // Reserved
	maxPDULength[2] = dicomReservedByte     // Length MSB
	maxPDULength[3] = dicomMaxPDULengthByte // Length LSB (4 bytes)
	maxPDULength[4] = dicomReservedByte     // Max PDU length (16384 = 0x00004000)
	maxPDULength[5] = dicomReservedByte     // ...
	maxPDULength[6] = dicomMaxPDUValue      // ...
	maxPDULength[7] = dicomReservedByte     // ...
	userInfoItem := buildDICOMItem(dicomItemUserInfo, maxPDULength)

	// Build PDU content
	pduContent := make([]byte, 0)
	pduContent = append(pduContent, dicomProtocolVersionMSB, dicomProtocolVersionLSB) // Protocol Version
	pduContent = append(pduContent, dicomReservedByte, dicomReservedByte)             // Reserved
	pduContent = append(pduContent, []byte(calledAE)...)
	pduContent = append(pduContent, []byte(callingAE)...)
	pduContent = append(pduContent, make([]byte, dicomReservedBlockLength)...) // Reserved (32 bytes)
	pduContent = append(pduContent, appContextItem...)
	pduContent = append(pduContent, presentationContextItem...)
	pduContent = append(pduContent, userInfoItem...)

	// Build full PDU
	pdu := make([]byte, 0)
	pdu = append(pdu, dicomPDUAssocRQ)   // PDU Type
	pdu = append(pdu, dicomReservedByte) // Reserved
	// PDU Length (4 bytes, big endian)
	pduLen := len(pduContent)
	pdu = append(pdu, byte(pduLen>>dicomShiftByteHigh), byte(pduLen>>dicomShiftByteMidHigh),
		byte(pduLen>>dicomShiftByteMidLow), byte(pduLen))
	pdu = append(pdu, pduContent...)

	return pdu
}

// buildDICOMItem builds a DICOM item with type, reserved, length, and content.
func buildDICOMItem(itemType byte, content []byte) []byte {
	item := make([]byte, 0)
	item = append(item, itemType)
	item = append(item, dicomReservedByte) // Reserved
	// Length (2 bytes, big endian)
	length := len(content)
	item = append(item, byte(length>>dicomShiftByteMidLow), byte(length))
	item = append(item, content...)
	return item
}

// padAETitle pads an AE title to dicomAETitleLength characters with spaces.
func padAETitle(ae string) string {
	if len(ae) > dicomAETitleLength {
		return ae[:dicomAETitleLength]
	}
	return ae + strings.Repeat(" ", dicomAETitleLength-len(ae))
}
