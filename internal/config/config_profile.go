package config

// config_profile.go contains the per-profile JSON export/import contract.
// The ProfileExportFields struct is the single source of truth for which
// fields are user-portable between deployments; global settings (Server,
// Auth, Security, Logging, Database) are deliberately excluded.

import (
	"encoding/json"
	"fmt"
)

// ProfileExportFields defines which Config sections are stored per-profile.
// These are the user-configurable settings that vary between deployment profiles.
// Global settings (Server, Auth, Security, Logging, Database) are NOT included.
type ProfileExportFields struct {
	Thresholds       ThresholdsConfig       `json:"thresholds"`
	HealthChecks     HealthChecksConfig     `json:"healthChecks"`
	Speedtest        SpeedtestConfig        `json:"speedtest"`
	Iperf            IperfConfig            `json:"iperf"`
	FABOptions       FABOptionsConfig       `json:"fabOptions"`
	DisplayOptions   DisplayOptionsConfig   `json:"displayOptions"`
	DNS              DNSConfig              `json:"dns"`
	SNMP             SNMPConfig             `json:"snmp"`
	NetworkDiscovery NetworkDiscoveryConfig `json:"networkDiscovery"`
	Link             LinkConfig             `json:"link,omitzero"`
	CableTest        CableTestConfig        `json:"cableTest,omitzero"`
}

// ToProfileJSON exports profile-specific settings as JSON.
// Only settings that vary per-profile are included.
// Global settings (Server, Auth, Security) are excluded.
func (c *Config) ToProfileJSON() (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	export := ProfileExportFields{
		Thresholds:       c.Thresholds,
		HealthChecks:     c.HealthChecks,
		Speedtest:        c.Speedtest,
		Iperf:            c.Iperf,
		FABOptions:       c.FABOptions,
		DisplayOptions:   c.DisplayOptions,
		DNS:              c.DNS,
		SNMP:             c.SNMP,
		NetworkDiscovery: c.NetworkDiscovery,
		Link:             c.Link,
		CableTest:        c.CableTest,
	}

	data, err := json.Marshal(export)
	if err != nil {
		return "", fmt.Errorf("marshal profile settings: %w", err)
	}
	return string(data), nil
}

// ApplyProfileJSON applies profile settings from JSON to this config.
// Only profile-specific fields are updated; global settings are preserved.
func (c *Config) ApplyProfileJSON(jsonStr string) error {
	if jsonStr == "" {
		return nil
	}

	var imported ProfileExportFields
	if err := json.Unmarshal([]byte(jsonStr), &imported); err != nil {
		return fmt.Errorf("unmarshal profile settings: %w", err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.Thresholds = imported.Thresholds
	c.HealthChecks = imported.HealthChecks
	c.Speedtest = imported.Speedtest
	c.Iperf = imported.Iperf
	c.FABOptions = imported.FABOptions
	c.DisplayOptions = imported.DisplayOptions
	c.DNS = imported.DNS
	c.SNMP = imported.SNMP
	c.NetworkDiscovery = imported.NetworkDiscovery
	c.Link = imported.Link
	c.CableTest = imported.CableTest

	return nil
}
