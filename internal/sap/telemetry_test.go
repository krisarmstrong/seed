package sap_test

import (
	"context"
	"testing"
	"time"

	"github.com/krisarmstrong/seed/internal/config"
	"github.com/krisarmstrong/seed/internal/sap"
)

// TestNewTelemetryService verifies that NewTelemetryService creates a valid service.
func TestNewTelemetryService(t *testing.T) {
	cfg := config.DefaultConfig()
	service := sap.NewTelemetryService(cfg, nil)

	if service == nil {
		t.Fatal("expected non-nil TelemetryService")
	}
}

// TestTelemetryServiceStartStop verifies Start and Stop don't panic.
func TestTelemetryServiceStartStop(t *testing.T) {
	cfg := config.DefaultConfig()
	service := sap.NewTelemetryService(cfg, nil)

	ctx := context.Background()
	if err := service.Start(ctx); err != nil {
		t.Errorf("Start() returned error: %v", err)
	}

	// Stop should not panic.
	service.Stop()

	// Multiple stops should be safe.
	service.Stop()
}

// TestTelemetryServiceGetSnapshotNotImplemented verifies GetSnapshot returns ErrNotImplemented.
func TestTelemetryServiceGetSnapshotNotImplemented(t *testing.T) {
	cfg := config.DefaultConfig()
	service := sap.NewTelemetryService(cfg, nil)

	ctx := context.Background()
	snapshot, err := service.GetSnapshot(ctx)

	if err != sap.ErrNotImplemented {
		t.Errorf("expected ErrNotImplemented, got %v", err)
	}
	if snapshot != nil {
		t.Errorf("expected nil snapshot, got %+v", snapshot)
	}
}

// TestTelemetryServiceGetHistoryNotImplemented verifies GetHistory returns ErrNotImplemented.
func TestTelemetryServiceGetHistoryNotImplemented(t *testing.T) {
	cfg := config.DefaultConfig()
	service := sap.NewTelemetryService(cfg, nil)

	ctx := context.Background()
	history, err := service.GetHistory(ctx, "", "")

	if err != sap.ErrNotImplemented {
		t.Errorf("expected ErrNotImplemented, got %v", err)
	}
	if history != nil {
		t.Errorf("expected nil history, got %+v", history)
	}
}

// TestTelemetryServiceWithCanceledContext verifies behavior with canceled context.
func TestTelemetryServiceWithCanceledContext(t *testing.T) {
	cfg := config.DefaultConfig()
	service := sap.NewTelemetryService(cfg, nil)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately.

	// Start with canceled context should not panic.
	if err := service.Start(ctx); err != nil {
		t.Errorf("Start() with canceled context returned error: %v", err)
	}

	service.Stop()
}

// TestTelemetrySnapshotFields verifies TelemetrySnapshot struct fields.
func TestTelemetrySnapshotFields(t *testing.T) {
	now := time.Now()

	snapshot := sap.TelemetrySnapshot{
		Timestamp: now,
		Links:     []sap.LinkStatus{},
		Gateway:   nil,
		DNS:       nil,
		DHCP:      nil,
		Bandwidth: nil,
	}

	if snapshot.Timestamp != now {
		t.Errorf("expected Timestamp %v, got %v", now, snapshot.Timestamp)
	}
	if snapshot.Links == nil {
		t.Error("expected non-nil Links slice")
	}
	if len(snapshot.Links) != 0 {
		t.Errorf("expected empty Links slice, got %d elements", len(snapshot.Links))
	}
	if snapshot.Gateway != nil {
		t.Error("expected nil Gateway")
	}
	if snapshot.DNS != nil {
		t.Error("expected nil DNS")
	}
	if snapshot.DHCP != nil {
		t.Error("expected nil DHCP")
	}
	if snapshot.Bandwidth != nil {
		t.Error("expected nil Bandwidth")
	}
}

