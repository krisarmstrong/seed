// Package survey provides WiFi site survey functionality.
package survey

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/krisarmstrong/seed/internal/iperf"
	"github.com/krisarmstrong/seed/internal/wifi"
)

// Type indicates the type of survey being conducted.
type Type string

// WiFi survey type constants.
const (
	TypePassive    Type = "passive"    // Passive scan (all visible networks)
	TypeActive     Type = "active"     // Active monitoring (current connection)
	TypeThroughput Type = "throughput" // Throughput testing with iperf3
)

// Status indicates the current status of a survey.
type Status string

// WiFi survey status constants.
const (
	StatusCreated    Status = "created"
	StatusInProgress Status = "in_progress"
	StatusPaused     Status = "paused"
	StatusCompleted  Status = "completed"
)

// FloorPlan contains floor plan image and metadata.
type FloorPlan struct {
	ImageData string  `json:"imageData"` // Base64-encoded image
	Width     int     `json:"width"`     // Image width in pixels
	Height    int     `json:"height"`    // Image height in pixels
	ScaleM    float64 `json:"scaleM"`    // Meters per pixel
}

// PassiveSample contains data from a passive WiFi scan.
type PassiveSample struct {
	Networks []*wifi.ScannedNetwork `json:"networks"` // All visible APs

	// Aggregated statistics for heatmap visualization
	UniqueSSIDs   int `json:"uniqueSSIDs"`   // Number of unique network names
	UniqueBSSIDs  int `json:"uniqueBSSIDs"`  // Number of unique access points
	APCount2_4    int `json:"apCount2_4"`    // Number of APs on 2.4 GHz band (2400-2500 MHz)
	APCount5      int `json:"apCount5"`      // Number of APs on 5 GHz band (5000-5900 MHz)
	APCount6      int `json:"apCount6"`      // Number of APs on 6 GHz band (5900+ MHz)
	CoChannelAPs  int `json:"coChannelAPs"`  // Number of APs on same channel as strongest AP
	AdjChannelAPs int `json:"adjChannelAPs"` // Number of APs on adjacent channels (±1-2 from strongest)
}

// CalculateAggregations computes aggregated statistics from the networks array.
// This should be called after populating the Networks field to generate heatmap data.
//
// The function calculates:
//   - Unique SSIDs and BSSIDs for network density metrics
//   - AP counts per frequency band (2.4 GHz, 5 GHz, 6 GHz) for band utilization
//   - Co-channel interference: APs on the same channel as the strongest AP
//   - Adjacent channel interference: APs on channels ±1 or ±2 from the strongest AP
//
// Handles nil or empty networks array gracefully.
func (p *PassiveSample) CalculateAggregations() {
	// Reset all aggregated fields
	p.UniqueSSIDs = 0
	p.UniqueBSSIDs = 0
	p.APCount2_4 = 0
	p.APCount5 = 0
	p.APCount6 = 0
	p.CoChannelAPs = 0
	p.AdjChannelAPs = 0

	// Handle nil or empty networks
	if len(p.Networks) == 0 {
		return
	}

	// Track unique SSIDs and BSSIDs using maps
	uniqueSSIDs := make(map[string]bool)
	uniqueBSSIDs := make(map[string]bool)

	// Find the strongest AP (first in the array, as they're sorted by signal strength)
	strongestChannel := p.Networks[0].Channel

	// Process each network
	for _, network := range p.Networks {
		// Count unique SSIDs (skip empty/hidden SSIDs)
		if network.SSID != "" {
			uniqueSSIDs[network.SSID] = true
		}

		// Count unique BSSIDs
		if network.BSSID != "" {
			uniqueBSSIDs[network.BSSID] = true
		}

		// Count APs by frequency band
		// 2.4 GHz: 2400-2500 MHz (channels 1-14)
		// 5 GHz: 5000-5900 MHz (channels 36-165)
		// 6 GHz: 5900+ MHz (channels 1-233 in 6GHz band)
		switch {
		case network.Frequency >= 2400 && network.Frequency < 2500:
			p.APCount2_4++
		case network.Frequency >= 5000 && network.Frequency < 5900:
			p.APCount5++
		case network.Frequency >= 5900:
			p.APCount6++
		}

		// Count co-channel interference (same channel as strongest)
		if network.Channel == strongestChannel {
			p.CoChannelAPs++
		}

		// Count adjacent channel interference (±1 or ±2 channels from strongest)
		// For 2.4 GHz, channels 1-11 are commonly used with ±2 overlap
		// For 5 GHz, channels are typically spaced further apart
		channelDiff := network.Channel - strongestChannel
		if channelDiff < 0 {
			channelDiff = -channelDiff
		}
		if channelDiff >= 1 && channelDiff <= 2 {
			p.AdjChannelAPs++
		}
	}

	// Set the counts
	p.UniqueSSIDs = len(uniqueSSIDs)
	p.UniqueBSSIDs = len(uniqueBSSIDs)
}

