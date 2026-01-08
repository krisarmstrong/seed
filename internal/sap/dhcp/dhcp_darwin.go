//go:build darwin

package dhcp

import (
	"bufio"
	"context"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// testDHCPPlatform performs platform-specific DHCP testing on macOS.
// On macOS, we query the current DHCP lease information using ipconfig.
func testDHCPPlatform(ctx context.Context, interfaceName string) *TestResult {
	result := &TestResult{
		Interface: interfaceName,
		TestedAt:  time.Now(),
	}

	// Use ipconfig to get DHCP info
	cmd := exec.CommandContext(ctx, "ipconfig", "getpacket", interfaceName)
	output, err := cmd.Output()
	if err != nil {
		result.Success = false
		result.Error = "failed to get DHCP packet: " + err.Error()
		return result
	}

	// Parse the output
	parseIPConfigOutput(string(output), result)

	if result.OfferedIP != "" {
		result.Success = true
	} else {
		result.Success = false
		result.Error = "no DHCP lease found"
	}

	return result
}

// parseIPConfigOutput parses the output from ipconfig getpacket.
func parseIPConfigOutput(output string, result *TestResult) {
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		parseDHCPLine(line, result)
	}
}

// parseDHCPLine parses a single line from ipconfig getpacket output.
func parseDHCPLine(line string, result *TestResult) {
	switch {
	case strings.HasPrefix(line, "yiaddr = "):
		result.OfferedIP = strings.TrimPrefix(line, "yiaddr = ")
	case strings.HasPrefix(line, "siaddr = "):
		result.ServerIP = strings.TrimPrefix(line, "siaddr = ")
	case strings.HasPrefix(line, "subnet_mask"):
		// Format: subnet_mask (ip): 255.255.255.0
		if idx := strings.LastIndex(line, ": "); idx != -1 {
			result.SubnetMask = strings.TrimSpace(line[idx+2:])
		}
	case strings.HasPrefix(line, "router"):
		// Format: router (ip_mult): {192.168.1.1}
		if idx := strings.Index(line, "{"); idx != -1 {
			end := strings.Index(line, "}")
			if end > idx {
				result.Gateway = strings.TrimSpace(line[idx+1 : end])
			}
		}
	case strings.HasPrefix(line, "domain_name_server"):
		// Format: domain_name_server (ip_mult): {8.8.8.8, 8.8.4.4}
		if idx := strings.Index(line, "{"); idx != -1 {
			end := strings.Index(line, "}")
			if end > idx {
				servers := strings.Split(line[idx+1:end], ",")
				for _, s := range servers {
					s = strings.TrimSpace(s)
					if s != "" {
						result.DNSServers = append(result.DNSServers, s)
					}
				}
			}
		}
	case strings.HasPrefix(line, "domain_name"):
		// Format: domain_name (string): local
		if idx := strings.LastIndex(line, ": "); idx != -1 {
			result.DomainName = strings.TrimSpace(line[idx+2:])
		}
	case strings.HasPrefix(line, "lease_time"):
		// Format: lease_time (uint32): 0x15180
		if idx := strings.LastIndex(line, ": "); idx != -1 {
			valStr := strings.TrimSpace(line[idx+2:])
			if seconds, err := parseLeaseTime(valStr); err == nil {
				result.LeaseTime = time.Duration(seconds) * time.Second
				result.LeaseTimeSec = seconds
			}
		}
	case strings.HasPrefix(line, "server_identifier"):
		// Format: server_identifier (ip): 192.168.1.1
		if idx := strings.LastIndex(line, ": "); idx != -1 {
			if result.ServerIP == "" {
				result.ServerIP = strings.TrimSpace(line[idx+2:])
			}
		}
	}
}

// parseLeaseTime parses a lease time value which may be hex or decimal.
func parseLeaseTime(val string) (int, error) {
	val = strings.TrimSpace(val)
	if strings.HasPrefix(val, "0x") || strings.HasPrefix(val, "0X") {
		// Hex format
		parsed, err := strconv.ParseInt(val[2:], 16, 64)
		return int(parsed), err
	}
	// Decimal format
	parsed, err := strconv.Atoi(val)
	return parsed, err
}

// getCurrentLeasePlatform retrieves the current DHCP lease on macOS.
func getCurrentLeasePlatform(interfaceName string) (*LeaseInfo, error) {
	lease := &LeaseInfo{
		Interface: interfaceName,
	}

	// Get current IP info from interface
	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return nil, &InterfaceError{Message: "interface not found: " + err.Error()}
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return nil, &InterfaceError{Message: "failed to get addresses: " + err.Error()}
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.To4() != nil {
			lease.IPAddress = ipnet.IP.String()
			lease.SubnetMask = net.IP(ipnet.Mask).String()
			break
		}
	}

	// Get DHCP-specific info using ipconfig
	cmd := exec.Command("ipconfig", "getpacket", interfaceName)
	output, err := cmd.Output()
	if err != nil {
		// No DHCP lease, but we might still have static IP info
		if lease.IPAddress != "" {
			return lease, nil
		}
		return nil, &InterfaceError{Message: "no DHCP lease: " + err.Error()}
	}

	// Parse DHCP packet info
	result := &TestResult{}
	parseIPConfigOutput(string(output), result)

	// Copy DHCP info to lease
	if result.ServerIP != "" {
		lease.ServerIP = result.ServerIP
	}
	if result.Gateway != "" {
		lease.Gateway = result.Gateway
	}
	if result.DNSServers != nil {
		lease.DNSServers = result.DNSServers
	}
	if result.DomainName != "" {
		lease.DomainName = result.DomainName
	}
	if result.LeaseTime > 0 {
		lease.LeaseTime = result.LeaseTime
		lease.LeaseTimeSec = int(result.LeaseTime.Seconds())
	}

	return lease, nil
}
