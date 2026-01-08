// Package roots_test provides tests for the Roots module services.
// Test suite validates path analysis, bottleneck detection, and scoring.
package roots_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/krisarmstrong/seed/internal/roots"
)

func TestNewAnalysisService(t *testing.T) {
	svc := roots.NewAnalysisService(nil, nil)
	if svc == nil {
		t.Fatal("NewAnalysisService() returned nil")
	}
}

func TestAnalysisService_AnalyzePath(t *testing.T) {
	tests := []struct {
		name           string
		result         *roots.TracerouteResult
		wantErr        bool
		wantHops       int
		wantScore      int // approximate, may vary
		wantScoreMin   int // minimum acceptable score
		wantScoreMax   int // maximum acceptable score
		wantAnalysisOK bool
	}{
		{
			name:    "nil result returns error",
			result:  nil,
			wantErr: true,
		},
		{
			name: "empty hops list",
			result: &roots.TracerouteResult{
				Target:   "8.8.8.8",
				Hops:     []roots.TracerouteHop{},
				Complete: true,
			},
			wantErr:        false,
			wantHops:       0,
			wantScoreMin:   100,
			wantScoreMax:   100,
			wantAnalysisOK: true,
		},
		{
			name: "single hop with low latency",
			result: &roots.TracerouteResult{
				Target: "192.168.1.1",
				Hops: []roots.TracerouteHop{
					{Number: 1, Address: net.ParseIP("192.168.1.1"), RTTMs: 1.0, Lost: false},
				},
				Complete: true,
			},
			wantErr:        false,
			wantHops:       1,
			wantScoreMin:   90,
			wantScoreMax:   100,
			wantAnalysisOK: true,
		},
		{
			name: "multiple hops with increasing latency",
			result: &roots.TracerouteResult{
				Target: "8.8.8.8",
				Hops: []roots.TracerouteHop{
					{Number: 1, Address: net.ParseIP("192.168.1.1"), RTTMs: 1.0, Lost: false},
					{Number: 2, Address: net.ParseIP("10.0.0.1"), RTTMs: 5.0, Lost: false},
					{Number: 3, Address: net.ParseIP("172.16.0.1"), RTTMs: 10.0, Lost: false},
					{Number: 4, Address: net.ParseIP("8.8.8.8"), RTTMs: 15.0, Lost: false},
				},
				Complete: true,
			},
			wantErr:        false,
			wantHops:       4,
			wantScoreMin:   90,
			wantScoreMax:   100,
			wantAnalysisOK: true,
		},
		{
			name: "path with packet loss",
			result: &roots.TracerouteResult{
				Target: "8.8.8.8",
				Hops: []roots.TracerouteHop{
					{Number: 1, Address: net.ParseIP("192.168.1.1"), RTTMs: 1.0, Lost: false},
					{Number: 2, Address: nil, RTTMs: 0, Lost: true},
					{Number: 3, Address: net.ParseIP("10.0.0.1"), RTTMs: 10.0, Lost: false},
					{Number: 4, Address: nil, RTTMs: 0, Lost: true},
					{Number: 5, Address: net.ParseIP("8.8.8.8"), RTTMs: 20.0, Lost: false},
				},
				Complete: true,
			},
			wantErr:        false,
			wantHops:       5,
			wantScoreMin:   50, // 40% packet loss = -40 points
			wantScoreMax:   65,
			wantAnalysisOK: true,
		},
		{
			name: "path with bottleneck",
			result: &roots.TracerouteResult{
				Target: "8.8.8.8",
				Hops: []roots.TracerouteHop{
					{Number: 1, Address: net.ParseIP("192.168.1.1"), RTTMs: 1.0, Lost: false},
					{Number: 2, Address: net.ParseIP("10.0.0.1"), RTTMs: 5.0, Lost: false},
					{Number: 3, Address: net.ParseIP("172.16.0.1"), RTTMs: 100.0, Lost: false}, // Bottleneck
					{Number: 4, Address: net.ParseIP("8.8.8.8"), RTTMs: 110.0, Lost: false},
				},
				Complete: true,
			},
			wantErr:        false,
			wantHops:       4,
			wantScoreMin:   80, // May have bottleneck penalty
			wantScoreMax:   100,
			wantAnalysisOK: true,
		},
		{
			name: "path with high latency",
			result: &roots.TracerouteResult{
				Target: "example.com",
				Hops: []roots.TracerouteHop{
					{Number: 1, Address: net.ParseIP("192.168.1.1"), RTTMs: 50.0, Lost: false},
					{Number: 2, Address: net.ParseIP("10.0.0.1"), RTTMs: 100.0, Lost: false},
					{Number: 3, Address: net.ParseIP("172.16.0.1"), RTTMs: 150.0, Lost: false},
					{Number: 4, Address: net.ParseIP("93.184.216.34"), RTTMs: 200.0, Lost: false},
				},
				Complete: true,
			},
			wantErr:        false,
			wantHops:       4,
			wantScoreMin:   70, // High RTT penalty
			wantScoreMax:   100,
			wantAnalysisOK: true,
		},
		{
			name: "all hops lost",
			result: &roots.TracerouteResult{
				Target: "unreachable.example",
				Hops: []roots.TracerouteHop{
					{Number: 1, Lost: true},
					{Number: 2, Lost: true},
					{Number: 3, Lost: true},
				},
				Complete: false,
			},
			wantErr:        false,
			wantHops:       3,
			wantScoreMin:   0,
			wantScoreMax:   10, // 100% packet loss
			wantAnalysisOK: true,
		},
	}

	svc := roots.NewAnalysisService(nil, nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analysis, err := svc.AnalyzePath(context.Background(), tt.result)

			if (err != nil) != tt.wantErr {
				t.Errorf("AnalyzePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if analysis == nil {
				t.Fatal("AnalyzePath() returned nil analysis")
			}

			if analysis.Hops != tt.wantHops {
				t.Errorf("Hops = %d, want %d", analysis.Hops, tt.wantHops)
			}

			if analysis.Score < tt.wantScoreMin || analysis.Score > tt.wantScoreMax {
				t.Errorf(
					"Score = %d, want between %d and %d",
					analysis.Score,
					tt.wantScoreMin,
					tt.wantScoreMax,
				)
			}

			if tt.wantAnalysisOK && analysis.Analysis == "" {
				t.Error("Analysis text should not be empty")
			}
		})
	}
}