// ActiveSample contains data from active connection monitoring.
type ActiveSample struct {
	SSID          string  `json:"ssid"`
	BSSID         string  `json:"bssid"`
	RSSI          int     `json:"rssi"`                    // Signal strength in dBm
	DataRate      float64 `json:"dataRate"`                // Mbps
	RoamingEvent  bool    `json:"roamingEvent"`            // true if BSSID changed since last sample
	PreviousBSSID string  `json:"previousBssid,omitempty"` // BSSID before roaming event
	RoamCount     int     `json:"roamCount,omitempty"`     // Total number of roaming events during survey
}

// ThroughputSample contains data from iperf3 throughput testing.
type ThroughputSample struct {
	SSID         string  `json:"ssid"`
	BSSID        string  `json:"bssid"`
	RSSI         int     `json:"rssi"`
	DownloadMbps float64 `json:"downloadMbps"`
	UploadMbps   float64 `json:"uploadMbps"`
	Latency      float64 `json:"latency"`    // milliseconds
	Jitter       float64 `json:"jitter"`     // milliseconds
	PacketLoss   float64 `json:"packetLoss"` // percentage
}

// SamplePoint represents a measurement at a specific location.
type SamplePoint struct {
	X          int         `json:"x"` // Pixel X coordinate on floor plan
	Y          int         `json:"y"` // Pixel Y coordinate on floor plan
	Timestamp  time.Time   `json:"timestamp"`
	SampleData interface{} `json:"sampleData"` // PassiveSample | ActiveSample | ThroughputSample
}

