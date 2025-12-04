package iperf

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// ClientConfig holds iperf3 client test configuration
type ClientConfig struct {
	Server   string `json:"server"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"` // "tcp" or "udp"
	Reverse  bool   `json:"reverse"`  // true = download (server sends), false = upload (client sends)
	Duration int    `json:"duration"` // seconds
	Parallel int    `json:"parallel"` // number of streams
}

// Result contains the iperf3 test results
type Result struct {
	BitsPerSecond   float64   `json:"bitsPerSecond"`
	Bandwidth       float64   `json:"bandwidth"`       // Mbps
	Transfer        float64   `json:"transfer"`        // MB
	Retransmits     int       `json:"retransmits"`     // TCP only
	Jitter          float64   `json:"jitter"`          // UDP only, ms
	LostPackets     int       `json:"lostPackets"`     // UDP only
	LostPercent     float64   `json:"lostPercent"`     // UDP only
	Protocol        string    `json:"protocol"`
	Direction       string    `json:"direction"`       // "upload" or "download"
	Duration        float64   `json:"duration"`        // seconds
	Server          string    `json:"server"`
	Port            int       `json:"port"`
	Timestamp       time.Time `json:"timestamp"`
}

// ServerStatus represents the iperf3 server status
type ServerStatus struct {
	Running bool   `json:"running"`
	Port    int    `json:"port"`
	PID     int    `json:"pid"`
	Error   string `json:"error,omitempty"`
}

// ClientStatus represents the client test status
type ClientStatus struct {
	Running  bool    `json:"running"`
	Phase    string  `json:"phase"` // "idle", "connecting", "testing", "complete"
	Progress float64 `json:"progress"`
}

// iperfJSON is the structure of iperf3 -J output
type iperfJSON struct {
	Start struct {
		Connected []struct {
			Socket     int    `json:"socket"`
			LocalHost  string `json:"local_host"`
			LocalPort  int    `json:"local_port"`
			RemoteHost string `json:"remote_host"`
			RemotePort int    `json:"remote_port"`
		} `json:"connected"`
		TestStart struct {
			Protocol   string `json:"protocol"`
			NumStreams int    `json:"num_streams"`
			Duration   int    `json:"duration"`
			Reverse    int    `json:"reverse"`
		} `json:"test_start"`
	} `json:"start"`
	End struct {
		SumSent struct {
			Seconds       float64 `json:"seconds"`
			Bytes         float64 `json:"bytes"`
			BitsPerSecond float64 `json:"bits_per_second"`
			Retransmits   int     `json:"retransmits"`
		} `json:"sum_sent"`
		SumReceived struct {
			Seconds       float64 `json:"seconds"`
			Bytes         float64 `json:"bytes"`
			BitsPerSecond float64 `json:"bits_per_second"`
		} `json:"sum_received"`
		Sum struct {
			Seconds       float64 `json:"seconds"`
			Bytes         float64 `json:"bytes"`
			BitsPerSecond float64 `json:"bits_per_second"`
			JitterMs      float64 `json:"jitter_ms"`
			LostPackets   int     `json:"lost_packets"`
			Packets       int     `json:"packets"`
			LostPercent   float64 `json:"lost_percent"`
		} `json:"sum"`
	} `json:"end"`
	Error string `json:"error"`
}

// Manager handles iperf3 client and server operations
type Manager struct {
	mu           sync.RWMutex
	serverStatus ServerStatus
	clientStatus ClientStatus
	lastResult   *Result
	serverCmd    *exec.Cmd
	serverCancel context.CancelFunc
}

// NewManager creates a new iperf3 manager
func NewManager() *Manager {
	return &Manager{
		clientStatus: ClientStatus{Phase: "idle"},
	}
}

// CheckInstalled checks if iperf3 is available
func CheckInstalled() error {
	_, err := exec.LookPath("iperf3")
	if err != nil {
		return fmt.Errorf("iperf3 not found: %w", err)
	}
	return nil
}

// GetVersion returns the installed iperf3 version
func GetVersion() (string, error) {
	if err := CheckInstalled(); err != nil {
		return "", err
	}
	out, err := exec.Command("iperf3", "--version").Output()
	if err != nil {
		return "", fmt.Errorf("failed to get iperf3 version: %w", err)
	}
	// Output is like "iperf 3.16 (cJSON 1.7.17)\n..."
	lines := strings.Split(string(out), "\n")
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0]), nil
	}
	return "unknown", nil
}

// GetServerStatus returns the current server status
func (m *Manager) GetServerStatus() ServerStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.serverStatus
}

// GetClientStatus returns the current client status
func (m *Manager) GetClientStatus() ClientStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.clientStatus
}

// GetLastResult returns the last test result
func (m *Manager) GetLastResult() *Result {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.lastResult
}

