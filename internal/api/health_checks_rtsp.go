package api

// health_checks_rtsp.go contains RTSP OPTIONS reachability tests for the
// health-check pipeline (Issue #778).

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

// RTSP protocol constants.
const (
	rtspDefaultPort     = 554
	rtspTestTimeoutSec  = 10
	rtspReadBufferBytes = 1024
)

// runRTSPTests runs all configured RTSP stream tests and returns results.
func (s *Server) runRTSPTests(ctx context.Context) []CustomTestResult {
	results := make([]CustomTestResult, 0, len(s.config.HealthChecks.RTSPEndpoints))
	threshold := s.config.Thresholds.CustomTests.TCP // Use TCP thresholds for RTSP

	for _, endpoint := range s.config.HealthChecks.RTSPEndpoints {
		if !endpoint.Enabled {
			continue
		}

		name := endpoint.Name
		if name == "" {
			name = endpoint.URL
		}

		testResult := CustomTestResult{Name: name, URL: endpoint.URL}
		latency, err := runRTSPTest(ctx, endpoint.URL)

		if err != nil {
			testResult.Success = false
			testResult.Error = "RTSP test failed: " + err.Error()
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

// runRTSPTest runs an RTSP OPTIONS request and returns latency in ms.
// RTSP uses a simple text-based protocol similar to HTTP.
func runRTSPTest(ctx context.Context, rtspURL string) (float64, error) {
	// Parse RTSP URL to extract host and port
	// rtsp://host:port/path
	url := strings.TrimPrefix(rtspURL, "rtsp://")
	url = strings.TrimPrefix(url, "rtsps://")

	// Split host:port from path
	hostPort, _, _ := strings.Cut(url, "/")

	// Add default port if not specified
	if !strings.Contains(hostPort, ":") {
		hostPort = hostPort + ":" + strconv.Itoa(rtspDefaultPort)
	}

	ctx, cancel := context.WithTimeout(ctx, rtspTestTimeoutSec*time.Second)
	defer cancel()

	start := time.Now()

	// Connect via TCP
	dialer := net.Dialer{Timeout: rtspTestTimeoutSec * time.Second}
	conn, err := dialer.DialContext(ctx, "tcp", hostPort)
	if err != nil {
		return 0, fmt.Errorf("connection failed: %w", err)
	}
	defer func() { _ = conn.Close() }()

	// Set deadline for the entire exchange
	if deadlineErr := conn.SetDeadline(time.Now().Add(rtspTestTimeoutSec * time.Second)); deadlineErr != nil {
		return 0, deadlineErr
	}

	// Send RTSP OPTIONS request
	request := fmt.Sprintf("OPTIONS %s RTSP/1.0\r\nCSeq: 1\r\n\r\n", rtspURL)
	_, err = conn.Write([]byte(request))
	if err != nil {
		return 0, fmt.Errorf("failed to send OPTIONS: %w", err)
	}

	// Read response
	buf := make([]byte, rtspReadBufferBytes)
	n, err := conn.Read(buf)
	if err != nil {
		return 0, fmt.Errorf("failed to read response: %w", err)
	}

	latency := time.Since(start).Seconds() * millisecondsPerSecond

	// Check for valid RTSP response
	response := string(buf[:n])
	if !strings.Contains(response, "RTSP/1.0") {
		return latency, errors.New("invalid RTSP response")
	}

	// Check for 200 OK status
	if !strings.Contains(response, "200 OK") && !strings.Contains(response, "200") {
		return latency, errors.New("RTSP server returned error")
	}

	return latency, nil
}