// Survey represents a WiFi site survey.
type Survey struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	FloorPlan   *FloorPlan     `json:"floorPlan,omitempty"`
	SurveyType  Type           `json:"surveyType"`
	Status      Status         `json:"status"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	Samples     []*SamplePoint `json:"samples"`

	// Configuration
	Interface    string `json:"interface"`              // WiFi interface to use
	IperfServer  string `json:"iperfServer,omitempty"`  // For throughput surveys
	TestDuration int    `json:"testDuration,omitempty"` // seconds, for throughput tests
}

// Manager manages WiFi site surveys.
type Manager struct {
	mu           sync.RWMutex
	surveys      map[string]*Survey // key is survey ID
	storagePath  string
	wifiScanner  *wifi.Scanner
	wifiManager  *wifi.Manager
	iperfManager *iperf.Manager
}

// NewManager creates a new survey manager.
func NewManager(storagePath string, wifiScanner *wifi.Scanner, wifiManager *wifi.Manager, iperfManager *iperf.Manager) *Manager {
	return &Manager{
		surveys:      make(map[string]*Survey),
		storagePath:  storagePath,
		wifiScanner:  wifiScanner,
		wifiManager:  wifiManager,
		iperfManager: iperfManager,
	}
}

// CreateSurvey creates a new survey.
func (m *Manager) CreateSurvey(name, description, iface string, surveyType Type) (*Survey, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	survey := &Survey{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		SurveyType:  surveyType,
		Status:      StatusCreated,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Interface:   iface,
		Samples:     make([]*SamplePoint, 0),
	}

	m.surveys[survey.ID] = survey

	// Persist to disk
	if err := m.saveSurvey(survey); err != nil {
		return nil, fmt.Errorf("failed to save survey: %w", err)
	}

	return survey, nil
}

// GetSurvey retrieves a survey by ID.
func (m *Manager) GetSurvey(id string) (*Survey, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	survey, exists := m.surveys[id]
	if !exists {
		return nil, fmt.Errorf("survey not found: %s", id)
	}

	return survey, nil
}

// ListSurveys returns all surveys.
func (m *Manager) ListSurveys() []*Survey {
	m.mu.RLock()
	defer m.mu.RUnlock()

	surveys := make([]*Survey, 0, len(m.surveys))
	for _, survey := range m.surveys {
		surveys = append(surveys, survey)
	}

	return surveys
}

// DeleteSurvey deletes a survey.
func (m *Manager) DeleteSurvey(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.surveys[id]; !exists {
		return fmt.Errorf("survey not found: %s", id)
	}

	// Delete from disk
	filename := filepath.Join(m.storagePath, fmt.Sprintf("%s.json", id))
	if err := os.Remove(filename); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete survey file: %w", err)
	}

	delete(m.surveys, id)

	return nil
}

// UpdateFloorPlan updates the floor plan for a survey.
func (m *Manager) UpdateFloorPlan(id string, floorPlan *FloorPlan) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	survey, exists := m.surveys[id]
	if !exists {
		return fmt.Errorf("survey not found: %s", id)
	}

	survey.FloorPlan = floorPlan
	survey.UpdatedAt = time.Now()

	return m.saveSurvey(survey)
}

// StartSurvey starts a survey.
func (m *Manager) StartSurvey(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	survey, exists := m.surveys[id]
	if !exists {
		return fmt.Errorf("survey not found: %s", id)
	}

	if survey.Status == StatusInProgress {
		return fmt.Errorf("survey already in progress")
	}

	survey.Status = StatusInProgress
	survey.UpdatedAt = time.Now()

	return m.saveSurvey(survey)
}

// PauseSurvey pauses a survey.
func (m *Manager) PauseSurvey(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	survey, exists := m.surveys[id]
	if !exists {
		return fmt.Errorf("survey not found: %s", id)
	}

	survey.Status = StatusPaused
	survey.UpdatedAt = time.Now()

	return m.saveSurvey(survey)
}

// CompleteSurvey completes a survey.
func (m *Manager) CompleteSurvey(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	survey, exists := m.surveys[id]
	if !exists {
		return fmt.Errorf("survey not found: %s", id)
	}

	survey.Status = StatusCompleted
	survey.UpdatedAt = time.Now()

	return m.saveSurvey(survey)
}

// UpdateSurveySettings updates survey settings (only when survey is in created state).
func (m *Manager) UpdateSurveySettings(id string, surveyType Type, iperfServer string, testDuration int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	survey, exists := m.surveys[id]
	if !exists {
		return fmt.Errorf("survey not found: %s", id)
	}

	if survey.Status != StatusCreated {
		return fmt.Errorf("cannot update settings after survey has started")
	}

	// Validate survey type
	switch surveyType {
	case TypePassive, TypeActive, TypeThroughput:
		survey.SurveyType = surveyType
	default:
		return fmt.Errorf("invalid survey type: %s", surveyType)
	}

	// Validate test duration
	if testDuration < 1 {
		testDuration = 3 // Default
	}
	if testDuration > 60 {
		testDuration = 60 // Max
	}
	survey.TestDuration = testDuration
	survey.IperfServer = iperfServer
	survey.UpdatedAt = time.Now()

	return m.saveSurvey(survey)
}

// AddSample adds a measurement sample to a survey.
func (m *Manager) AddSample(id string, x, y int, sampleData interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	survey, exists := m.surveys[id]
	if !exists {
		return fmt.Errorf("survey not found: %s", id)
	}

	if survey.Status != StatusInProgress {
		return fmt.Errorf("survey not in progress")
	}

	sample := &SamplePoint{
		X:          x,
		Y:          y,
		Timestamp:  time.Now(),
		SampleData: sampleData,
	}

	survey.Samples = append(survey.Samples, sample)
	survey.UpdatedAt = time.Now()

	return m.saveSurvey(survey)
}

// saveSurvey persists a survey to disk.
func (m *Manager) saveSurvey(survey *Survey) error {
	// Ensure storage directory exists with restrictive permissions
	if err := os.MkdirAll(m.storagePath, 0o750); err != nil {
		return fmt.Errorf("failed to create storage directory: %w", err)
	}

	filename := filepath.Join(m.storagePath, fmt.Sprintf("%s.json", survey.ID))
	data, err := json.MarshalIndent(survey, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal survey: %w", err)
	}

	if err := os.WriteFile(filename, data, 0o600); err != nil {
		return fmt.Errorf("failed to write survey file: %w", err)
	}

	return nil
}

// LoadSurveys loads all surveys from disk.
func (m *Manager) LoadSurveys() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Create storage directory if it doesn't exist with restrictive permissions
	if err := os.MkdirAll(m.storagePath, 0o750); err != nil {
		return fmt.Errorf("failed to create storage directory: %w", err)
	}

	// Read all JSON files
	files, err := filepath.Glob(filepath.Join(m.storagePath, "*.json"))
	if err != nil {
		return fmt.Errorf("failed to list survey files: %w", err)
	}

	for _, file := range files {
		// file comes from filepath.Glob with our controlled pattern, not user input
		data, err := os.ReadFile(file) //nolint:gosec // G304: file path from Glob with controlled pattern
		if err != nil {
			continue // Skip files that can't be read
		}

		var survey Survey
		if err := json.Unmarshal(data, &survey); err != nil {
			continue // Skip invalid JSON
		}

		m.surveys[survey.ID] = &survey
	}

	return nil
}