// TestTelemetrySnapshotWithData verifies TelemetrySnapshot with populated fields.
func TestTelemetrySnapshotWithData(t *testing.T) {
	now := time.Now()

	linkStatus := sap.LinkStatus{
		Interface:  "eth0",
		State:      sap.LinkStateUp,
		Speed:      "1000Mbps",
		Duplex:     "full",
		MTU:        1500,
		MACAddress: "00:11:22:33:44:55",
		IPAddress:  "192.168.1.100",
		Carrier:    true,
		UpdatedAt:  now,
	}

	gatewayHealth := &sap.GatewayHealth{
		IP:         "192.168.1.1",
		Reachable:  true,
		RTT:        time.Millisecond * 5,
		RTTMs:      5.0,
		PacketLoss: 0.0,
		Status:     sap.HealthStatusHealthy,
		LastCheck:  now,
	}

	bandwidth := &sap.BandwidthSample{
		Interface:     "eth0",
		TxBytesPerSec: 1000000,
		RxBytesPerSec: 5000000,
		TxMbps:        8.0,
		RxMbps:        40.0,
		Utilization:   5.0,
		SampledAt:     now,
	}

	snapshot := sap.TelemetrySnapshot{
		Timestamp: now,
		Links:     []sap.LinkStatus{linkStatus},
		Gateway:   gatewayHealth,
		Bandwidth: bandwidth,
	}

	if len(snapshot.Links) != 1 {
		t.Fatalf("expected 1 link, got %d", len(snapshot.Links))
	}
	if snapshot.Links[0].Interface != "eth0" {
		t.Errorf("expected Interface 'eth0', got %q", snapshot.Links[0].Interface)
	}
	if snapshot.Links[0].State != sap.LinkStateUp {
		t.Errorf("expected State LinkStateUp, got %v", snapshot.Links[0].State)
	}
	if snapshot.Gateway == nil {
		t.Fatal("expected non-nil Gateway")
	}
	if snapshot.Gateway.IP != "192.168.1.1" {
		t.Errorf("expected Gateway IP '192.168.1.1', got %q", snapshot.Gateway.IP)
	}
	if snapshot.Bandwidth == nil {
		t.Fatal("expected non-nil Bandwidth")
	}
	if snapshot.Bandwidth.RxMbps != 40.0 {
		t.Errorf("expected RxMbps 40.0, got %v", snapshot.Bandwidth.RxMbps)
	}
}

// TestBandwidthSampleFields verifies BandwidthSample struct fields.
func TestBandwidthSampleFields(t *testing.T) {
	tests := []struct {
		name          string
		sample        sap.BandwidthSample
		wantInterface string
		wantTxBytes   float64
		wantRxBytes   float64
		wantTxMbps    float64
		wantRxMbps    float64
		wantUtil      float64
	}{
		{
			name: "zero values",
			sample: sap.BandwidthSample{
				Interface: "",
			},
			wantInterface: "",
			wantTxBytes:   0,
			wantRxBytes:   0,
			wantTxMbps:    0,
			wantRxMbps:    0,
			wantUtil:      0,
		},
		{
			name: "typical values",
			sample: sap.BandwidthSample{
				Interface:     "eth0",
				TxBytesPerSec: 125000000, // 1 Gbps
				RxBytesPerSec: 62500000,  // 500 Mbps
				TxMbps:        1000.0,
				RxMbps:        500.0,
				Utilization:   50.0,
				SampledAt:     time.Now(),
			},
			wantInterface: "eth0",
			wantTxBytes:   125000000,
			wantRxBytes:   62500000,
			wantTxMbps:    1000.0,
			wantRxMbps:    500.0,
			wantUtil:      50.0,
		},
		{
			name: "high utilization",
			sample: sap.BandwidthSample{
				Interface:     "enp0s31f6",
				TxBytesPerSec: 118750000, // 950 Mbps
				RxBytesPerSec: 118750000, // 950 Mbps
				TxMbps:        950.0,
				RxMbps:        950.0,
				Utilization:   95.0,
				SampledAt:     time.Now(),
			},
			wantInterface: "enp0s31f6",
			wantTxBytes:   118750000,
			wantRxBytes:   118750000,
			wantTxMbps:    950.0,
			wantRxMbps:    950.0,
			wantUtil:      95.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.sample.Interface != tt.wantInterface {
				t.Errorf("Interface = %q, want %q", tt.sample.Interface, tt.wantInterface)
			}
			if tt.sample.TxBytesPerSec != tt.wantTxBytes {
				t.Errorf("TxBytesPerSec = %v, want %v", tt.sample.TxBytesPerSec, tt.wantTxBytes)
			}
			if tt.sample.RxBytesPerSec != tt.wantRxBytes {
				t.Errorf("RxBytesPerSec = %v, want %v", tt.sample.RxBytesPerSec, tt.wantRxBytes)
			}
			if tt.sample.TxMbps != tt.wantTxMbps {
				t.Errorf("TxMbps = %v, want %v", tt.sample.TxMbps, tt.wantTxMbps)
			}
			if tt.sample.RxMbps != tt.wantRxMbps {
				t.Errorf("RxMbps = %v, want %v", tt.sample.RxMbps, tt.wantRxMbps)
			}
			if tt.sample.Utilization != tt.wantUtil {
				t.Errorf("Utilization = %v, want %v", tt.sample.Utilization, tt.wantUtil)
			}
		})
	}
}