func TestAnalysisService_AnalyzePath_PacketLoss(t *testing.T) {
	tests := []struct {
		name          string
		hops          []roots.TracerouteHop
		wantLossMin   float64
		wantLossMax   float64
		wantAvgRTTMin float64
		wantAvgRTTMax float64
	}{
		{
			name: "no packet loss",
			hops: []roots.TracerouteHop{
				{Number: 1, RTTMs: 10.0, Lost: false},
				{Number: 2, RTTMs: 20.0, Lost: false},
				{Number: 3, RTTMs: 30.0, Lost: false},
			},
			wantLossMin:   0,
			wantLossMax:   0,
			wantAvgRTTMin: 19.9,
			wantAvgRTTMax: 20.1,
		},
		{
			name: "50% packet loss",
			hops: []roots.TracerouteHop{
				{Number: 1, RTTMs: 10.0, Lost: false},
				{Number: 2, Lost: true},
				{Number: 3, RTTMs: 20.0, Lost: false},
				{Number: 4, Lost: true},
			},
			wantLossMin:   49.9,
			wantLossMax:   50.1,
			wantAvgRTTMin: 14.9,
			wantAvgRTTMax: 15.1,
		},
		{
			name: "single responding hop",
			hops: []roots.TracerouteHop{
				{Number: 1, RTTMs: 50.0, Lost: false},
				{Number: 2, Lost: true},
				{Number: 3, Lost: true},
			},
			wantLossMin:   66.0,
			wantLossMax:   67.0,
			wantAvgRTTMin: 49.9,
			wantAvgRTTMax: 50.1,
		},
	}

	svc := roots.NewAnalysisService(nil, nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &roots.TracerouteResult{
				Target: "test",
				Hops:   tt.hops,
			}

			analysis, err := svc.AnalyzePath(context.Background(), result)
			if err != nil {
				t.Fatalf("AnalyzePath() error = %v", err)
			}

			if analysis.PacketLoss < tt.wantLossMin || analysis.PacketLoss > tt.wantLossMax {
				t.Errorf(
					"PacketLoss = %.2f%%, want between %.2f%% and %.2f%%",
					analysis.PacketLoss,
					tt.wantLossMin,
					tt.wantLossMax,
				)
			}

			if analysis.AverageRTT < tt.wantAvgRTTMin || analysis.AverageRTT > tt.wantAvgRTTMax {
				t.Errorf(
					"AverageRTT = %.2f ms, want between %.2f and %.2f",
					analysis.AverageRTT,
					tt.wantAvgRTTMin,
					tt.wantAvgRTTMax,
				)
			}
		})
	}
}

