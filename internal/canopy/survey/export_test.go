// Package survey exports internal functions for testing.
package survey

import (
	"image"
)

// InterpolateColor exports interpolateColor for testing.
var InterpolateColor = interpolateColor

// FilterWeakSamples exports filterWeakSamples for testing.
var FilterWeakSamples = filterWeakSamples

// CalculateCoverageScore exports calculateCoverageScore for testing.
var CalculateCoverageScore = calculateCoverageScore

// DetermineSeverity exports determineSeverity for testing.
var DetermineSeverity = determineSeverity

// GenerateRecommendations exports generateRecommendations for testing.
var GenerateRecommendations = generateRecommendations

// SetSurvey sets a survey in the manager for testing.
func (m *Manager) SetSurvey(s *Survey) {
	m.surveys[s.ID] = s
}

// GetStoragePath returns the storage path for testing.
func (m *Manager) GetStoragePath() string {
	return m.storagePath
}

// GetSurveys returns the surveys map for testing.
func (m *Manager) GetSurveys() map[string]*Survey {
	return m.surveys
}

// GetWifiScanner returns the WiFi scanner for testing.
func (m *Manager) GetWifiScanner() any {
	return m.wifiScanner
}

// GetWifiManager returns the WiFi manager for testing.
func (m *Manager) GetWifiManager() any {
	return m.wifiManager
}

// GetIperfManager returns the iperf manager for testing.
func (m *Manager) GetIperfManager() any {
	return m.iperfManager
}

// Distance exports distance for testing.
var Distance = distance

// ExtractValue exports extractValue for testing.
var ExtractValue = extractValue

// ExtractPassiveValue exports extractPassiveValue for testing.
var ExtractPassiveValue = extractPassiveValue

// ExtractActiveValue exports extractActiveValue for testing.
var ExtractActiveValue = extractActiveValue

// ExtractThroughputValue exports extractThroughputValue for testing.
var ExtractThroughputValue = extractThroughputValue

// ExtractMapValue exports extractMapValue for testing.
var ExtractMapValue = extractMapValue

// GetHeatmapDimensions exports getHeatmapDimensions for testing.
var GetHeatmapDimensions = getHeatmapDimensions

// MapHeatmapTypeToValueType exports mapHeatmapTypeToValueType for testing.
var MapHeatmapTypeToValueType = mapHeatmapTypeToValueType

// GetColorScaleForType exports getColorScaleForType for testing.
var GetColorScaleForType = getColorScaleForType

// RenderHeatmapToImage exports renderHeatmapToImage for testing.
var RenderHeatmapToImage = renderHeatmapToImage

// RenderSamplePoints exports renderSamplePoints for testing.
var RenderSamplePoints = renderSamplePoints

// RenderGrid exports renderGrid for testing.
var RenderGrid = renderGrid

// CreateTestImage creates a test image for testing.
func CreateTestImage(width, height int) *image.RGBA {
	return image.NewRGBA(image.Rect(0, 0, width, height))
}