// StartServer starts the iperf3 server
func (m *Manager) StartServer(port int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.serverStatus.Running {
		return fmt.Errorf("server already running on port %d", m.serverStatus.Port)
	}

	if err := CheckInstalled(); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	m.serverCancel = cancel

	// Start iperf3 server: iperf3 -s -p <port> -D (daemon mode doesn't work well, use background)
	cmd := exec.CommandContext(ctx, "iperf3", "-s", "-p", fmt.Sprintf("%d", port))
	if err := cmd.Start(); err != nil {
		cancel()
		return fmt.Errorf("failed to start iperf3 server: %w", err)
	}

	m.serverCmd = cmd
	m.serverStatus = ServerStatus{
		Running: true,
		Port:    port,
		PID:     cmd.Process.Pid,
	}

	// Monitor the server process
	go func() {
		err := cmd.Wait()
		m.mu.Lock()
		m.serverStatus.Running = false
		if err != nil && ctx.Err() == nil {
			m.serverStatus.Error = err.Error()
		}
		m.mu.Unlock()
	}()

	return nil
}

// StopServer stops the iperf3 server
func (m *Manager) StopServer() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.serverStatus.Running {
		return fmt.Errorf("server not running")
	}

	if m.serverCancel != nil {
		m.serverCancel()
	}

	if m.serverCmd != nil && m.serverCmd.Process != nil {
		m.serverCmd.Process.Kill()
	}

	m.serverStatus = ServerStatus{Running: false}
	m.serverCmd = nil
	m.serverCancel = nil

	return nil
}

// RunClient runs an iperf3 client test
func (m *Manager) RunClient(ctx context.Context, config ClientConfig) (*Result, error) {
	m.mu.Lock()
	if m.clientStatus.Running {
		m.mu.Unlock()
		return nil, fmt.Errorf("test already in progress")
	}
	m.clientStatus = ClientStatus{Running: true, Phase: "connecting", Progress: 10}
	m.mu.Unlock()

	defer func() {
		m.mu.Lock()
		m.clientStatus = ClientStatus{Running: false, Phase: "idle", Progress: 0}
		m.mu.Unlock()
	}()

	if err := CheckInstalled(); err != nil {
		return nil, err
	}

	// Set defaults
	if config.Port == 0 {
		config.Port = 5201
	}
	if config.Duration == 0 {
		config.Duration = 10
	}
	if config.Parallel == 0 {
		config.Parallel = 1
	}
	if config.Protocol == "" {
		config.Protocol = "tcp"
	}

	// Build command
	args := []string{
		"-c", config.Server,
		"-p", fmt.Sprintf("%d", config.Port),
		"-t", fmt.Sprintf("%d", config.Duration),
		"-P", fmt.Sprintf("%d", config.Parallel),
		"-J", // JSON output
	}

	if config.Protocol == "udp" {
		args = append(args, "-u")
		args = append(args, "-b", "0") // Unlimited bandwidth for UDP
	}

	if config.Reverse {
		args = append(args, "-R") // Reverse mode (server sends, client receives)
	}

	m.mu.Lock()
	m.clientStatus.Phase = "testing"
	m.clientStatus.Progress = 30
	m.mu.Unlock()

	// Run iperf3
	cmd := exec.CommandContext(ctx, "iperf3", args...)
	output, err := cmd.Output()
	if err != nil {
		// Try to parse error from output
		var iperfOut iperfJSON
		if jsonErr := json.Unmarshal(output, &iperfOut); jsonErr == nil && iperfOut.Error != "" {
			return nil, fmt.Errorf("iperf3 error: %s", iperfOut.Error)
		}
		return nil, fmt.Errorf("iperf3 failed: %w", err)
	}

	m.mu.Lock()
	m.clientStatus.Progress = 80
	m.mu.Unlock()

	// Parse JSON output
	var iperfOut iperfJSON
	if err := json.Unmarshal(output, &iperfOut); err != nil {
		return nil, fmt.Errorf("failed to parse iperf3 output: %w", err)
	}

	if iperfOut.Error != "" {
		return nil, fmt.Errorf("iperf3 error: %s", iperfOut.Error)
	}

	// Build result
	result := &Result{
		Protocol:  config.Protocol,
		Server:    config.Server,
		Port:      config.Port,
		Timestamp: time.Now(),
	}

	if config.Reverse {
		result.Direction = "download"
		// In reverse mode, we care about what we received
		result.BitsPerSecond = iperfOut.End.SumReceived.BitsPerSecond
		result.Bandwidth = iperfOut.End.SumReceived.BitsPerSecond / 1_000_000
		result.Transfer = iperfOut.End.SumReceived.Bytes / 1_000_000
		result.Duration = iperfOut.End.SumReceived.Seconds
	} else {
		result.Direction = "upload"
		// In normal mode, we care about what we sent
		result.BitsPerSecond = iperfOut.End.SumSent.BitsPerSecond
		result.Bandwidth = iperfOut.End.SumSent.BitsPerSecond / 1_000_000
		result.Transfer = iperfOut.End.SumSent.Bytes / 1_000_000
		result.Duration = iperfOut.End.SumSent.Seconds
		result.Retransmits = iperfOut.End.SumSent.Retransmits
	}

	// UDP-specific fields
	if config.Protocol == "udp" {
		result.Jitter = iperfOut.End.Sum.JitterMs
		result.LostPackets = iperfOut.End.Sum.LostPackets
		result.LostPercent = iperfOut.End.Sum.LostPercent
	}

	m.mu.Lock()
	m.lastResult = result
	m.clientStatus.Phase = "complete"
	m.clientStatus.Progress = 100
	m.mu.Unlock()

	return result, nil
}