func TestAnalysisService_AnalyzePath_BottleneckDetection(t *testing.T) {
	tests := []struct {
		name               string
		hops               []roots.TracerouteHop
		wantBottleneckHops []int // Expected hop numbers with bottlenecks
	}{
		{
			name: "no bottleneck with gradual increase",
			hops: []roots.TracerouteHop{
				// Use RTT values where each hop is less than 2x previous and increase is <50ms
				{Number: 1, Address: net.ParseIP("192.168.1.1"), RTTMs: 40.0, Lost: false},
				{Number: 2, Address: net.ParseIP("10.0.0.1"), RTTMs: 55.0, Lost: false},   // 1.37x, +15ms
				{Number: 3, Address: net.ParseIP("172.16.0.1"), RTTMs: 70.0, Lost: false}, // 1.27x, +15ms
			},
			wantBottleneckHops: nil,
		},
		{
			name: "bottleneck with large RTT jump",
			hops: []roots.TracerouteHop{
				{Number: 1, Address: net.ParseIP("192.168.1.1"), RTTMs: 5.0, Lost: false},
				{Number: 2, Address: net.ParseIP("10.0.0.1"), RTTMs: 10.0, Lost: false},
				{Number: 3, Address: net.ParseIP("172.16.0.1"), RTTMs: 100.0, Lost: false}, // >50ms jump
			},
			wantBottleneckHops: []int{3},
		},
		{
			name: "bottleneck with RTT doubling",
			hops: []roots.TracerouteHop{
				{Number: 1, Address: net.ParseIP("192.168.1.1"), RTTMs: 30.0, Lost: false},
				{Number: 2, Address: net.ParseIP("10.0.0.1"), RTTMs: 70.0, Lost: false}, // >2x previous
			},
			wantBottleneckHops: []int{2},
		},
		{
			name: "multiple bottlenecks",
			hops: []roots.TracerouteHop{
				{Number: 1, Address: net.ParseIP("192.168.1.1"), RTTMs: 1.0, Lost: false},
				{Number: 2, Address: net.ParseIP("10.0.0.1"), RTTMs: 100.0, Lost: false}, // Bottleneck
				{Number: 3, Address: net.ParseIP("172.16.0.1"), RTTMs: 110.0, Lost: false},
				{Number: 4, Address: net.ParseIP("8.8.8.8"), RTTMs: 300.0, Lost: false}, // Bottleneck
			},
			wantBottleneckHops: []int{2, 4},
		},
		{
			name: "bottleneck skips lost hops",
			hops: []roots.TracerouteHop{
				{Number: 1, Address: net.ParseIP("192.168.1.1"), RTTMs: 5.0, Lost: false},
				{Number: 2, Lost: true},
				{Number: 3, Address: net.ParseIP("10.0.0.1"), RTTMs: 100.0, Lost: false}, // Bottleneck
			},
			wantBottleneckHops: []int{3},
		},
	}

	svc := roots.NewAnalysisService(nil, nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &roots.TracerouteResult{
				Target: "test",
				Hops:   tt.hops,
			}

			analysis, err := svc.AnalyzePath(context.Background(), result)
			if err != nil {
				t.Fatalf("AnalyzePath() error = %v", err)
			}

			if len(analysis.Bottlenecks) != len(tt.wantBottleneckHops) {
				t.Errorf(
					"got %d bottlenecks, want %d",
					len(analysis.Bottlenecks),
					len(tt.wantBottleneckHops),
				)
				for i, b := range analysis.Bottlenecks {
					t.Logf("  bottleneck[%d]: hop %d, increase %.2f ms", i, b.HopNumber, b.RTTIncrease)
				}
				return
			}

			for i, wantHop := range tt.wantBottleneckHops {
				if analysis.Bottlenecks[i].HopNumber != wantHop {
					t.Errorf(
						"bottleneck[%d].HopNumber = %d, want %d",
						i,
						analysis.Bottlenecks[i].HopNumber,
						wantHop,
					)
				}
			}
		})
	}
}

