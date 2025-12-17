// Package api provides the HTTP/WebSocket server.
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/krisarmstrong/seed/internal/config"
	"github.com/krisarmstrong/seed/internal/logging"
)

// ============================================================================
// Device Discovery Handlers (fixes #544 - split from handlers_discovery.go)
// ============================================================================

// handleDevices returns discovered devices and status.
func (s *Server) handleDevices(w http.ResponseWriter, r *http.Request) {
	logger := logging.FromContext(r.Context())
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if s.deviceDiscovery == nil {
		http.Error(w, "Device discovery not available", http.StatusServiceUnavailable)
		return
	}

	devices := s.deviceDiscovery.GetDevices()
	status := s.deviceDiscovery.GetStatus()

	resp := map[string]interface{}{
		"devices": devices,
		"status":  status,
	}

	sendJSONResponse(w, logger, http.StatusOK, resp)
}

// handleDevicesScan triggers a network device scan.
func (s *Server) handleDevicesScan(w http.ResponseWriter, r *http.Request) {
	logger := logging.FromContext(r.Context())
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if s.deviceDiscovery == nil {
		http.Error(w, "Device discovery not available", http.StatusServiceUnavailable)
		return
	}

	// Check if scan is already in progress
	if s.deviceDiscovery.IsScanning() {
		sendJSONResponse(w, logger, http.StatusOK, map[string]interface{}{
			"message":  "Scan already in progress",
			"scanning": true,
		})
		return
	}

	// Start scan in background
	go func(reqCtx context.Context) {
		logger := logging.FromContext(reqCtx)
		ctx, cancel := context.WithTimeout(reqCtx, 2*time.Minute)
		defer cancel()

		logger.Info("Starting background device scan")
		start := time.Now()
		defer func() {
			logger.Info("Background device scan finished", "duration_ms", time.Since(start).Milliseconds())
		}()

		if err := s.deviceDiscovery.Scan(ctx); err != nil {
			logger.Error("Device scan error", "error", err)
		}

		// Auto-scan for vulnerabilities if enabled
		s.postScanVulnerabilityCheck(logger)

		// Notify WebSocket clients when scan completes
		s.wsHub.Broadcast(Message{
			Type: "deviceScanComplete",
			Payload: map[string]interface{}{
				"deviceCount": s.deviceDiscovery.Count(),
				"timestamp":   time.Now().Format(time.RFC3339),
			},
		})
	}(r.Context())

	sendJSONResponse(w, logger, http.StatusOK, map[string]interface{}{
		"message":  "Scan started",
		"scanning": true,
	})
}

// postScanVulnerabilityCheck runs vulnerability scans after device discovery if auto-scan is enabled.
// This method extracts business logic from the handler for better separation of concerns.
func (s *Server) postScanVulnerabilityCheck(logger *slog.Logger) {
	if s.vulnScanner == nil || !s.config.Security.VulnerabilityScanning.Enabled ||
		!s.config.Security.VulnerabilityScanning.AutoScan {
		return
	}

	logger.Info("Auto-scan: triggering vulnerability scan", "device_count", s.deviceDiscovery.Count())
	devices := s.deviceDiscovery.GetDevices()

	vulnCtx, vulnCancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer vulnCancel()

	for _, device := range devices {
		if _, err := s.vulnScanner.ScanDevice(vulnCtx, device); err != nil {
			logger.Warn("Auto vulnerability scan failed", "device_ip", device.IP, "error", err)
		}
	}

	// Broadcast vulnerability results
	results := s.vulnScanner.GetAllVulnerabilities()
	s.wsHub.BroadcastCardUpdate("vulnerabilities", map[string]interface{}{
		"results": results,
		"count":   len(results),
	})
	logger.Info("Auto-scan: completed vulnerability scan", "vulnerable_devices", len(results))
}

// handleDevicesStatus returns the current device discovery status.
func (s *Server) handleDevicesStatus(w http.ResponseWriter, r *http.Request) {
	logger := logging.FromContext(r.Context())
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if s.deviceDiscovery == nil {
		http.Error(w, "Device discovery not available", http.StatusServiceUnavailable)
		return
	}

	status := s.deviceDiscovery.GetStatus()
	sendJSONResponse(w, logger, http.StatusOK, status)
}

// NetworkDiscoverySettingsResponse represents network discovery settings.
type NetworkDiscoverySettingsResponse struct {
	Enabled        bool   `json:"enabled"`
	ARPScanWorkers int    `json:"arpScanWorkers"`
	PingTimeoutMs  int64  `json:"pingTimeoutMs"`
	ScanTimeoutMs  int64  `json:"scanTimeoutMs"`
	AutoScan       bool   `json:"autoScan"`
	ScanIntervalMs int64  `json:"scanIntervalMs"`
	OUIFilePath    string `json:"ouiFilePath"`
}

