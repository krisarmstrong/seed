// Package network handles network interface management.
package network

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
)

// InterfaceType represents the type of network interface.
type InterfaceType string

const (
	InterfaceTypeEthernet InterfaceType = "ethernet"
	InterfaceTypeWiFi     InterfaceType = "wifi"
	InterfaceTypeLoopback InterfaceType = "loopback"
	InterfaceTypeOther    InterfaceType = "other"
)

// InterfaceInfo contains information about a network interface.
type InterfaceInfo struct {
	Name         string        `json:"name"`
	Type         InterfaceType `json:"type"`
	Up           bool          `json:"up"`
	Running      bool          `json:"running"`
	HardwareAddr string        `json:"hardwareAddr"`
	MTU          int           `json:"mtu"`
	Addresses    []string      `json:"addresses"`
}

// LinkStatus contains link layer status information.
type LinkStatus struct {
	Speed      string   `json:"speed"`      // e.g., "1000Mb/s"
	Duplex     string   `json:"duplex"`     // "full" or "half"
	LinkUp     bool     `json:"linkUp"`
	Advertised []string `json:"advertised"` // Advertised link modes
	AutoNeg    bool     `json:"autoNeg"`    // Auto-negotiation enabled
}

// Manager handles network interface operations.
type Manager struct {
	currentInterface string
	interfaces       map[string]*InterfaceInfo
}

// NewManager creates a new network manager.
func NewManager(defaultInterface string) *Manager {
	m := &Manager{
		currentInterface: defaultInterface,
		interfaces:       make(map[string]*InterfaceInfo),
	}
	m.RefreshInterfaces()
	return m
}

// RefreshInterfaces updates the list of available interfaces.
func (m *Manager) RefreshInterfaces() error {
	ifaces, err := net.Interfaces()
	if err != nil {
		return fmt.Errorf("failed to get interfaces: %w", err)
	}

	m.interfaces = make(map[string]*InterfaceInfo)

	for _, iface := range ifaces {
		info := &InterfaceInfo{
			Name:         iface.Name,
			Type:         detectInterfaceType(iface.Name),
			Up:           iface.Flags&net.FlagUp != 0,
			Running:      iface.Flags&net.FlagRunning != 0,
			HardwareAddr: iface.HardwareAddr.String(),
			MTU:          iface.MTU,
			Addresses:    []string{},
		}

		// Get IP addresses
		addrs, err := iface.Addrs()
		if err == nil {
			for _, addr := range addrs {
				info.Addresses = append(info.Addresses, addr.String())
			}
		}

		m.interfaces[iface.Name] = info
	}

	return nil
}

// GetInterfaces returns all available interfaces.
func (m *Manager) GetInterfaces() []*InterfaceInfo {
	result := make([]*InterfaceInfo, 0, len(m.interfaces))
	for _, info := range m.interfaces {
		result = append(result, info)
	}
	return result
}

// GetInterface returns information about a specific interface.
func (m *Manager) GetInterface(name string) (*InterfaceInfo, error) {
	info, ok := m.interfaces[name]
	if !ok {
		return nil, fmt.Errorf("interface %s not found", name)
	}
	return info, nil
}

// GetCurrentInterface returns the currently selected interface.
func (m *Manager) GetCurrentInterface() string {
	return m.currentInterface
}

// SetCurrentInterface sets the active interface.
func (m *Manager) SetCurrentInterface(name string) error {
	if _, ok := m.interfaces[name]; !ok {
		return fmt.Errorf("interface %s not found", name)
	}
	m.currentInterface = name
	return nil
}

// FindFirstAvailable finds the first available interface from a list.
func (m *Manager) FindFirstAvailable(preferred []string) string {
	for _, name := range preferred {
		if info, ok := m.interfaces[name]; ok && info.Up {
			return name
		}
	}

	// Fall back to first non-loopback interface
	for name, info := range m.interfaces {
		if info.Type != InterfaceTypeLoopback && info.Up {
			return name
		}
	}

	return ""
}

// GetLinkStatus returns the link status for an interface.
func (m *Manager) GetLinkStatus(name string) (*LinkStatus, error) {
	info, ok := m.interfaces[name]
	if !ok {
		return nil, fmt.Errorf("interface %s not found", name)
	}

	status := &LinkStatus{
		LinkUp: info.Running,
	}

	// Try to read speed from sysfs (Linux)
	speedPath := filepath.Join("/sys/class/net", name, "speed")
	if data, err := os.ReadFile(speedPath); err == nil {
		speed := strings.TrimSpace(string(data))
		if speed != "" && speed != "-1" {
			status.Speed = speed + "Mb/s"
		}
	}

	// Try to read duplex from sysfs (Linux)
	duplexPath := filepath.Join("/sys/class/net", name, "duplex")
	if data, err := os.ReadFile(duplexPath); err == nil {
		status.Duplex = strings.TrimSpace(string(data))
	}

	return status, nil
}

// detectInterfaceType determines the type of interface from its name.
func detectInterfaceType(name string) InterfaceType {
	// Loopback
	if name == "lo" || name == "lo0" {
		return InterfaceTypeLoopback
	}

	// WiFi interfaces
	wifiPrefixes := []string{"wlan", "wlp", "wifi", "ath", "ra", "wl"}
	for _, prefix := range wifiPrefixes {
		if strings.HasPrefix(name, prefix) {
			return InterfaceTypeWiFi
		}
	}

	// Ethernet interfaces
	ethPrefixes := []string{"eth", "enp", "ens", "eno", "em", "en"}
	for _, prefix := range ethPrefixes {
		if strings.HasPrefix(name, prefix) {
			return InterfaceTypeEthernet
		}
	}

	return InterfaceTypeOther
}

// IsWireless returns true if the interface is a wireless interface.
func (m *Manager) IsWireless(name string) bool {
	info, ok := m.interfaces[name]
	if !ok {
		return false
	}
	return info.Type == InterfaceTypeWiFi
}