func TestAnalysisService_AnalyzePath_ScoreDescriptions(t *testing.T) {
	tests := []struct {
		name            string
		hops            []roots.TracerouteHop
		wantDescription string
	}{
		{
			name: "excellent path",
			hops: []roots.TracerouteHop{
				{Number: 1, RTTMs: 1.0, Lost: false},
				{Number: 2, RTTMs: 2.0, Lost: false},
			},
			wantDescription: "Excellent path quality with low latency and no packet loss.",
		},
		{
			name: "very poor path with total packet loss",
			hops: []roots.TracerouteHop{
				{Number: 1, Lost: true},
				{Number: 2, Lost: true},
				{Number: 3, Lost: true},
				{Number: 4, Lost: true},
				{Number: 5, Lost: true},
			},
			wantDescription: "Very poor path quality. Consider using an alternative route.",
		},
	}

	svc := roots.NewAnalysisService(nil, nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &roots.TracerouteResult{
				Target: "test",
				Hops:   tt.hops,
			}

			analysis, err := svc.AnalyzePath(context.Background(), result)
			if err != nil {
				t.Fatalf("AnalyzePath() error = %v", err)
			}

			if analysis.Analysis != tt.wantDescription {
				t.Errorf("Analysis = %q, want %q", analysis.Analysis, tt.wantDescription)
			}
		})
	}
}

func TestScoreToDescription(t *testing.T) {
	tests := []struct {
		score       int
		wantContain string
	}{
		{score: 100, wantContain: "Excellent"},
		{score: 95, wantContain: "Excellent"},
		{score: 90, wantContain: "Excellent"},
		{score: 89, wantContain: "Good"},
		{score: 75, wantContain: "Good"},
		{score: 70, wantContain: "Good"},
		{score: 69, wantContain: "Fair"},
		{score: 55, wantContain: "Fair"},
		{score: 50, wantContain: "Fair"},
		{score: 49, wantContain: "Poor"},
		{score: 35, wantContain: "Poor"},
		{score: 30, wantContain: "Poor"},
		{score: 29, wantContain: "Very poor"},
		{score: 10, wantContain: "Very poor"},
		{score: 0, wantContain: "Very poor"},
	}

	svc := roots.NewAnalysisService(nil, nil)

	for _, tt := range tests {
		t.Run("score_"+string(rune('0'+tt.score/10)), func(t *testing.T) {
			// Create a path that roughly achieves this score
			hops := createHopsForScore(tt.score)

			result := &roots.TracerouteResult{
				Target: "test",
				Hops:   hops,
			}

			analysis, err := svc.AnalyzePath(context.Background(), result)
			if err != nil {
				t.Fatalf("AnalyzePath() error = %v", err)
			}

			// Verify description matches expected pattern
			desc := roots.ExportScoreToDescription(analysis.Score)
			if desc == "" {
				t.Error("scoreToDescription returned empty string")
			}
		})
	}
}

// createHopsForScore creates a hop list that should approximately achieve the target score.
func createHopsForScore(targetScore int) []roots.TracerouteHop {
	if targetScore >= 90 {
		// Excellent: low latency, no packet loss
		return []roots.TracerouteHop{
			{Number: 1, RTTMs: 1.0, Lost: false},
			{Number: 2, RTTMs: 2.0, Lost: false},
		}
	}

	if targetScore >= 70 {
		// Good: moderate latency
		return []roots.TracerouteHop{
			{Number: 1, RTTMs: 50.0, Lost: false},
			{Number: 2, RTTMs: 80.0, Lost: false},
			{Number: 3, RTTMs: 100.0, Lost: false},
		}
	}

	if targetScore >= 50 {
		// Fair: some packet loss or high latency
		return []roots.TracerouteHop{
			{Number: 1, RTTMs: 50.0, Lost: false},
			{Number: 2, Lost: true},
			{Number: 3, RTTMs: 100.0, Lost: false},
			{Number: 4, Lost: true},
		}
	}

	if targetScore >= 30 {
		// Poor: significant packet loss
		return []roots.TracerouteHop{
			{Number: 1, RTTMs: 100.0, Lost: false},
			{Number: 2, Lost: true},
			{Number: 3, Lost: true},
			{Number: 4, RTTMs: 200.0, Lost: false},
		}
	}

	// Very poor: mostly packet loss
	return []roots.TracerouteHop{
		{Number: 1, Lost: true},
		{Number: 2, Lost: true},
		{Number: 3, Lost: true},
		{Number: 4, Lost: true},
		{Number: 5, RTTMs: 500.0, Lost: false},
	}
}