// handleDevicesSettings handles GET/PUT for network discovery settings.
func (s *Server) handleDevicesSettings(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.getDevicesSettings(w, r)
	case http.MethodPut:
		s.updateDevicesSettings(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) getDevicesSettings(w http.ResponseWriter, r *http.Request) {
	logger := logging.FromContext(r.Context())
	resp := NetworkDiscoverySettingsResponse{
		Enabled:        s.config.NetworkDiscovery.Enabled,
		ARPScanWorkers: s.config.NetworkDiscovery.ARPScanWorkers,
		PingTimeoutMs:  s.config.NetworkDiscovery.PingTimeout.Milliseconds(),
		ScanTimeoutMs:  s.config.NetworkDiscovery.ScanTimeout.Milliseconds(),
		AutoScan:       s.config.NetworkDiscovery.AutoScan,
		ScanIntervalMs: s.config.NetworkDiscovery.ScanInterval.Milliseconds(),
		OUIFilePath:    s.config.NetworkDiscovery.OUIFilePath,
	}

	sendJSONResponse(w, logger, http.StatusOK, resp)
}

func (s *Server) updateDevicesSettings(w http.ResponseWriter, r *http.Request) {
	logger := logging.FromContext(r.Context())
	// Limit request body size to prevent DoS attacks (fixes #693)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodySizeJSON)

	var req NetworkDiscoverySettingsResponse
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update config
	s.config.NetworkDiscovery.Enabled = req.Enabled
	if req.ARPScanWorkers > 0 {
		s.config.NetworkDiscovery.ARPScanWorkers = req.ARPScanWorkers
	}
	if req.PingTimeoutMs > 0 {
		s.config.NetworkDiscovery.PingTimeout = time.Duration(req.PingTimeoutMs) * time.Millisecond
	}
	if req.ScanTimeoutMs > 0 {
		s.config.NetworkDiscovery.ScanTimeout = time.Duration(req.ScanTimeoutMs) * time.Millisecond
	}
	s.config.NetworkDiscovery.AutoScan = req.AutoScan
	s.config.NetworkDiscovery.ScanInterval = time.Duration(req.ScanIntervalMs) * time.Millisecond
	if req.OUIFilePath != "" {
		s.config.NetworkDiscovery.OUIFilePath = req.OUIFilePath
	}

	// Save config to file
	if err := s.config.Save(s.configPath); err != nil {
		logger.Warn("Failed to save config", "error", err)
	}

	sendJSONResponse(w, logger, http.StatusOK, map[string]string{
		"status":  "success",
		"message": "Network discovery settings updated",
	})
}

// SubnetRequest represents a subnet configuration request.
type SubnetRequest struct {
	CIDR    string `json:"cidr"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

// SubnetResponse represents a subnet in API responses.
type SubnetResponse struct {
	CIDR    string `json:"cidr"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

// handleDevicesSubnets handles GET/POST/DELETE for additional subnets.
func (s *Server) handleDevicesSubnets(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.getDevicesSubnets(w, r)
	case http.MethodPost:
		s.addDevicesSubnet(w, r)
	case http.MethodPut:
		s.updateDevicesSubnet(w, r)
	case http.MethodDelete:
		s.deleteDevicesSubnet(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) getDevicesSubnets(w http.ResponseWriter, r *http.Request) {
	logger := logging.FromContext(r.Context())
	subnets := make([]SubnetResponse, 0, len(s.config.NetworkDiscovery.AdditionalSubnets))
	for _, subnet := range s.config.NetworkDiscovery.AdditionalSubnets {
		subnets = append(subnets, SubnetResponse{
			CIDR:    subnet.CIDR,
			Name:    subnet.Name,
			Enabled: subnet.Enabled,
		})
	}

	sendJSONResponse(w, logger, http.StatusOK, subnets)
}

func (s *Server) addDevicesSubnet(w http.ResponseWriter, r *http.Request) {
	logger := logging.FromContext(r.Context())
	// Limit request body size to prevent DoS attacks (fixes #693)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodySizeJSON)

	var req SubnetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate CIDR format
	_, _, err := net.ParseCIDR(req.CIDR)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid CIDR format: %v", err), http.StatusBadRequest)
		return
	}

	// Check for duplicates
	for _, existing := range s.config.NetworkDiscovery.AdditionalSubnets {
		if existing.CIDR == req.CIDR {
			http.Error(w, "Subnet already exists", http.StatusConflict)
			return
		}
	}

	// Add the new subnet
	newSubnet := config.SubnetConfig{
		CIDR:    req.CIDR,
		Name:    req.Name,
		Enabled: req.Enabled,
	}
	s.config.NetworkDiscovery.AdditionalSubnets = append(
		s.config.NetworkDiscovery.AdditionalSubnets,
		newSubnet,
	)

	// Update the device discovery scanner
	s.syncDeviceDiscoverySubnets(logger)

	// Save config to file
	if err := s.config.Save(s.configPath); err != nil {
		logger.Warn("Failed to save config", "error", err)
	}

	sendJSONResponse(w, logger, http.StatusOK, map[string]string{
		"status":  "success",
		"message": "Subnet added",
	})
}

func (s *Server) updateDevicesSubnet(w http.ResponseWriter, r *http.Request) {
	logger := logging.FromContext(r.Context())
	// Limit request body size to prevent DoS attacks (fixes #693)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodySizeJSON)

	var req SubnetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate CIDR format
	if _, _, err := net.ParseCIDR(req.CIDR); err != nil {
		http.Error(w, fmt.Sprintf("Invalid CIDR format: %v", err), http.StatusBadRequest)
		return
	}

	// Find and update the subnet
	found := false
	for i, existing := range s.config.NetworkDiscovery.AdditionalSubnets {
		if existing.CIDR == req.CIDR {
			s.config.NetworkDiscovery.AdditionalSubnets[i].Name = req.Name
			s.config.NetworkDiscovery.AdditionalSubnets[i].Enabled = req.Enabled
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Subnet not found", http.StatusNotFound)
		return
	}

	// Update the device discovery scanner
	s.syncDeviceDiscoverySubnets(logger)

	// Save config to file
	if err := s.config.Save(s.configPath); err != nil {
		logger.Warn("Failed to save config", "error", err)
	}

	sendJSONResponse(w, logger, http.StatusOK, map[string]string{
		"status":  "success",
		"message": "Subnet updated",
	})
}

func (s *Server) deleteDevicesSubnet(w http.ResponseWriter, r *http.Request) {
	logger := logging.FromContext(r.Context())
	cidr := r.URL.Query().Get("cidr")
	if cidr == "" {
		http.Error(w, "CIDR parameter required", http.StatusBadRequest)
		return
	}

	// Find and remove the subnet
	found := false
	newSubnets := make([]config.SubnetConfig, 0, len(s.config.NetworkDiscovery.AdditionalSubnets))
	for _, existing := range s.config.NetworkDiscovery.AdditionalSubnets {
		if existing.CIDR == cidr {
			found = true
			continue
		}
		newSubnets = append(newSubnets, existing)
	}

	if !found {
		http.Error(w, "Subnet not found", http.StatusNotFound)
		return
	}

	s.config.NetworkDiscovery.AdditionalSubnets = newSubnets

	// Update the device discovery scanner
	s.syncDeviceDiscoverySubnets(logger)

	// Save config to file
	if err := s.config.Save(s.configPath); err != nil {
		logger.Warn("Failed to save config", "error", err)
	}

	sendJSONResponse(w, logger, http.StatusOK, map[string]string{
		"status":  "success",
		"message": "Subnet deleted",
	})
}

// syncDeviceDiscoverySubnets synchronizes enabled subnets from config to the device discovery scanner.
// This helper method eliminates DRY violation across add/update/delete subnet handlers.
func (s *Server) syncDeviceDiscoverySubnets(logger *slog.Logger) {
	if s.deviceDiscovery == nil {
		return
	}

	enabledCIDRs := make([]string, 0, len(s.config.NetworkDiscovery.AdditionalSubnets))
	for _, subnet := range s.config.NetworkDiscovery.AdditionalSubnets {
		if subnet.Enabled {
			enabledCIDRs = append(enabledCIDRs, subnet.CIDR)
		}
	}

	if err := s.deviceDiscovery.SetAdditionalSubnets(enabledCIDRs); err != nil {
		logger.Warn("Failed to update scanner subnets", "error", err)
	}
}

// handlePublicIP returns the public IPv4 and IPv6 addresses.
func (s *Server) handlePublicIP(w http.ResponseWriter, r *http.Request) {
	logger := logging.FromContext(r.Context())
	if s.publicipChecker == nil {
		http.Error(w, "Public IP checker not available", http.StatusServiceUnavailable)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Return cached result or fetch if cache expired
		result := s.publicipChecker.GetPublicIP(r.Context())
		sendJSONResponse(w, logger, http.StatusOK, result)

	case http.MethodPost:
		// Force refresh
		result := s.publicipChecker.Refresh(r.Context())
		sendJSONResponse(w, logger, http.StatusOK, result)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
