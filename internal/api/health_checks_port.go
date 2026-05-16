package api

// health_checks_port.go contains TCP, UDP, and DNS port-reachability tests
// for health checks.

import (
	"context"
	"errors"
	"net"
	"strconv"
	"time"
)

// runTCPTests runs all configured TCP port tests and returns results.
func (s *Server) runTCPTests(ctx context.Context) []CustomTestResult {
	results := make([]CustomTestResult, 0, len(s.config.HealthChecks.TCPPorts))
	threshold := s.config.Thresholds.CustomTests.TCP

	for _, target := range s.config.HealthChecks.TCPPorts {
		if !target.Enabled {
			continue
		}

		name := target.Name
		if name == "" {
			name = net.JoinHostPort(target.Host, strconv.Itoa(target.Port))
		}

		testResult := CustomTestResult{Name: name, Host: target.Host, Port: target.Port}
		latency, err := runTCPTest(ctx, target.Host, target.Port)

		if err != nil {
			testResult.Success = false
			testResult.Error = "TCP connection test failed"
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

// runTCPTest runs a TCP port test and returns latency in ms.
func runTCPTest(ctx context.Context, host string, port int) (float64, error) {
	ctx, cancel := context.WithTimeout(ctx, tcpTestTimeoutSec*time.Second)
	defer cancel()

	addr := net.JoinHostPort(host, strconv.Itoa(port))
	start := time.Now()
	conn, err := (&net.Dialer{}).DialContext(ctx, "tcp", addr)
	if err != nil {
		return 0, err
	}
	latency := time.Since(start).Seconds() * millisecondsPerSecond
	_ = conn.Close()
	return latency, nil
}

// runUDPTests runs all configured UDP port tests and returns results.
func (s *Server) runUDPTests() []CustomTestResult {
	results := make([]CustomTestResult, 0, len(s.config.HealthChecks.UDPPorts))
	threshold := s.config.Thresholds.CustomTests.UDP

	for _, target := range s.config.HealthChecks.UDPPorts {
		if !target.Enabled {
			continue
		}

		name := target.Name
		if name == "" {
			name = net.JoinHostPort(target.Host, strconv.Itoa(target.Port))
		}

		testResult := CustomTestResult{Name: name, Host: target.Host, Port: target.Port}
		latency, err := runUDPTest(target.Host, target.Port)

		if err != nil {
			testResult.Success = false
			testResult.Error = "UDP connection test failed"
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

// runUDPTest runs a UDP port test and returns latency in ms.
// Note: UDP is connectionless, so we send a packet and wait for ICMP unreachable
// or application response. For DNS (53), NTP (123), etc. we can get actual responses.
func runUDPTest(host string, port int) (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), udpTestTimeoutSec*time.Second)
	defer cancel()

	addr := net.JoinHostPort(host, strconv.Itoa(port))

	// For DNS port, try a simple DNS query
	if port == dnsPort {
		return testDNSPort(ctx, host)
	}

	// For other UDP ports, we try to connect (which on UDP just sets up local state)
	// and send a small probe packet
	start := time.Now()

	dialer := net.Dialer{Timeout: udpTestTimeoutSec * time.Second}
	conn, err := dialer.DialContext(ctx, "udp", addr)
	if err != nil {
		return 0, err
	}
	defer func() { _ = conn.Close() }()

	// Set deadline for response
	if deadlineErr := conn.SetDeadline(time.Now().Add(udpReadDeadlineSec * time.Second)); deadlineErr != nil {
		return 0, deadlineErr
	}

	// Send a small probe packet
	_, err = conn.Write([]byte{0x00})
	if err != nil {
		return 0, err
	}

	// Try to read response (may timeout for non-responding services)
	buf := make([]byte, udpReadBufferBytes)
	_, err = conn.Read(buf)

	latency := time.Since(start).Seconds() * millisecondsPerSecond

	// For UDP, no error on Write means the port is likely open
	// (no ICMP unreachable received)
	if err != nil {
		// Check if it's a timeout (which for UDP often means the port is open but not responding)
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			// Port is likely open but service didn't respond - still count as success
			return latency, nil
		}
		// Connection refused or other error means port is closed
		return 0, errors.New("port closed or filtered")
	}

	return latency, nil
}

// testDNSPort tests DNS port by sending a simple query.
func testDNSPort(ctx context.Context, host string) (float64, error) {
	// Use Go's resolver to test DNS
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(dialCtx context.Context, _, _ string) (net.Conn, error) {
			d := net.Dialer{Timeout: udpTestTimeoutSec * time.Second}
			return d.DialContext(dialCtx, "udp", host+":53")
		},
	}

	start := time.Now()
	_, err := resolver.LookupHost(ctx, "google.com")
	latency := time.Since(start).Seconds() * millisecondsPerSecond

	if err != nil {
		return 0, err
	}
	return latency, nil
}