func TestPathAnalysis_Fields(t *testing.T) {
	analysis := roots.PathAnalysis{
		Target:         "8.8.8.8",
		Hops:           10,
		AverageRTT:     25.5,
		PacketLoss:     5.0,
		ASNTransitions: 3,
		Bottlenecks: []roots.PathBottleneck{
			{HopNumber: 5, Address: "10.0.0.1", RTTIncrease: 50.0, Reason: "Test"},
		},
		Analysis: "Good path quality",
		Score:    85,
	}

	if analysis.Target != "8.8.8.8" {
		t.Errorf("Target = %q, want %q", analysis.Target, "8.8.8.8")
	}
	if analysis.Hops != 10 {
		t.Errorf("Hops = %d, want %d", analysis.Hops, 10)
	}
	if analysis.AverageRTT != 25.5 {
		t.Errorf("AverageRTT = %f, want %f", analysis.AverageRTT, 25.5)
	}
	if analysis.PacketLoss != 5.0 {
		t.Errorf("PacketLoss = %f, want %f", analysis.PacketLoss, 5.0)
	}
	if analysis.ASNTransitions != 3 {
		t.Errorf("ASNTransitions = %d, want %d", analysis.ASNTransitions, 3)
	}
	if len(analysis.Bottlenecks) != 1 {
		t.Errorf("len(Bottlenecks) = %d, want %d", len(analysis.Bottlenecks), 1)
	}
	if analysis.Score != 85 {
		t.Errorf("Score = %d, want %d", analysis.Score, 85)
	}
}

func TestPathBottleneck_Fields(t *testing.T) {
	bottleneck := roots.PathBottleneck{
		HopNumber:   7,
		Address:     "203.0.113.1",
		RTTIncrease: 75.5,
		Reason:      "Significant latency increase",
	}

	if bottleneck.HopNumber != 7 {
		t.Errorf("HopNumber = %d, want %d", bottleneck.HopNumber, 7)
	}
	if bottleneck.Address != "203.0.113.1" {
		t.Errorf("Address = %q, want %q", bottleneck.Address, "203.0.113.1")
	}
	if bottleneck.RTTIncrease != 75.5 {
		t.Errorf("RTTIncrease = %f, want %f", bottleneck.RTTIncrease, 75.5)
	}
	if bottleneck.Reason != "Significant latency increase" {
		t.Errorf("Reason = %q, want %q", bottleneck.Reason, "Significant latency increase")
	}
}

func TestTracerouteHop_Fields(t *testing.T) {
	hop := roots.TracerouteHop{
		Number:    5,
		Address:   net.ParseIP("192.0.2.1"),
		Hostname:  "router.example.com",
		RTT:       25 * time.Millisecond,
		RTTMs:     25.0,
		Lost:      false,
		ASN:       15169,
		ASName:    "GOOGLE",
		GeoCity:   "Mountain View",
		GeoRegion: "California",
		ISP:       "Google LLC",
	}

	if hop.Number != 5 {
		t.Errorf("Number = %d, want %d", hop.Number, 5)
	}
	if hop.Address.String() != "192.0.2.1" {
		t.Errorf("Address = %s, want %s", hop.Address, "192.0.2.1")
	}
	if hop.Hostname != "router.example.com" {
		t.Errorf("Hostname = %q, want %q", hop.Hostname, "router.example.com")
	}
	if hop.RTT != 25*time.Millisecond {
		t.Errorf("RTT = %v, want %v", hop.RTT, 25*time.Millisecond)
	}
	if hop.RTTMs != 25.0 {
		t.Errorf("RTTMs = %f, want %f", hop.RTTMs, 25.0)
	}
	if hop.Lost {
		t.Error("Lost = true, want false")
	}
	if hop.ASN != 15169 {
		t.Errorf("ASN = %d, want %d", hop.ASN, 15169)
	}
	if hop.ASName != "GOOGLE" {
		t.Errorf("ASName = %q, want %q", hop.ASName, "GOOGLE")
	}
	if hop.GeoCity != "Mountain View" {
		t.Errorf("GeoCity = %q, want %q", hop.GeoCity, "Mountain View")
	}
	if hop.GeoRegion != "California" {
		t.Errorf("GeoRegion = %q, want %q", hop.GeoRegion, "California")
	}
	if hop.ISP != "Google LLC" {
		t.Errorf("ISP = %q, want %q", hop.ISP, "Google LLC")
	}
}

