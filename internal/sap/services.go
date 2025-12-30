package sap

import (
	"context"

	"github.com/krisarmstrong/seed/internal/config"
	"github.com/krisarmstrong/seed/internal/database"
)

// LinkService monitors network link status.
type LinkService struct {
	cfg    *config.Config
	cancel context.CancelFunc
}

// NewLinkService creates a new link service.
func NewLinkService(cfg *config.Config) *LinkService {
	return &LinkService{cfg: cfg}
}

// Start begins link monitoring.
func (s *LinkService) Start(ctx context.Context) error {
	ctx, s.cancel = context.WithCancel(ctx)
	// TODO: Migrate from internal/network + internal/phy
	_ = ctx
	return nil
}

// Stop halts link monitoring.
func (s *LinkService) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
}

// GetStatus returns current link status for all interfaces.
func (s *LinkService) GetStatus(_ context.Context) ([]LinkStatus, error) {
	// TODO: Migrate from internal/network
	return nil, ErrNotImplemented
}

// GetInterfaceStatus returns status for a specific interface.
func (s *LinkService) GetInterfaceStatus(_ context.Context, _ string) (*LinkStatus, error) {
	// TODO: Migrate from internal/network
	return nil, ErrNotImplemented
}

// CableService provides TDR cable testing.
type CableService struct {
	cfg *config.Config
}

// NewCableService creates a new cable service.
func NewCableService(cfg *config.Config) *CableService {
	return &CableService{cfg: cfg}
}

// Test performs a cable test on the specified interface.
func (s *CableService) Test(_ context.Context, _ string) (*CableTestResult, error) {
	// TODO: Migrate from internal/cable
	return nil, ErrNotImplemented
}

// DHCPService provides DHCP testing.
type DHCPService struct {
	cfg *config.Config
}

// NewDHCPService creates a new DHCP service.
func NewDHCPService(cfg *config.Config) *DHCPService {
	return &DHCPService{cfg: cfg}
}

// Test performs a DHCP discovery test.
func (s *DHCPService) Test(_ context.Context, _ string) (*DHCPTestResult, error) {
	// TODO: Migrate from internal/dhcp
	return nil, ErrNotImplemented
}

// DNSService provides DNS testing.
type DNSService struct {
	cfg *config.Config
}

// NewDNSService creates a new DNS service.
func NewDNSService(cfg *config.Config) *DNSService {
	return &DNSService{cfg: cfg}
}

// Test performs a DNS query test.
func (s *DNSService) Test(_ context.Context, _, _ string) (*DNSTestResult, error) {
	// TODO: Migrate from internal/dns
	return nil, ErrNotImplemented
}

// GatewayService monitors gateway health.
type GatewayService struct {
	cfg    *config.Config
	cancel context.CancelFunc
}

// NewGatewayService creates a new gateway service.
func NewGatewayService(cfg *config.Config) *GatewayService {
	return &GatewayService{cfg: cfg}
}

// Start begins gateway monitoring.
func (s *GatewayService) Start(ctx context.Context) error {
	ctx, s.cancel = context.WithCancel(ctx)
	// TODO: Migrate from internal/gateway
	_ = ctx
	return nil
}

// Stop halts gateway monitoring.
func (s *GatewayService) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
}

// GetHealth returns current gateway health.
func (s *GatewayService) GetHealth(_ context.Context) (*GatewayHealth, error) {
	// TODO: Migrate from internal/gateway
	return nil, ErrNotImplemented
}

// SNMPService collects SNMP data from network devices.
type SNMPService struct {
	cfg *config.Config
}

// NewSNMPService creates a new SNMP service.
func NewSNMPService(cfg *config.Config) *SNMPService {
	return &SNMPService{cfg: cfg}
}

// Collect gathers SNMP data from a device.
func (s *SNMPService) Collect(_ context.Context, _, _ string) (*SNMPDevice, error) {
	// TODO: Migrate from internal/snmp
	return nil, ErrNotImplemented
}

// GetInterfaces retrieves interface data via SNMP.
func (s *SNMPService) GetInterfaces(_ context.Context, _, _ string) ([]SNMPInterface, error) {
	// TODO: Migrate from internal/snmp
	return nil, ErrNotImplemented
}

// GetMACTable retrieves MAC address table via SNMP.
func (s *SNMPService) GetMACTable(_ context.Context, _, _ string) ([]MACTableEntry, error) {
	// TODO: Migrate from internal/snmp
	return nil, ErrNotImplemented
}

// PerformanceService provides speed and throughput testing.
type PerformanceService struct {
	cfg    *config.Config
	cancel context.CancelFunc
}

// NewPerformanceService creates a new performance service.
func NewPerformanceService(cfg *config.Config) *PerformanceService {
	return &PerformanceService{cfg: cfg}
}

// Stop halts any running tests.
func (s *PerformanceService) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
}

// Speedtest runs an internet speed test.
func (s *PerformanceService) Speedtest(_ context.Context) (*SpeedtestResult, error) {
	// TODO: Migrate from internal/speedtest
	return nil, ErrNotImplemented
}

// IPerf runs an iPerf test.
func (s *PerformanceService) IPerf(_ context.Context, _ string, _ map[string]any) (*IPerfResult, error) {
	// TODO: Migrate from internal/iperf
	return nil, ErrNotImplemented
}

// VLANService manages VLAN configuration.
type VLANService struct {
	cfg *config.Config
}

// NewVLANService creates a new VLAN service.
func NewVLANService(cfg *config.Config) *VLANService {
	return &VLANService{cfg: cfg}
}

// Create creates a VLAN subinterface.
func (s *VLANService) Create(_ context.Context, _ string, _ int) error {
	// TODO: Migrate from internal/vlan
	return ErrNotImplemented
}

// Delete removes a VLAN subinterface.
func (s *VLANService) Delete(_ context.Context, _ string) error {
	// TODO: Migrate from internal/vlan
	return ErrNotImplemented
}

// List returns all VLAN configurations.
func (s *VLANService) List(_ context.Context) ([]VLANConfig, error) {
	// TODO: Migrate from internal/vlan
	return nil, ErrNotImplemented
}

// TelemetryService aggregates and stores telemetry data.
type TelemetryService struct {
	cfg    *config.Config
	db     *database.DB
	cancel context.CancelFunc
}

// NewTelemetryService creates a new telemetry service.
func NewTelemetryService(cfg *config.Config, db *database.DB) *TelemetryService {
	return &TelemetryService{
		cfg: cfg,
		db:  db,
	}
}

// Start begins telemetry collection.
func (s *TelemetryService) Start(ctx context.Context) error {
	ctx, s.cancel = context.WithCancel(ctx)
	// TODO: Implement telemetry aggregation
	_ = ctx
	return nil
}

// Stop halts telemetry collection.
func (s *TelemetryService) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
}

// GetSnapshot returns the latest telemetry snapshot.
func (s *TelemetryService) GetSnapshot(_ context.Context) (*TelemetrySnapshot, error) {
	// TODO: Implement snapshot retrieval
	return nil, ErrNotImplemented
}

// GetHistory returns historical telemetry data.
func (s *TelemetryService) GetHistory(_ context.Context, _, _ string) ([]TelemetrySnapshot, error) {
	// TODO: Implement history retrieval
	return nil, ErrNotImplemented
}
