// Package phy provides PHY-level interface detection including PoE and SFP DDM.
package phy

// PoEStatus represents Power over Ethernet detection status.
type PoEStatus struct {
	Detected bool    `json:"detected"`           // Whether PoE is being received
	Standard string  `json:"standard,omitempty"` // 802.3af, 802.3at, 802.3bt
	Class    int     `json:"class,omitempty"`    // PoE class 0-8
	PowerMw  float64 `json:"powerMw,omitempty"`  // Power in milliwatts
	Voltage  float64 `json:"voltage,omitempty"`  // Voltage
}

// SFPInfo represents SFP/SFP+/QSFP module information and DDM (Digital Diagnostics Monitoring).
type SFPInfo struct {
	Present    bool     `json:"present"`              // SFP module detected
	Vendor     string   `json:"vendor,omitempty"`     // Module vendor
	PartNumber string   `json:"partNumber,omitempty"` // Vendor part number
	Serial     string   `json:"serial,omitempty"`     // Serial number
	Type       string   `json:"type,omitempty"`       // SR, LR, ER, etc.
	Wavelength int      `json:"wavelength,omitempty"` // nm (e.g., 850, 1310, 1550)
	Distance   int      `json:"distance,omitempty"`   // Max distance in meters
	Connector  string   `json:"connector,omitempty"`  // LC, SC, etc.
	DDMSupport bool     `json:"ddmSupport"`           // Whether DDM is supported
	DDM        *DDMInfo `json:"ddm,omitempty"`        // DDM readings if supported
}

// DDMInfo contains Digital Diagnostics Monitoring readings from SFP.
type DDMInfo struct {
	Temperature float64  `json:"temperature"` // Celsius
	Voltage     float64  `json:"voltage"`     // Volts
	TxPowerDbm  float64  `json:"txPowerDbm"`  // dBm
	TxPowerMw   float64  `json:"txPowerMw"`   // mW
	RxPowerDbm  float64  `json:"rxPowerDbm"`  // dBm
	RxPowerMw   float64  `json:"rxPowerMw"`   // mW
	LaserBiasMa float64  `json:"laserBiasMa"` // mA
	Alarms      []string `json:"alarms,omitempty"`
	Warnings    []string `json:"warnings,omitempty"`
}

// Detector provides PHY-level detection capabilities.
type Detector struct {
	interfaceName string
}

// NewDetector creates a new PHY detector for the given interface.
func NewDetector(interfaceName string) *Detector {
	return &Detector{
		interfaceName: interfaceName,
	}
}

// SetInterface updates the interface to detect.
func (d *Detector) SetInterface(name string) {
	d.interfaceName = name
}

// GetPoEStatus returns the PoE status for the interface.
func (d *Detector) GetPoEStatus() *PoEStatus {
	return getPoEStatus(d.interfaceName)
}

// GetSFPInfo returns SFP module info and DDM for the interface.
func (d *Detector) GetSFPInfo() *SFPInfo {
	return getSFPInfo(d.interfaceName)
}
