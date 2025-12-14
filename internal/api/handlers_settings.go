// Package api provides the HTTP/WebSocket server.
package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/krisarmstrong/luminetiq/internal/validation"
)

// ============================================================================
// Settings Handlers (fixes #544 - split from handlers.go)
// ============================================================================

func (s *Server) handleSettings(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.getSettings(w, r)
	case http.MethodPut:
		s.updateSettings(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) getSettings(w http.ResponseWriter, _ *http.Request) {
	// Lock config for read access
	s.config.RLock()
	defer s.config.RUnlock()

	settings := map[string]interface{}{
		"interface": map[string]interface{}{
			"current":   s.config.Interface.Default,
			"available": []string{}, // Will be populated by network module
		},
		"vlan": map[string]interface{}{
			"enabled": s.config.VLAN.Enabled,
			"id":      s.config.VLAN.ID,
		},
		"ip": map[string]interface{}{
			"mode": s.config.IP.Mode,
		},
		"thresholds": map[string]interface{}{
			"dns": map[string]int64{
				"good":    s.config.Thresholds.DNS.Warning.Milliseconds(),
				"warning": s.config.Thresholds.DNS.Critical.Milliseconds(),
			},
			"gateway": map[string]int64{
				"good":    s.config.Thresholds.Ping.Warning.Milliseconds(),
				"warning": s.config.Thresholds.Ping.Critical.Milliseconds(),
			},
			"wifi": map[string]int{
				"good":    s.config.Thresholds.WiFi.Signal.Warning,
				"warning": s.config.Thresholds.WiFi.Signal.Critical,
			},
			"customPing": map[string]int64{
				"good":    s.config.Thresholds.CustomTests.Ping.Warning.Milliseconds(),
				"warning": s.config.Thresholds.CustomTests.Ping.Critical.Milliseconds(),
			},
			"customTcp": map[string]int64{
				"good":    s.config.Thresholds.CustomTests.TCP.Warning.Milliseconds(),
				"warning": s.config.Thresholds.CustomTests.TCP.Critical.Milliseconds(),
			},
			"customHttp": map[string]int64{
				"good":    s.config.Thresholds.CustomTests.HTTP.Warning.Milliseconds(),
				"warning": s.config.Thresholds.CustomTests.HTTP.Critical.Milliseconds(),
			},
			"httpTimings": map[string]map[string]int64{
				"dns": {
					"good":    s.config.Thresholds.CustomTests.HTTPTimings.DNS.Warning.Milliseconds(),
					"warning": s.config.Thresholds.CustomTests.HTTPTimings.DNS.Critical.Milliseconds(),
				},
				"tcp": {
					"good":    s.config.Thresholds.CustomTests.HTTPTimings.TCP.Warning.Milliseconds(),
					"warning": s.config.Thresholds.CustomTests.HTTPTimings.TCP.Critical.Milliseconds(),
				},
				"tls": {
					"good":    s.config.Thresholds.CustomTests.HTTPTimings.TLS.Warning.Milliseconds(),
					"warning": s.config.Thresholds.CustomTests.HTTPTimings.TLS.Critical.Milliseconds(),
				},
				"ttfb": {
					"good":    s.config.Thresholds.CustomTests.HTTPTimings.TTFB.Warning.Milliseconds(),
					"warning": s.config.Thresholds.CustomTests.HTTPTimings.TTFB.Critical.Milliseconds(),
				},
			},
		},
		"tests": map[string]interface{}{
			"runPerformance": s.config.Tests.RunPerformance,
			"runSpeedtest":   s.config.Tests.RunSpeedtest,
			"runIperf":       s.config.Tests.RunIperf,
			"runDiscovery":   s.config.Tests.RunDiscovery,
		},
		"speedtest": map[string]interface{}{
			"serverId":      s.config.Speedtest.ServerID,
			"autoRunOnLink": s.config.Speedtest.AutoRunOnLink,
		},
		"iperf": map[string]interface{}{
			"autoRunOnLink": s.config.Iperf.AutoRunOnLink,
			"server":        s.config.Iperf.Server,
			"port":          s.config.Iperf.Port,
			"protocol":      s.config.Iperf.Protocol,
			"direction":     s.config.Iperf.Direction,
			"duration":      s.config.Iperf.Duration,
			"serverPort":    s.config.Iperf.ServerPort,
			"enableServer":  s.config.Iperf.EnableServer,
		},
		"fabOptions": map[string]interface{}{
			"runLink":             s.config.FABOptions.RunLink,
			"runSwitch":           s.config.FABOptions.RunSwitch,
			"runVLAN":             s.config.FABOptions.RunVLAN,
			"runIPConfig":         s.config.FABOptions.RunIPConfig,
			"runGateway":          s.config.FABOptions.RunGateway,
			"runDNS":              s.config.FABOptions.RunDNS,
			"runHealthChecks":     s.config.FABOptions.RunHealthChecks,
			"runNetworkDiscovery": s.config.FABOptions.RunNetworkDiscovery,
			"runSpeedtest":        s.config.FABOptions.RunSpeedtest,
			"runIperf":            s.config.FABOptions.RunIperf,
			"runPerformance":      s.config.FABOptions.RunPerformance,
			"autoScanOnLink":      s.config.FABOptions.AutoScanOnLink,
		},
		"displayOptions": map[string]interface{}{
			"showPublicIP": s.config.DisplayOptions.ShowPublicIP,
		},
	}

	sendJSONResponse(w, http.StatusOK, settings)
}

func (s *Server) updateSettings(w http.ResponseWriter, r *http.Request) {
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Lock config for write access
	s.config.Lock()
	defer s.config.Unlock()

	// Apply threshold updates
	if thresholds, ok := updates["thresholds"].(map[string]interface{}); ok {
		if dnsThresh, ok := thresholds["dns"].(map[string]interface{}); ok {
			if good, ok := dnsThresh["good"].(float64); ok {
				s.config.Thresholds.DNS.Warning = time.Duration(good) * time.Millisecond
			}
			if warning, ok := dnsThresh["warning"].(float64); ok {
				s.config.Thresholds.DNS.Critical = time.Duration(warning) * time.Millisecond
			}
		}
		if gwThresh, ok := thresholds["gateway"].(map[string]interface{}); ok {
			if good, ok := gwThresh["good"].(float64); ok {
				s.config.Thresholds.Ping.Warning = time.Duration(good) * time.Millisecond
			}
			if warning, ok := gwThresh["warning"].(float64); ok {
				s.config.Thresholds.Ping.Critical = time.Duration(warning) * time.Millisecond
			}
		}
		if wifi, ok := thresholds["wifi"].(map[string]interface{}); ok {
			if good, ok := wifi["good"].(float64); ok {
				s.config.Thresholds.WiFi.Signal.Warning = int(good)
			}
			if warning, ok := wifi["warning"].(float64); ok {
				s.config.Thresholds.WiFi.Signal.Critical = int(warning)
			}
		}
		if customPing, ok := thresholds["customPing"].(map[string]interface{}); ok {
			if good, ok := customPing["good"].(float64); ok {
				s.config.Thresholds.CustomTests.Ping.Warning = time.Duration(good) * time.Millisecond
			}
			if warning, ok := customPing["warning"].(float64); ok {
				s.config.Thresholds.CustomTests.Ping.Critical = time.Duration(warning) * time.Millisecond
			}
		}
		if customTcp, ok := thresholds["customTcp"].(map[string]interface{}); ok {
			if good, ok := customTcp["good"].(float64); ok {
				s.config.Thresholds.CustomTests.TCP.Warning = time.Duration(good) * time.Millisecond
			}
			if warning, ok := customTcp["warning"].(float64); ok {
				s.config.Thresholds.CustomTests.TCP.Critical = time.Duration(warning) * time.Millisecond
			}
		}
		if customHttp, ok := thresholds["customHttp"].(map[string]interface{}); ok {
			if good, ok := customHttp["good"].(float64); ok {
				s.config.Thresholds.CustomTests.HTTP.Warning = time.Duration(good) * time.Millisecond
			}
			if warning, ok := customHttp["warning"].(float64); ok {
				s.config.Thresholds.CustomTests.HTTP.Critical = time.Duration(warning) * time.Millisecond
			}
		}
		if httpTimings, ok := thresholds["httpTimings"].(map[string]interface{}); ok {
			if dnsT, ok := httpTimings["dns"].(map[string]interface{}); ok {
				if good, ok := dnsT["good"].(float64); ok {
					s.config.Thresholds.CustomTests.HTTPTimings.DNS.Warning = time.Duration(good) * time.Millisecond
				}
				if warning, ok := dnsT["warning"].(float64); ok {
					s.config.Thresholds.CustomTests.HTTPTimings.DNS.Critical = time.Duration(warning) * time.Millisecond
				}
			}
			if tcpT, ok := httpTimings["tcp"].(map[string]interface{}); ok {
				if good, ok := tcpT["good"].(float64); ok {
					s.config.Thresholds.CustomTests.HTTPTimings.TCP.Warning = time.Duration(good) * time.Millisecond
				}
				if warning, ok := tcpT["warning"].(float64); ok {
					s.config.Thresholds.CustomTests.HTTPTimings.TCP.Critical = time.Duration(warning) * time.Millisecond
				}
			}
			if tlsT, ok := httpTimings["tls"].(map[string]interface{}); ok {
				if good, ok := tlsT["good"].(float64); ok {
					s.config.Thresholds.CustomTests.HTTPTimings.TLS.Warning = time.Duration(good) * time.Millisecond
				}
				if warning, ok := tlsT["warning"].(float64); ok {
					s.config.Thresholds.CustomTests.HTTPTimings.TLS.Critical = time.Duration(warning) * time.Millisecond
				}
			}
			if ttfb, ok := httpTimings["ttfb"].(map[string]interface{}); ok {
				if good, ok := ttfb["good"].(float64); ok {
					s.config.Thresholds.CustomTests.HTTPTimings.TTFB.Warning = time.Duration(good) * time.Millisecond
				}
				if warning, ok := ttfb["warning"].(float64); ok {
					s.config.Thresholds.CustomTests.HTTPTimings.TTFB.Critical = time.Duration(warning) * time.Millisecond
				}
			}
		}
	}

	// Apply tests updates
	if tests, ok := updates["tests"].(map[string]interface{}); ok {
		if runPerformance, ok := tests["runPerformance"].(bool); ok {
			s.config.Tests.RunPerformance = runPerformance
		}
		if runSpeedtest, ok := tests["runSpeedtest"].(bool); ok {
			s.config.Tests.RunSpeedtest = runSpeedtest
		}
		if runIperf, ok := tests["runIperf"].(bool); ok {
			s.config.Tests.RunIperf = runIperf
		}
		if runDiscovery, ok := tests["runDiscovery"].(bool); ok {
			s.config.Tests.RunDiscovery = runDiscovery
		}
	}

	// Apply speedtest updates
	if speedtest, ok := updates["speedtest"].(map[string]interface{}); ok {
		if serverId, ok := speedtest["serverId"].(string); ok {
			s.config.Speedtest.ServerID = serverId
		}
		if autoRunOnLink, ok := speedtest["autoRunOnLink"].(bool); ok {
			s.config.Speedtest.AutoRunOnLink = autoRunOnLink
		}
	}

	// Apply iperf updates
	if iperf, ok := updates["iperf"].(map[string]interface{}); ok {
		if autoRunOnLink, ok := iperf["autoRunOnLink"].(bool); ok {
			s.config.Iperf.AutoRunOnLink = autoRunOnLink
		}
		if server, ok := iperf["server"].(string); ok {
			s.config.Iperf.Server = server
		}
		if port, ok := iperf["port"].(float64); ok {
			p := int(port)
			if validation.ValidatePort(p) == nil {
				s.config.Iperf.Port = p
			}
		}
		if protocol, ok := iperf["protocol"].(string); ok {
			s.config.Iperf.Protocol = protocol
		}
		if direction, ok := iperf["direction"].(string); ok {
			s.config.Iperf.Direction = direction
		}
		if duration, ok := iperf["duration"].(float64); ok {
			s.config.Iperf.Duration = int(duration)
		}
		if serverPort, ok := iperf["serverPort"].(float64); ok {
			p := int(serverPort)
			if validation.ValidatePort(p) == nil {
				s.config.Iperf.ServerPort = p
			}
		}
		if enableServer, ok := iperf["enableServer"].(bool); ok {
			s.config.Iperf.EnableServer = enableServer
		}
	}

	// Apply fabOptions updates
	if fabOptions, ok := updates["fabOptions"].(map[string]interface{}); ok {
		if runLink, ok := fabOptions["runLink"].(bool); ok {
			s.config.FABOptions.RunLink = runLink
		}
		if runSwitch, ok := fabOptions["runSwitch"].(bool); ok {
			s.config.FABOptions.RunSwitch = runSwitch
		}
		if runVLAN, ok := fabOptions["runVLAN"].(bool); ok {
			s.config.FABOptions.RunVLAN = runVLAN
		}
		if runIPConfig, ok := fabOptions["runIPConfig"].(bool); ok {
			s.config.FABOptions.RunIPConfig = runIPConfig
		}
		if runGateway, ok := fabOptions["runGateway"].(bool); ok {
			s.config.FABOptions.RunGateway = runGateway
		}
		if runDNS, ok := fabOptions["runDNS"].(bool); ok {
			s.config.FABOptions.RunDNS = runDNS
		}
		if runHealthChecks, ok := fabOptions["runHealthChecks"].(bool); ok {
			s.config.FABOptions.RunHealthChecks = runHealthChecks
		}
		if runNetworkDiscovery, ok := fabOptions["runNetworkDiscovery"].(bool); ok {
			s.config.FABOptions.RunNetworkDiscovery = runNetworkDiscovery
		}
		if runSpeedtest, ok := fabOptions["runSpeedtest"].(bool); ok {
			s.config.FABOptions.RunSpeedtest = runSpeedtest
		}
		if runIperf, ok := fabOptions["runIperf"].(bool); ok {
			s.config.FABOptions.RunIperf = runIperf
		}
		if runPerformance, ok := fabOptions["runPerformance"].(bool); ok {
			s.config.FABOptions.RunPerformance = runPerformance
		}
		if autoScanOnLink, ok := fabOptions["autoScanOnLink"].(bool); ok {
			s.config.FABOptions.AutoScanOnLink = autoScanOnLink
		}
	}

	// Apply displayOptions updates
	if displayOptions, ok := updates["displayOptions"].(map[string]interface{}); ok {
		if showPublicIP, ok := displayOptions["showPublicIP"].(bool); ok {
			s.config.DisplayOptions.ShowPublicIP = showPublicIP
		}
	}

	// Save config to file
	if err := s.config.Save(s.configPath); err != nil {
		http.Error(w, "Failed to save config: "+err.Error(), http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, http.StatusOK, map[string]string{"status": "updated"})
}