// TestSystemHealthFields verifies SystemHealth struct fields.
func TestSystemHealthFields(t *testing.T) {
	tests := []struct {
		name        string
		health      sap.SystemHealth
		wantCPU     float64
		wantMemory  float64
		wantDisk    float64
		wantTemp    float64
		wantUptime  time.Duration
		wantLoadAvg []float64
	}{
		{
			name:        "zero values",
			health:      sap.SystemHealth{},
			wantCPU:     0,
			wantMemory:  0,
			wantDisk:    0,
			wantTemp:    0,
			wantUptime:  0,
			wantLoadAvg: nil,
		},
		{
			name: "typical server values",
			health: sap.SystemHealth{
				CPUPercent:    25.5,
				MemoryPercent: 60.0,
				DiskPercent:   45.0,
				Temperature:   55.0,
				Uptime:        24 * time.Hour * 30, // 30 days
				LoadAverage:   []float64{1.5, 2.0, 1.8},
				SampledAt:     time.Now(),
			},
			wantCPU:     25.5,
			wantMemory:  60.0,
			wantDisk:    45.0,
			wantTemp:    55.0,
			wantUptime:  24 * time.Hour * 30,
			wantLoadAvg: []float64{1.5, 2.0, 1.8},
		},
		{
			name: "high load values",
			health: sap.SystemHealth{
				CPUPercent:    95.0,
				MemoryPercent: 90.0,
				DiskPercent:   85.0,
				Temperature:   80.0,
				Uptime:        time.Hour * 2,
				LoadAverage:   []float64{8.0, 7.5, 6.0},
				SampledAt:     time.Now(),
			},
			wantCPU:     95.0,
			wantMemory:  90.0,
			wantDisk:    85.0,
			wantTemp:    80.0,
			wantUptime:  time.Hour * 2,
			wantLoadAvg: []float64{8.0, 7.5, 6.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.health.CPUPercent != tt.wantCPU {
				t.Errorf("CPUPercent = %v, want %v", tt.health.CPUPercent, tt.wantCPU)
			}
			if tt.health.MemoryPercent != tt.wantMemory {
				t.Errorf("MemoryPercent = %v, want %v", tt.health.MemoryPercent, tt.wantMemory)
			}
			if tt.health.DiskPercent != tt.wantDisk {
				t.Errorf("DiskPercent = %v, want %v", tt.health.DiskPercent, tt.wantDisk)
			}
			if tt.health.Temperature != tt.wantTemp {
				t.Errorf("Temperature = %v, want %v", tt.health.Temperature, tt.wantTemp)
			}
			if tt.health.Uptime != tt.wantUptime {
				t.Errorf("Uptime = %v, want %v", tt.health.Uptime, tt.wantUptime)
			}
			if len(tt.health.LoadAverage) != len(tt.wantLoadAvg) {
				t.Errorf("LoadAverage length = %d, want %d",
					len(tt.health.LoadAverage), len(tt.wantLoadAvg))
			}
		})
	}
}