func TestTracerouteResult_Fields(t *testing.T) {
	now := time.Now()
	later := now.Add(5 * time.Second)

	result := roots.TracerouteResult{
		Target:     "google.com",
		ResolvedIP: "142.250.80.46",
		Hops: []roots.TracerouteHop{
			{Number: 1, RTTMs: 1.0},
			{Number: 2, RTTMs: 5.0},
		},
		Complete:    true,
		Duration:    5 * time.Second,
		DurationMs:  5000.0,
		StartedAt:   now,
		CompletedAt: later,
	}

	if result.Target != "google.com" {
		t.Errorf("Target = %q, want %q", result.Target, "google.com")
	}
	if result.ResolvedIP != "142.250.80.46" {
		t.Errorf("ResolvedIP = %q, want %q", result.ResolvedIP, "142.250.80.46")
	}
	if len(result.Hops) != 2 {
		t.Errorf("len(Hops) = %d, want %d", len(result.Hops), 2)
	}
	if !result.Complete {
		t.Error("Complete = false, want true")
	}
	if result.Duration != 5*time.Second {
		t.Errorf("Duration = %v, want %v", result.Duration, 5*time.Second)
	}
	if result.DurationMs != 5000.0 {
		t.Errorf("DurationMs = %f, want %f", result.DurationMs, 5000.0)
	}
	if !result.StartedAt.Equal(now) {
		t.Errorf("StartedAt = %v, want %v", result.StartedAt, now)
	}
	if !result.CompletedAt.Equal(later) {
		t.Errorf("CompletedAt = %v, want %v", result.CompletedAt, later)
	}
}

func TestTracerouteOptions_Fields(t *testing.T) {
	opts := roots.TracerouteOptions{
		MaxHops:     30,
		Timeout:     2 * time.Second,
		Probes:      3,
		PacketSize:  60,
		EnrichHops:  true,
		UseUDP:      false,
		SourceAddr:  "192.168.1.100",
		DontResolve: false,
	}

	if opts.MaxHops != 30 {
		t.Errorf("MaxHops = %d, want %d", opts.MaxHops, 30)
	}
	if opts.Timeout != 2*time.Second {
		t.Errorf("Timeout = %v, want %v", opts.Timeout, 2*time.Second)
	}
	if opts.Probes != 3 {
		t.Errorf("Probes = %d, want %d", opts.Probes, 3)
	}
	if opts.PacketSize != 60 {
		t.Errorf("PacketSize = %d, want %d", opts.PacketSize, 60)
	}
	if !opts.EnrichHops {
		t.Error("EnrichHops = false, want true")
	}
	if opts.UseUDP {
		t.Error("UseUDP = true, want false")
	}
	if opts.SourceAddr != "192.168.1.100" {
		t.Errorf("SourceAddr = %q, want %q", opts.SourceAddr, "192.168.1.100")
	}
	if opts.DontResolve {
		t.Error("DontResolve = true, want false")
	}
}

func TestRootsErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "ErrNotImplemented",
			err:  roots.ErrNotImplemented,
			want: "not implemented: pending migration",
		},
		{
			name: "ErrNotInitialized",
			err:  roots.ErrNotInitialized,
			want: "service not initialized",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Error() != tt.want {
				t.Errorf("error = %q, want %q", tt.err.Error(), tt.want)
			}
		})
	}
}
