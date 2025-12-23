//go:build darwin

package phy

// getPoEStatus returns PoE status on macOS.
// macOS doesn't expose PoE information through standard interfaces.
func getPoEStatus(_ string) *PoEStatus {
	return &PoEStatus{
		Detected: false,
	}
}

// getSFPInfo returns SFP info on macOS.
// macOS doesn't expose SFP DDM through standard interfaces.
func getSFPInfo(_ string) *SFPInfo {
	return &SFPInfo{
		Present:    false,
		DDMSupport: false,
	}
}