// TestTelemetrySnapshotWithSystemHealth verifies TelemetrySnapshot with SystemHealth.
func TestTelemetrySnapshotWithSystemHealth(t *testing.T) {
	now := time.Now()

	systemHealth := &sap.SystemHealth{
		CPUPercent:    35.5,
		MemoryPercent: 55.0,
		DiskPercent:   40.0,
		Temperature:   45.0,
		Uptime:        time.Hour * 48,
		LoadAverage:   []float64{1.0, 0.8, 0.9},
		SampledAt:     now,
	}

	snapshot := sap.TelemetrySnapshot{
		Timestamp:    now,
		SystemHealth: systemHealth,
	}

	if snapshot.SystemHealth == nil {
		t.Fatal("expected non-nil SystemHealth")
	}
	if snapshot.SystemHealth.CPUPercent != 35.5 {
		t.Errorf("expected CPUPercent 35.5, got %v", snapshot.SystemHealth.CPUPercent)
	}
	if snapshot.SystemHealth.Uptime != time.Hour*48 {
		t.Errorf("expected Uptime 48h, got %v", snapshot.SystemHealth.Uptime)
	}
	if len(snapshot.SystemHealth.LoadAverage) != 3 {
		t.Errorf("expected 3 LoadAverage values, got %d", len(snapshot.SystemHealth.LoadAverage))
	}
}

// TestLinkStatusFields verifies LinkStatus struct fields.
func TestLinkStatusFields(t *testing.T) {
	tests := []struct {
		name   string
		status sap.LinkStatus
	}{
		{
			name: "up interface with full data",
			status: sap.LinkStatus{
				Interface:  "eth0",
				State:      sap.LinkStateUp,
				Speed:      "1000Mbps",
				Duplex:     "full",
				MTU:        1500,
				MACAddress: "00:11:22:33:44:55",
				IPAddress:  "192.168.1.100",
				Gateway:    "192.168.1.1",
				Carrier:    true,
				TxBytes:    1000000,
				RxBytes:    5000000,
				TxPackets:  1000,
				RxPackets:  5000,
				TxErrors:   0,
				RxErrors:   0,
				TxDropped:  0,
				RxDropped:  0,
				UpdatedAt:  time.Now(),
			},
		},
		{
			name: "down interface",
			status: sap.LinkStatus{
				Interface:  "eth1",
				State:      sap.LinkStateDown,
				Speed:      "",
				Duplex:     "",
				MTU:        1500,
				MACAddress: "00:11:22:33:44:66",
				Carrier:    false,
				UpdatedAt:  time.Now(),
			},
		},
		{
			name: "interface with errors",
			status: sap.LinkStatus{
				Interface:  "eth2",
				State:      sap.LinkStateUp,
				Speed:      "100Mbps",
				Duplex:     "half",
				MTU:        1500,
				MACAddress: "00:11:22:33:44:77",
				Carrier:    true,
				TxBytes:    500000,
				RxBytes:    2500000,
				TxPackets:  500,
				RxPackets:  2500,
				TxErrors:   10,
				RxErrors:   5,
				TxDropped:  2,
				RxDropped:  1,
				UpdatedAt:  time.Now(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.status.Interface == "" {
				t.Error("expected non-empty Interface")
			}
			if tt.status.State == "" {
				t.Error("expected non-empty State")
			}
			if tt.status.UpdatedAt.IsZero() {
				t.Error("expected non-zero UpdatedAt")
			}
		})
	}
}

// TestLinkStateConstants verifies LinkState constant values.
func TestLinkStateConstants(t *testing.T) {
	tests := []struct {
		state    sap.LinkState
		expected string
	}{
		{sap.LinkStateUp, "up"},
		{sap.LinkStateDown, "down"},
		{sap.LinkStateDormant, "dormant"},
		{sap.LinkStateUnknown, "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if string(tt.state) != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, string(tt.state))
			}
		})
	}
}

// TestHealthStatusConstants verifies HealthStatus constant values.
func TestHealthStatusConstants(t *testing.T) {
	tests := []struct {
		status   sap.HealthStatus
		expected string
	}{
		{sap.HealthStatusHealthy, "healthy"},
		{sap.HealthStatusDegraded, "degraded"},
		{sap.HealthStatusUnhealthy, "unhealthy"},
		{sap.HealthStatusUnknown, "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if string(tt.status) != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, string(tt.status))
			}
		})
	}
}

// TestCableStatusConstants verifies CableStatus constant values.
func TestCableStatusConstants(t *testing.T) {
	tests := []struct {
		status   sap.CableStatus
		expected string
	}{
		{sap.CableStatusOK, "ok"},
		{sap.CableStatusOpen, "open"},
		{sap.CableStatusShort, "short"},
		{sap.CableStatusImpedance, "impedance_mismatch"},
		{sap.CableStatusUnknown, "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if string(tt.status) != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, string(tt.status))
			}
		})
	}
}

