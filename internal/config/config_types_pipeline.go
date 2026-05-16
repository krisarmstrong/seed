package config

// config_types_pipeline.go contains the PipelineConfig type tree (phases,
// timing, port-scan, SNMP collection, persistence) and the read-only
// adapter getters the discovery package uses.

import "time"

// PipelineConfig controls the sequential discovery pipeline.
type PipelineConfig struct {
	// Phases controls which pipeline phases are enabled.
	Phases PipelinePhaseConfig `json:"phases"`

	// Timing controls rate limiting and delays.
	Timing PipelineTimingConfig `json:"timing"`

	// PortScan controls port scanning behavior and intensity.
	PortScan PipelinePortScanConfig `json:"port_scan"`

	// SNMPCollection controls extended SNMP MIB collection.
	SNMPCollection PipelineSNMPConfig `json:"snmp_collection"`

	// Persistence controls how results are stored.
	Persistence PipelinePersistenceConfig `json:"persistence"`
}

// PipelinePhaseConfig controls which phases are executed.
type PipelinePhaseConfig struct {
	Enumeration      bool `json:"enumeration"`       // Always true - core functionality
	NameResolution   bool `json:"name_resolution"`   // Default: true
	ServiceDiscovery bool `json:"service_discovery"` // Default: false (passive only)
	VulnAssessment   bool `json:"vuln_assessment"`   // Default: false
}

// PipelineTimingConfig controls scan rate limiting.
type PipelineTimingConfig struct {
	// ProbeDelay is the minimum time between probes to a single host.
	ProbeDelay time.Duration `json:"probe_delay"`

	// HostDelay is the minimum time between starting scans of different hosts.
	HostDelay time.Duration `json:"host_delay"`

	// MaxConcurrentHosts limits parallel host scanning.
	MaxConcurrentHosts int `json:"max_concurrent_hosts"`

	// PhaseTimeout is the max duration for any single phase.
	PhaseTimeout time.Duration `json:"phase_timeout"`

	// Profile selects a pre-defined timing profile: polite, normal, aggressive.
	Profile string `json:"profile"`
}

// PipelinePortScanConfig controls port scanning intensity.
type PipelinePortScanConfig struct {
	// Intensity controls which ports are scanned: off, quick, standard, comprehensive, custom.
	Intensity string `json:"intensity"`

	// CustomPorts for Intensity="custom".
	CustomPorts []int `json:"custom_ports,omitempty"`

	// BannerGrab enables service banner reading.
	BannerGrab bool `json:"banner_grab"`

	// ConnectTimeout for port connections.
	ConnectTimeout time.Duration `json:"connect_timeout"`
}

// PipelineSNMPConfig controls extended SNMP data collection.
type PipelineSNMPConfig struct {
	// Enabled turns on extended SNMP collection in Phase 3.
	Enabled bool `json:"enabled"`

	// MIBs specifies which MIB groups to collect.
	MIBs PipelineSNMPMIBs `json:"mibs"`

	// WalkTimeout per MIB walk operation.
	WalkTimeout time.Duration `json:"walk_timeout"`

	// MaxOIDsPerRequest for bulk requests.
	MaxOIDsPerRequest int `json:"max_oids_per_request"`
}

// PipelineSNMPMIBs controls which MIBs are collected.
type PipelineSNMPMIBs struct {
	System      bool `json:"system"`       // SNMPv2-MIB::system (always on)
	Interfaces  bool `json:"interfaces"`   // IF-MIB (ifTable, ifXTable)
	IPAddresses bool `json:"ip_addresses"` // IP-MIB (ipAddrTable)
	Routing     bool `json:"routing"`      // IP-FORWARD-MIB
	Bridge      bool `json:"bridge"`       // BRIDGE-MIB (MAC table)
	Entity      bool `json:"entity"`       // ENTITY-MIB (physical inventory)
	LLDP        bool `json:"lldp"`         // LLDP-MIB
	VLAN        bool `json:"vlan"`         // Q-BRIDGE-MIB
}

// PipelinePersistenceConfig controls database storage.
type PipelinePersistenceConfig struct {
	// StoreHistory keeps historical device state.
	StoreHistory bool `json:"store_history"`

	// StalenessThreshold marks devices inactive after this duration.
	StalenessThreshold time.Duration `json:"staleness_threshold"`

	// PurgeAfter removes inactive devices after this duration.
	PurgeAfter time.Duration `json:"purge_after"`
}

// GetPhases implements discovery.ConfigPipelineAdapter.
func (c *PipelineConfig) GetPhases() (bool, bool, bool, bool) {
	return c.Phases.Enumeration, c.Phases.NameResolution, c.Phases.ServiceDiscovery, c.Phases.VulnAssessment
}

// GetTiming implements discovery.ConfigPipelineAdapter.
func (c *PipelineConfig) GetTiming() (time.Duration, time.Duration, time.Duration, int, string) {
	return c.Timing.ProbeDelay, c.Timing.HostDelay, c.Timing.PhaseTimeout, c.Timing.MaxConcurrentHosts, c.Timing.Profile
}

// GetPortScan implements discovery.ConfigPipelineAdapter.
func (c *PipelineConfig) GetPortScan() (string, []int, bool, time.Duration) {
	// Fixes #959: Deep copy CustomPorts to prevent caller mutation
	var portsCopy []int
	if len(c.PortScan.CustomPorts) > 0 {
		portsCopy = make([]int, len(c.PortScan.CustomPorts))
		copy(portsCopy, c.PortScan.CustomPorts)
	}
	return c.PortScan.Intensity, portsCopy, c.PortScan.BannerGrab, c.PortScan.ConnectTimeout
}

// GetSNMP implements discovery.ConfigPipelineAdapter.
//

func (c *PipelineConfig) GetSNMP() (bool, bool, bool, bool, bool, bool, bool, bool, bool, time.Duration, int) {
	return c.SNMPCollection.Enabled,
		c.SNMPCollection.MIBs.System,
		c.SNMPCollection.MIBs.Interfaces,
		c.SNMPCollection.MIBs.IPAddresses,
		c.SNMPCollection.MIBs.Routing,
		c.SNMPCollection.MIBs.Bridge,
		c.SNMPCollection.MIBs.Entity,
		c.SNMPCollection.MIBs.LLDP,
		c.SNMPCollection.MIBs.VLAN,
		c.SNMPCollection.WalkTimeout,
		c.SNMPCollection.MaxOIDsPerRequest
}

// GetPersistence implements discovery.ConfigPipelineAdapter.
func (c *PipelineConfig) GetPersistence() (bool, time.Duration, time.Duration) {
	return c.Persistence.StoreHistory,
		c.Persistence.StalenessThreshold,
		c.Persistence.PurgeAfter
}