// TestGatewayHealthFields verifies GatewayHealth struct fields.
func TestGatewayHealthFields(t *testing.T) {
	tests := []struct {
		name   string
		health sap.GatewayHealth
	}{
		{
			name: "healthy gateway",
			health: sap.GatewayHealth{
				IP:         "192.168.1.1",
				Reachable:  true,
				RTT:        time.Millisecond * 2,
				RTTMs:      2.0,
				PacketLoss: 0.0,
				Jitter:     0.5,
				Status:     sap.HealthStatusHealthy,
				Uptime:     time.Hour * 24 * 7,
				LastCheck:  time.Now(),
			},
		},
		{
			name: "degraded gateway",
			health: sap.GatewayHealth{
				IP:         "10.0.0.1",
				Reachable:  true,
				RTT:        time.Millisecond * 100,
				RTTMs:      100.0,
				PacketLoss: 5.0,
				Jitter:     10.0,
				Status:     sap.HealthStatusDegraded,
				LastCheck:  time.Now(),
			},
		},
		{
			name: "unreachable gateway",
			health: sap.GatewayHealth{
				IP:         "172.16.0.1",
				Reachable:  false,
				RTT:        0,
				RTTMs:      0,
				PacketLoss: 100.0,
				Jitter:     0,
				Status:     sap.HealthStatusUnhealthy,
				LastCheck:  time.Now(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.health.IP == "" {
				t.Error("expected non-empty IP")
			}
			if tt.health.LastCheck.IsZero() {
				t.Error("expected non-zero LastCheck")
			}
			if tt.health.Status == "" {
				t.Error("expected non-empty Status")
			}
		})
	}
}

// TestDNSTestResultFields verifies DNSTestResult struct fields.
func TestDNSTestResultFields(t *testing.T) {
	tests := []struct {
		name   string
		result sap.DNSTestResult
	}{
		{
			name: "successful DNS query",
			result: sap.DNSTestResult{
				Query:        "google.com",
				Server:       "8.8.8.8",
				Success:      true,
				Answers:      []sap.DNSAnswer{{Name: "google.com", Type: "A", Value: "142.250.80.46", TTL: 300}},
				ResponseTime: time.Millisecond * 25,
				ResponseMs:   25.0,
				DNSSEC:       false,
				TestedAt:     time.Now(),
			},
		},
		{
			name: "failed DNS query",
			result: sap.DNSTestResult{
				Query:        "nonexistent.example.invalid",
				Server:       "8.8.8.8",
				Success:      false,
				Error:        "NXDOMAIN",
				ResponseTime: time.Millisecond * 50,
				ResponseMs:   50.0,
				TestedAt:     time.Now(),
			},
		},
		{
			name: "authoritative response",
			result: sap.DNSTestResult{
				Query:         "example.com",
				Server:        "ns1.example.com",
				Success:       true,
				Answers:       []sap.DNSAnswer{{Name: "example.com", Type: "A", Value: "93.184.216.34", TTL: 3600}},
				ResponseTime:  time.Millisecond * 10,
				ResponseMs:    10.0,
				DNSSEC:        true,
				Authoritative: true,
				TestedAt:      time.Now(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.result.Query == "" {
				t.Error("expected non-empty Query")
			}
			if tt.result.Server == "" {
				t.Error("expected non-empty Server")
			}
			if tt.result.TestedAt.IsZero() {
				t.Error("expected non-zero TestedAt")
			}
		})
	}
}

// TestDNSAnswerFields verifies DNSAnswer struct fields.
func TestDNSAnswerFields(t *testing.T) {
	tests := []struct {
		name   string
		answer sap.DNSAnswer
	}{
		{
			name:   "A record",
			answer: sap.DNSAnswer{Name: "example.com", Type: "A", Value: "93.184.216.34", TTL: 3600},
		},
		{
			name: "AAAA record",
			answer: sap.DNSAnswer{
				Name:  "example.com",
				Type:  "AAAA",
				Value: "2606:2800:220:1:248:1893:25c8:1946",
				TTL:   3600,
			},
		},
		{
			name:   "CNAME record",
			answer: sap.DNSAnswer{Name: "www.example.com", Type: "CNAME", Value: "example.com", TTL: 1800},
		},
		{
			name:   "MX record",
			answer: sap.DNSAnswer{Name: "example.com", Type: "MX", Value: "mail.example.com", TTL: 7200},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.answer.Name == "" {
				t.Error("expected non-empty Name")
			}
			if tt.answer.Type == "" {
				t.Error("expected non-empty Type")
			}
			if tt.answer.Value == "" {
				t.Error("expected non-empty Value")
			}
			if tt.answer.TTL < 0 {
				t.Errorf("expected non-negative TTL, got %d", tt.answer.TTL)
			}
		})
	}
}

// TestDHCPTestResultFields verifies DHCPTestResult struct fields.
func TestDHCPTestResultFields(t *testing.T) {
	tests := []struct {
		name   string
		result sap.DHCPTestResult
	}{
		{
			name: "successful DHCP lease",
			result: sap.DHCPTestResult{
				Success:      true,
				ServerIP:     "192.168.1.1",
				OfferedIP:    "192.168.1.100",
				SubnetMask:   "255.255.255.0",
				Gateway:      "192.168.1.1",
				DNSServers:   []string{"8.8.8.8", "8.8.4.4"},
				LeaseTime:    time.Hour * 24,
				LeaseTimeSec: 86400,
				ResponseTime: time.Millisecond * 50,
				ResponseMs:   50.0,
				TestedAt:     time.Now(),
			},
		},
		{
			name: "failed DHCP request",
			result: sap.DHCPTestResult{
				Success:      false,
				Error:        "no DHCP server responding",
				ResponseTime: time.Second * 5,
				ResponseMs:   5000.0,
				TestedAt:     time.Now(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.result.TestedAt.IsZero() {
				t.Error("expected non-zero TestedAt")
			}
			if tt.result.Success && tt.result.ServerIP == "" {
				t.Error("expected non-empty ServerIP for successful result")
			}
			if !tt.result.Success && tt.result.Error == "" {
				t.Error("expected non-empty Error for failed result")
			}
		})
	}
}

// TestSpeedtestResultFields verifies SpeedtestResult struct fields.
func TestSpeedtestResultFields(t *testing.T) {
	tests := []struct {
		name   string
		result sap.SpeedtestResult
	}{
		{
			name: "typical home connection",
			result: sap.SpeedtestResult{
				DownloadMbps: 100.5,
				UploadMbps:   20.3,
				PingMs:       15.0,
				JitterMs:     2.5,
				ServerName:   "Speedtest Server - New York",
				ServerID:     "12345",
				ISP:          "Example ISP",
				TestDuration: time.Second * 30,
				TestedAt:     time.Now(),
			},
		},
		{
			name: "gigabit fiber connection",
			result: sap.SpeedtestResult{
				DownloadMbps: 940.2,
				UploadMbps:   925.8,
				PingMs:       3.0,
				JitterMs:     0.5,
				ServerName:   "Fiber Server - Los Angeles",
				ServerID:     "67890",
				ISP:          "Fiber Provider",
				TestDuration: time.Second * 20,
				TestedAt:     time.Now(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.result.DownloadMbps < 0 {
				t.Errorf("expected non-negative DownloadMbps, got %v", tt.result.DownloadMbps)
			}
			if tt.result.UploadMbps < 0 {
				t.Errorf("expected non-negative UploadMbps, got %v", tt.result.UploadMbps)
			}
			if tt.result.TestedAt.IsZero() {
				t.Error("expected non-zero TestedAt")
			}
		})
	}
}

// TestIPerfResultFields verifies IPerfResult struct fields.
func TestIPerfResultFields(t *testing.T) {
	tests := []struct {
		name   string
		result sap.IPerfResult
	}{
		{
			name: "TCP download test",
			result: sap.IPerfResult{
				Protocol:      "tcp",
				Direction:     "receive",
				BandwidthMbps: 940.5,
				TransferMB:    1175.6,
				Duration:      time.Second * 10,
				DurationSec:   10.0,
				Retransmits:   5,
				ServerAddr:    "192.168.1.100",
				TestedAt:      time.Now(),
			},
		},
		{
			name: "UDP bidirectional test",
			result: sap.IPerfResult{
				Protocol:      "udp",
				Direction:     "bidirectional",
				BandwidthMbps: 850.0,
				TransferMB:    1062.5,
				Duration:      time.Second * 10,
				DurationSec:   10.0,
				Jitter:        0.25,
				PacketLoss:    0.1,
				ServerAddr:    "10.0.0.50",
				TestedAt:      time.Now(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.result.Protocol != "tcp" && tt.result.Protocol != "udp" {
				t.Errorf("expected Protocol 'tcp' or 'udp', got %q", tt.result.Protocol)
			}
			if tt.result.BandwidthMbps < 0 {
				t.Errorf("expected non-negative BandwidthMbps, got %v", tt.result.BandwidthMbps)
			}
			if tt.result.ServerAddr == "" {
				t.Error("expected non-empty ServerAddr")
			}
			if tt.result.TestedAt.IsZero() {
				t.Error("expected non-zero TestedAt")
			}
		})
	}
}

// TestVLANConfigFields verifies VLANConfig struct fields.
func TestVLANConfigFields(t *testing.T) {
	tests := []struct {
		name   string
		config sap.VLANConfig
	}{
		{
			name: "tagged VLAN",
			config: sap.VLANConfig{
				ID:          100,
				Name:        "Management",
				Interface:   "eth0",
				IPAddress:   "192.168.100.1",
				SubnetMask:  "255.255.255.0",
				Gateway:     "192.168.100.254",
				Tagged:      true,
				MemberPorts: []string{"eth0", "eth1"},
			},
		},
		{
			name: "native VLAN",
			config: sap.VLANConfig{
				ID:        1,
				Name:      "Native",
				Interface: "eth0",
				Tagged:    false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.config.ID < 1 || tt.config.ID > 4094 {
				t.Errorf("expected VLAN ID between 1-4094, got %d", tt.config.ID)
			}
			if tt.config.Interface == "" {
				t.Error("expected non-empty Interface")
			}
		})
	}
}

// TestCableTestResultFields verifies CableTestResult struct fields.
func TestCableTestResultFields(t *testing.T) {
	tests := []struct {
		name   string
		result sap.CableTestResult
	}{
		{
			name: "cable OK",
			result: sap.CableTestResult{
				Interface: "eth0",
				Status:    sap.CableStatusOK,
				Length:    25.5,
				PairResults: []sap.PairResult{
					{Pair: 1, Status: sap.CableStatusOK, Length: 25.5, Impedance: 100.0},
					{Pair: 2, Status: sap.CableStatusOK, Length: 25.5, Impedance: 100.0},
					{Pair: 3, Status: sap.CableStatusOK, Length: 25.5, Impedance: 100.0},
					{Pair: 4, Status: sap.CableStatusOK, Length: 25.5, Impedance: 100.0},
				},
				TestedAt: time.Now(),
			},
		},
		{
			name: "cable with open pair",
			result: sap.CableTestResult{
				Interface: "eth1",
				Status:    sap.CableStatusOpen,
				Length:    10.0,
				PairResults: []sap.PairResult{
					{Pair: 1, Status: sap.CableStatusOK, Length: 10.0},
					{Pair: 2, Status: sap.CableStatusOpen, Length: 10.0},
					{Pair: 3, Status: sap.CableStatusOK, Length: 10.0},
					{Pair: 4, Status: sap.CableStatusOK, Length: 10.0},
				},
				TestedAt: time.Now(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.result.Interface == "" {
				t.Error("expected non-empty Interface")
			}
			if tt.result.Status == "" {
				t.Error("expected non-empty Status")
			}
			if tt.result.TestedAt.IsZero() {
				t.Error("expected non-zero TestedAt")
			}
		})
	}
}

// TestPairResultFields verifies PairResult struct fields.
func TestPairResultFields(t *testing.T) {
	tests := []struct {
		name   string
		result sap.PairResult
	}{
		{
			name:   "OK pair with full data",
			result: sap.PairResult{Pair: 1, Status: sap.CableStatusOK, Length: 30.0, Impedance: 100.0},
		},
		{
			name:   "shorted pair",
			result: sap.PairResult{Pair: 2, Status: sap.CableStatusShort, Length: 5.0},
		},
		{
			name:   "impedance mismatch",
			result: sap.PairResult{Pair: 3, Status: sap.CableStatusImpedance, Length: 20.0, Impedance: 75.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.result.Pair < 1 || tt.result.Pair > 4 {
				t.Errorf("expected Pair between 1-4, got %d", tt.result.Pair)
			}
			if tt.result.Status == "" {
				t.Error("expected non-empty Status")
			}
		})
	}
}

// TestSNMPDeviceFields verifies SNMPDevice struct fields.
func TestSNMPDeviceFields(t *testing.T) {
	tests := []struct {
		name   string
		device sap.SNMPDevice
	}{
		{
			name: "switch with full data",
			device: sap.SNMPDevice{
				IP:          "192.168.1.2",
				SysName:     "core-switch-01",
				SysDescr:    "Cisco IOS Software, Version 15.2",
				SysLocation: "Server Room A",
				SysContact:  "admin@example.com",
				SysUpTime:   time.Hour * 24 * 365,
				Interfaces: []sap.SNMPInterface{
					{Index: 1, Name: "GigabitEthernet0/1", Type: "ethernetCsmacd", Speed: 1000000000},
				},
				VLANs: []sap.SNMPVLAN{
					{ID: 1, Name: "default", Status: "active"},
					{ID: 100, Name: "Management", Status: "active"},
				},
				MACTable: []sap.MACTableEntry{
					{MACAddress: "00:11:22:33:44:55", Port: 1, VLANID: 1, Type: "dynamic"},
				},
				CollectedAt: time.Now(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.device.IP == "" {
				t.Error("expected non-empty IP")
			}
			if tt.device.CollectedAt.IsZero() {
				t.Error("expected non-zero CollectedAt")
			}
		})
	}
}

// TestSNMPInterfaceFields verifies SNMPInterface struct fields.
func TestSNMPInterfaceFields(t *testing.T) {
	iface := sap.SNMPInterface{
		Index:       1,
		Name:        "GigabitEthernet0/1",
		Description: "Uplink to Router",
		Type:        "ethernetCsmacd",
		Speed:       1000000000,
		AdminStatus: "up",
		OperStatus:  "up",
		InOctets:    1000000000,
		OutOctets:   500000000,
		InErrors:    0,
		OutErrors:   0,
	}

	if iface.Index < 1 {
		t.Errorf("expected positive Index, got %d", iface.Index)
	}
	if iface.Name == "" {
		t.Error("expected non-empty Name")
	}
	if iface.Type == "" {
		t.Error("expected non-empty Type")
	}
}

// TestSNMPVLANFields verifies SNMPVLAN struct fields.
func TestSNMPVLANFields(t *testing.T) {
	vlan := sap.SNMPVLAN{
		ID:     100,
		Name:   "Management",
		Status: "active",
		Ports:  []int{1, 2, 3, 4},
	}

	if vlan.ID < 1 || vlan.ID > 4094 {
		t.Errorf("expected VLAN ID between 1-4094, got %d", vlan.ID)
	}
	if vlan.Name == "" {
		t.Error("expected non-empty Name")
	}
}

// TestMACTableEntryFields verifies MACTableEntry struct fields.
func TestMACTableEntryFields(t *testing.T) {
	tests := []struct {
		name  string
		entry sap.MACTableEntry
	}{
		{
			name:  "dynamic entry",
			entry: sap.MACTableEntry{MACAddress: "00:11:22:33:44:55", Port: 1, VLANID: 1, Type: "dynamic"},
		},
		{
			name:  "static entry",
			entry: sap.MACTableEntry{MACAddress: "66:77:88:99:AA:BB", Port: 24, VLANID: 100, Type: "static"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.entry.MACAddress == "" {
				t.Error("expected non-empty MACAddress")
			}
			if tt.entry.Port < 1 {
				t.Errorf("expected positive Port, got %d", tt.entry.Port)
			}
			if tt.entry.Type != "dynamic" && tt.entry.Type != "static" {
				t.Errorf("expected Type 'dynamic' or 'static', got %q", tt.entry.Type)
			}
		})
	}
}

// TestErrorConstants verifies error constant values.
func TestErrorConstants(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{"ErrNotImplemented", sap.ErrNotImplemented},
		{"ErrNotInitialized", sap.ErrNotInitialized},
		{"ErrNotSupported", sap.ErrNotSupported},
		{"ErrTestFailed", sap.ErrTestFailed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Errorf("%s should not be nil", tt.name)
			}
			if tt.err.Error() == "" {
				t.Errorf("%s should have non-empty message", tt.name)
			}
		})
	}
}
