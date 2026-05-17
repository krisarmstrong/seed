// biome-ignore-all lint/complexity/noExcessiveCognitiveComplexity: Complex component
/**
 * SettingsDrawer Component (~1399 lines)
 *
 * Purpose: Master settings configuration panel managing all application settings.
 * Provides tabbed interface to Appearance, Discovery, Network, Performance, Health,
 * DNS, WiFi, SNMP, and Thresholds configuration sections with persistent storage.
 *
 * Key Features:
 * - Multi-section settings: Appearance, Discovery, Network, Performance, Health, DNS, WiFi, SNMP, Thresholds
 * - Persistent storage: Auto-saves settings to server with save status indication
 * - Real-time updates: Changes immediately reflected in UI
 * - Validation: Input validation before saving
 * - Import/Export: Settings backup and restore functionality
 * - Advanced options: Expandable sections for power users
 * - Save status: Shows success/error/unsaved indicators
 * - iperf3 suggestions: Autocomplete for iperf3 server discovery
 * - Keyboard shortcuts: Enter to save, ESC to close
 *
 * Usage:
 * ```typescript
 * <SettingsDrawer isOpen={showSettings} onClose={handleClose} />
 * ```
 *
 * Dependencies: useTheme, useAuth, useSettings context, all settings section components
 * State: All settings state, save status, validation errors, iperf3 suggestions
 */

import type React from 'react';
import { memo, useCallback, useEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useSettings } from '../../contexts/useSettings';
import { useDebouncedAutoSave } from '../../hooks/useDebouncedAutoSave';
import { useSubnetSettings } from '../../hooks/useSubnetSettings';
import { useTheme } from '../../hooks/useTheme';
import { useVulnerabilitySettings } from '../../hooks/useVulnerabilitySettings';
import { LogComponents, logger } from '../../lib/logger';
import { button, cn, icon as iconTokens, layout, radius, spacing } from '../../styles/theme';
import type {
  CableTestSettings as CableTestSettingsType,
  IperfSuggestion,
  IpSettings,
  LinkSettings as LinkSettingsType,
  LogsResponse,
  NetworkDiscoverySettings,
  SettingsThresholds,
  SnmpSettings as SnmpSettingsType,
  TestsSettings,
  WiFiSettings as WiFiSettingsType,
} from '../../types/settings';
import { SettingsDrawerFooter } from './SettingsDrawerFooter';
import { SettingsDrawerNetworkSection } from './SettingsDrawerNetworkSection';
import { AppearanceSettings } from './sections/AppearanceSettings';
import { CableTestSettings } from './sections/CableTestSettings';
import { ConfigBackupsSection } from './sections/ConfigBackupsSection';
import { DiscoverySettings } from './sections/DiscoverySettings';
import { DnsSettings } from './sections/DnsSettings';
import { HealthChecksSettings } from './sections/HealthChecksSettings';
import { LinkSettings } from './sections/LinkSettings';
import { PerformanceSettings } from './sections/PerformanceSettings';
import { ThresholdsSettings } from './sections/ThresholdsSettings';
import { UpdateSettings } from './sections/UpdateSettings';
import { VulnerabilitySettings } from './sections/VulnerabilitySettings';
import { WiFiSettings } from './sections/WiFiSettings';
import {
  INLINE_DEFAULT_CABLE_TEST_SETTINGS,
  INLINE_DEFAULT_LINK_SETTINGS,
  normalizeTestsSettingsForSave,
  withIds,
} from './settingsDrawerNormalizer';

const API_BASE: string = import.meta.env.VITE_API_BASE || '';

interface SettingsDrawerProps {
  isOpen: boolean;
  onClose: () => void;
  version?: string;
  /** Whether currently viewing WiFi mode (shows WiFi settings instead of Link/Cable). */
  isWifi?: boolean;
}

export const SettingsDrawer: React.MemoExoticComponent<
  (props: SettingsDrawerProps) => React.ReactElement | null
> = memo(function settingsDrawer({
  isOpen,
  onClose,
  version = 'dev',
  isWifi = false,
}: SettingsDrawerProps): React.ReactElement | null {
  const { t } = useTranslation('settings');
  const { theme, setTheme, isDark } = useTheme();

  // Get settings from context - single source of truth
  const {
    displayOptions,
    iperfSettings,
    cardSettings,
    status: settingsStatus,
    updateDisplayOptions,
    updateIperfSettings,
    updateCardSettings,
  } = useSettings();

  // Create setter wrappers that use context update methods
  const setDisplayOptions = useCallback(
    (updater: React.SetStateAction<typeof displayOptions>) => {
      const newValue = typeof updater === 'function' ? updater(displayOptions) : updater;
      updateDisplayOptions(newValue);
    },
    [displayOptions, updateDisplayOptions],
  );

  const setIperfSettings = useCallback(
    (updater: React.SetStateAction<typeof iperfSettings>) => {
      const newValue = typeof updater === 'function' ? updater(iperfSettings) : updater;
      updateIperfSettings(newValue);
    },
    [iperfSettings, updateIperfSettings],
  );

  const scrollRef = useRef<HTMLDivElement | null>(null);
  // Track if health check settings were modified - dispatch event on drawer close
  const testsSettingsChangedRef = useRef(false);
  const prevIsOpenRef = useRef(isOpen);

  // Smooth scroll to top when opened
  useEffect(() => {
    if (isOpen && scrollRef.current) {
      scrollRef.current.scrollTop = 0;
    }
  }, [isOpen]);

  // When drawer closes, dispatch healthChecksUpdated if test settings changed
  useEffect(() => {
    if (prevIsOpenRef.current && !isOpen && testsSettingsChangedRef.current) {
      // Drawer just closed and test settings were changed
      window.dispatchEvent(new CustomEvent('healthChecksUpdated'));
      testsSettingsChangedRef.current = false;
    }
    prevIsOpenRef.current = isOpen;
  }, [isOpen]);
  const [thresholds, setThresholds] = useState<SettingsThresholds>({
    dns: { good: 50, warning: 100 },
    gateway: { good: 20, warning: 50 },
    wifi: { good: -50, warning: -70 },
    customPing: { good: 50, warning: 100 },
    customTcp: { good: 100, warning: 500 },
    customHttp: { good: 500, warning: 2000 },
    httpTimings: {
      dns: { good: 100, warning: 500 },
      tcp: { good: 100, warning: 500 },
      tls: { good: 150, warning: 500 },
      ttfb: { good: 500, warning: 2000 },
    },
  });
  const [ipSettings, setIpSettings] = useState<IpSettings>({
    mode: 'dhcp',
    address: '',
    netmask: '24',
    gateway: '',
    dns: [],
  });
  const [testsSettings, setTestsSettings] = useState<TestsSettings>({
    dnsHostname: 'google.com',
    dnsServers: [],
    pingTargets: [],
    tcpPorts: [],
    udpPorts: [],
    httpEndpoints: [],
    runPerformance: true,
    runSpeedtest: true,
    runIperf: true,
    runDiscovery: true,
    speedtest: {
      serverId: '',
      autoRunOnLink: false,
    },
    iperf: {
      autoRunOnLink: false,
    },
  });

  // FAB Options, Display Options, and iperf Settings now come from SettingsContext above

  const [wifiSettings, setWifiSettings] = useState<WiFiSettingsType>({
    interface: '',
    availableWifi: [],
    isWireless: false,
  });
  // Link settings (speed/duplex) - #734
  const [linkSettings, setLinkSettings] = useState<LinkSettingsType>(INLINE_DEFAULT_LINK_SETTINGS);
  // Cable test settings - #734, #740
  const [cableTestSettings, setCableTestSettings] = useState<CableTestSettingsType>(
    INLINE_DEFAULT_CABLE_TEST_SETTINGS,
  );
  const [dnsInput, setDnsInput] = useState('');
  const [iperfSuggestions, setIperfSuggestions] = useState<IperfSuggestion[]>([]);
  const [iperfSuggestionsStatus, setIperfSuggestionsStatus] = useState<
    'idle' | 'loading' | 'error'
  >('idle');
  const [iperfSuggestionsError, setIperfSuggestionsError] = useState<string | null>(null);
  // Network Discovery settings
  const [networkDiscoverySettings, setNetworkDiscoverySettings] =
    useState<NetworkDiscoverySettings>({
      enabled: true,
      arpScanWorkers: 50,
      pingTimeoutMs: 500,
      scanTimeoutMs: 30000,
      autoScan: false,
      scanIntervalMs: 0,
      ipv6Enabled: true,
      options: {
        passiveProtocols: {
          lldp: true,
          cdp: true,
          edp: true,
          ndp: true,
        },
        arpScan: true,
        icmpScan: true,
        portScan: {
          enabled: false,
          preset: 'common',
          tcpPorts: '22,80,443,8080-8100',
          udpPorts: '53,123,161',
          bannerTimeoutMs: 2000,
        },
        tcpProbe: {
          timeoutMs: 2000,
          workers: 20,
        },
        traceroute: false,
        snmpQuery: false,
      },
      timing: {
        probeIntervalMs: 75,
        rescanIntervalMs: 600000,
        workers: 50,
      },
      profiler: {
        enabled: true,
        timeoutMs: 2000,
        maxConcurrent: 5,
        quickPorts: [22, 80, 443, 8080],
      },
      fingerprinting: {
        enabled: false,
        osDetection: false,
        serviceProbes: false,
      },
    });
  // SNMP settings
  const [snmpSettings, setSnmpSettings] = useState<SnmpSettingsType>({
    communities: ['public'],
    v3Credentials: [],
    timeout: 5000,
    retries: 2,
    port: 161,
  });
  // Subnet management lives in its own hook
  const {
    subnets,
    newSubnetCidr,
    setNewSubnetCidr,
    newSubnetName,
    setNewSubnetName,
    subnetError,
    setSubnetError,
    subnetsStatus,
    fetchSubnets,
    addSubnet,
    toggleSubnet,
    deleteSubnet,
  } = useSubnetSettings(isOpen);
  // Log preview (debug)
  const [logPreview, setLogPreview] = useState<string[]>([]);
  const [logLoading, setLogLoading] = useState(false);
  const [logError, setLogError] = useState<string | null>(null);
  // Auto-save status for each section
  type SaveStatus = 'idle' | 'saving' | 'saved' | 'error';
  const [thresholdsStatus, setThresholdsStatus] = useState<SaveStatus>('idle');
  const [testsStatus, setTestsStatus] = useState<SaveStatus>('idle');
  const [wifiStatus, setWifiStatus] = useState<SaveStatus>('idle');
  const [linkStatus, setLinkStatus] = useState<SaveStatus>('idle');
  const [cableTestStatus, setCableTestStatus] = useState<SaveStatus>('idle');
  const [snmpStatus, setSnmpStatus] = useState<SaveStatus>('idle');
  const {
    vulnSettings,
    setVulnSettings,
    vulnStatus,
    vulnInitRef,
    vulnTimerRef,
    fetchVulnSettings,
    saveVulnSettings,
  } = useVulnerabilitySettings();
  // Status for display, iperf comes from context (settingsStatus)
  const displayStatus = settingsStatus.display;
  const iperfStatus = settingsStatus.iperf;

  const [networkDiscoveryStatus, setNetworkDiscoveryStatus] = useState<SaveStatus>('idle');

  // Refs to track initial load (skip auto-save on first load)
  const initialLoadRef = useRef(true);
  const thresholdsInitRef = useRef(true);
  const testsInitRef = useRef(true);
  const wifiInitRef = useRef(true);
  const linkInitRef = useRef(true);
  const cableTestInitRef = useRef(true);
  const networkDiscoveryInitRef = useRef(true);
  const snmpInitRef = useRef(true);

  // Debounce timers (fab, display, iperf now handled by context)
  const thresholdsTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const testsTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const wifiTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const linkTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const cableTestTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const networkDiscoveryTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const snmpTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  // Legacy state (keep for IP settings which still needs manual apply)
  const [savingIp, setSavingIp] = useState(false);
  const [ipMessage, setIpMessage] = useState<string | null>(null);

  // Fetch current thresholds
  const fetchThresholds = useCallback(async () => {
    try {
      const response = await fetch(`${API_BASE}/api/settings`, {
        credentials: 'include',
      });
      if (response.ok) {
        const data = await (response.json() as Promise<Record<string, unknown>>);
        if (data.thresholds) {
          setThresholds((prev) => ({
            ...prev,
            ...data.thresholds,
          }));
        }
      }
    } catch (err) {
      logger.error(LogComponents.CONFIG, 'Failed to fetch thresholds', err);
    }
  }, []);

  // Fetch current IP settings
  const fetchIpSettings = useCallback(async () => {
    try {
      const response = await fetch(`${API_BASE}/api/ipconfig/settings`, {
        credentials: 'include',
      });
      if (response.ok) {
        const data = await (response.json() as Promise<Record<string, unknown>>);
        setIpSettings({
          mode: data.mode || 'dhcp',
          address: data.address || '',
          netmask: data.netmask || '24',
          gateway: data.gateway || '',
          dns: data.dns || [],
        });
        setDnsInput((data.dns || []).join(', '));
      }
    } catch (err) {
      logger.error(LogComponents.CONFIG, 'Failed to fetch IP settings', err);
    }
  }, []);

  // Fetch current tests settings
  const fetchTestsSettings = useCallback(async () => {
    try {
      const response = await fetch(`${API_BASE}/api/health-checks/settings`, {
        credentials: 'include',
      });
      if (response.ok) {
        const data = await (response.json() as Promise<Record<string, unknown>>);
        setTestsSettings({
          dnsHostname: data.dnsHostname || 'google.com',
          dnsServers: withIds(data.dnsServers || []).map((server) => ({
            ...server,
            enabled: server.enabled !== false,
          })),
          pingTargets: withIds(data.pingTargets || []).map((target) => ({
            ...target,
            enabled: target.enabled !== false,
          })),
          tcpPorts: withIds(data.tcpPorts || []).map((port) => ({
            ...port,
            port: port.port || 80,
            enabled: port.enabled !== false,
          })),
          udpPorts: withIds(data.udpPorts || []).map((port) => ({
            ...port,
            port: port.port || 53,
            enabled: port.enabled !== false,
          })),
          httpEndpoints: withIds(data.httpEndpoints || []).map((endpoint) => ({
            ...endpoint,
            expectedStatus: endpoint.expectedStatus || 200,
            enabled: endpoint.enabled !== false,
          })),
          runPerformance: data.runPerformance ?? true,
          runSpeedtest: data.runSpeedtest ?? true,
          runIperf: data.runIperf ?? true,
          runDiscovery: data.runDiscovery ?? true,
          speedtest: {
            serverId: data.speedtest?.serverId || '',
            autoRunOnLink: data.speedtest?.autoRunOnLink ?? true, // Default to true
          },
          iperf: {
            autoRunOnLink: data.iperf?.autoRunOnLink,
          },
        });
      }
    } catch (err) {
      logger.error(LogComponents.CONFIG, 'Failed to fetch tests settings', err);
    }
  }, []);

  const fetchIperfSuggestions = useCallback(async () => {
    setIperfSuggestionsStatus('loading');
    setIperfSuggestionsError(null);
    try {
      const response = await fetch(`${API_BASE}/api/sap/iperf/suggestions`, {
        credentials: 'include',
      });
      if (response.ok) {
        const data = await (response.json() as Promise<Record<string, unknown>>);
        setIperfSuggestions(Array.isArray(data) ? data : []);
        setIperfSuggestionsStatus('idle');
      } else {
        setIperfSuggestionsStatus('error');
        setIperfSuggestionsError('No iperf hosts found');
      }
    } catch (err) {
      setIperfSuggestionsStatus('error');
      setIperfSuggestionsError(err instanceof Error ? err.message : 'Failed to find iperf hosts');
    }
  }, []);

  // Fetch WiFi settings
  const fetchWifiSettings = useCallback(async () => {
    try {
      const response = await fetch(`${API_BASE}/api/canopy/wifi/settings`, {
        credentials: 'include',
      });
      if (response.ok) {
        const data = await (response.json() as Promise<Record<string, unknown>>);
        setWifiSettings({
          interface: data.interface || '',
          availableWifi: data.availableWifi || [],
          isWireless: data.isWireless,
        });
      }
    } catch (err) {
      logger.error(LogComponents.Wifi, 'Failed to fetch WiFi settings', err);
    }
  }, []);

  // FAB options, display options, and iperf settings now come from SettingsContext
  // (loaded automatically by the context provider)

  // Fetch Network Discovery settings from API
  const fetchNetworkDiscoverySettings = useCallback(async () => {
    try {
      const response = await fetch(`${API_BASE}/api/v1/shell/devices/settings`, {
        credentials: 'include',
      });
      if (response.ok) {
        const data = await (response.json() as Promise<Record<string, unknown>>);
        setNetworkDiscoverySettings({
          enabled: data.enabled ?? true,
          arpScanWorkers: data.arpScanWorkers ?? 50,
          pingTimeoutMs: data.pingTimeoutMs ?? 500,
          scanTimeoutMs: data.scanTimeoutMs ?? 30000,
          autoScan: data.autoScan ?? false,
          scanIntervalMs: data.scanIntervalMs ?? 0,
          ipv6Enabled: data.ipv6Enabled ?? true,
          options: data.options ?? {
            passiveProtocols: {
              lldp: true,
              cdp: true,
              edp: true,
              ndp: true,
            },
            arpScan: true,
            icmpScan: true,
            portScan: {
              enabled: false,
              preset: 'common',
              tcpPorts: '22,80,443,8080-8100',
              udpPorts: '53,123,161',
              bannerTimeoutMs: 2000,
            },
            tcpProbe: {
              timeoutMs: 2000,
              workers: 20,
            },
            traceroute: false,
            snmpQuery: false,
          },
          timing: data.timing ?? {
            probeIntervalMs: 75,
            rescanIntervalMs: 600000,
            workers: 50,
          },
          profiler: data.profiler ?? {
            enabled: true,
            timeoutMs: 2000,
            maxConcurrent: 5,
            quickPorts: [22, 80, 443, 8080],
          },
          fingerprinting: data.fingerprinting ?? {
            enabled: false,
            osDetection: false,
            serviceProbes: false,
          },
        });
      }
    } catch (err) {
      logger.error(LogComponents.Discovery, 'Failed to fetch network discovery settings', err);
    }
  }, []);

  // Fetch SNMP settings from API
  const fetchSnmpSettings = useCallback(async () => {
    try {
      const response = await fetch(`${API_BASE}/api/sap/snmp/settings`, {
        credentials: 'include',
      });
      if (response.ok) {
        const data = await (response.json() as Promise<Record<string, unknown>>);
        setSnmpSettings({
          communities: data.communities ?? ['public'],
          v3Credentials: data.v3Credentials ?? [],
          timeout: data.timeout ?? 5000,
          retries: data.retries ?? 2,
          port: data.port ?? 161,
        });
      }
    } catch (err) {
      logger.error(LogComponents.CONFIG, 'Failed to fetch SNMP settings', err);
    }
  }, []);

  // Fetch link settings from API (fixes #734)
  const fetchLinkSettings = useCallback(async () => {
    try {
      const response = await fetch(`${API_BASE}/api/settings/link`, {
        credentials: 'include',
      });
      if (response.ok) {
        const data = await (response.json() as Promise<Record<string, unknown>>);
        // Backend may send separate speed/duplex or combined mode
        const mode = data.mode ?? (data.auto_negotiation ? 'auto' : `${data.speed}/${data.duplex}`);
        setLinkSettings({
          mode: mode,
          availableModes: data.available_modes ?? [],
        });
      }
    } catch (err) {
      logger.error(LogComponents.CONFIG, 'Failed to fetch link settings', err);
    }
  }, []);

  // Fetch cable test settings from API (fixes #740)
  const fetchCableTestSettings = useCallback(async () => {
    try {
      const response = await fetch(`${API_BASE}/api/settings/cable`, {
        credentials: 'include',
      });
      if (response.ok) {
        const data = await (response.json() as Promise<Record<string, unknown>>);
        setCableTestSettings({
          enabled: data.enabled ?? true,
        });
      }
    } catch (err) {
      logger.error(LogComponents.CONFIG, 'Failed to fetch cable test settings', err);
    }
  }, []);

  // Fetch a small tail of the application log (debug)
  // Security fix #301: Removed VITE_LOG_ACCESS_TOKEN - JWT authentication is sufficient
  const fetchLogPreview = useCallback(async () => {
    setLogLoading(true);
    setLogError(null);
    try {
      const response = await fetch(`${API_BASE}/api/harvest/logs?lines=200`, {
        credentials: 'include',
      });
      if (!response.ok) {
        throw new Error('Unable to load logs');
      }
      const data = await (response.json() as Promise<LogsResponse>);
      setLogPreview(data.lines || []);
    } catch (err) {
      setLogPreview([]);
      setLogError(err instanceof Error ? err.message : 'Failed to load log file');
    } finally {
      setLogLoading(false);
    }
  }, []);

  // Save Network Discovery settings to API
  const saveNetworkDiscoverySettings = useCallback(async () => {
    setNetworkDiscoveryStatus('saving');
    try {
      const response = await fetch(`${API_BASE}/api/v1/shell/devices/settings`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify(networkDiscoverySettings),
      });
      if (response.ok) {
        setNetworkDiscoveryStatus('saved');
        setTimeout(() => setNetworkDiscoveryStatus('idle'), 2000);
      } else {
        setNetworkDiscoveryStatus('error');
      }
    } catch {
      setNetworkDiscoveryStatus('error');
    }
  }, [networkDiscoverySettings]);

  const saveSnmpSettings = useCallback(async () => {
    setSnmpStatus('saving');
    try {
      const response = await fetch(`${API_BASE}/api/sap/snmp/settings`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify(snmpSettings),
      });
      if (response.ok) {
        setSnmpStatus('saved');
        setTimeout(() => setSnmpStatus('idle'), 2000);
      } else {
        setSnmpStatus('error');
      }
    } catch {
      setSnmpStatus('error');
    }
  }, [snmpSettings]);
  useEffect(() => {
    if (isOpen) {
      // Reset init refs on open
      initialLoadRef.current = true;
      thresholdsInitRef.current = true;
      testsInitRef.current = true;
      wifiInitRef.current = true;
      linkInitRef.current = true;
      cableTestInitRef.current = true;
      networkDiscoveryInitRef.current = true;
      snmpInitRef.current = true;
      vulnInitRef.current = true;

      fetchThresholds().catch(() => undefined);
      fetchIpSettings().catch(() => undefined);
      fetchTestsSettings().catch(() => undefined);
      fetchWifiSettings().catch(() => undefined);
      // FAB options, display options, and iperf settings come from SettingsContext
      fetchNetworkDiscoverySettings().catch(() => undefined);
      fetchSnmpSettings().catch(() => undefined);
      fetchVulnSettings().catch(() => undefined);
      fetchLinkSettings().catch(() => undefined);
      fetchCableTestSettings().catch(() => undefined);
      fetchSubnets().catch(() => undefined);

      // Mark initial load as done after a short delay
      setTimeout(() => {
        initialLoadRef.current = false;
        thresholdsInitRef.current = false;
        testsInitRef.current = false;
        wifiInitRef.current = false;
        linkInitRef.current = false;
        cableTestInitRef.current = false;
        networkDiscoveryInitRef.current = false;
        snmpInitRef.current = false;
        vulnInitRef.current = false;
      }, 500);
    }
  }, [
    isOpen,
    fetchThresholds,
    fetchIpSettings,
    fetchTestsSettings,
    fetchWifiSettings,
    fetchNetworkDiscoverySettings,
    fetchSnmpSettings,
    fetchVulnSettings,
    fetchLinkSettings,
    fetchCableTestSettings,
    fetchSubnets,
  ]);

  const saveThresholds = useCallback(async () => {
    setThresholdsStatus('saving');
    try {
      const response = await fetch(`${API_BASE}/api/settings`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({ thresholds }),
      });
      if (response.ok) {
        setThresholdsStatus('saved');
        setTimeout(() => setThresholdsStatus('idle'), 2000);
      } else {
        setThresholdsStatus('error');
      }
    } catch {
      setThresholdsStatus('error');
    }
  }, [thresholds]);

  const saveIpSettings = async (): Promise<void> => {
    setSavingIp(true);
    setIpMessage(null);
    try {
      // Parse DNS from input
      const dns = dnsInput
        .split(',')
        .map((s) => s.trim())
        .filter((s) => s.length > 0);

      const response = await fetch(`${API_BASE}/api/ipconfig/settings`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({
          mode: ipSettings.mode,
          address: ipSettings.address,
          netmask: ipSettings.netmask,
          gateway: ipSettings.gateway,
          dns,
        }),
      });
      if (response.ok) {
        setIpMessage('IP settings applied');
        setTimeout(() => setIpMessage(null), 3000);
      } else {
        const error = await (response.text() as Promise<string>);
        setIpMessage(`Failed: ${error}`);
      }
    } catch {
      setIpMessage('Error applying IP settings');
    } finally {
      setSavingIp(false);
    }
  };

  const saveTestsSettings = useCallback(async () => {
    setTestsStatus('saving');
    try {
      const payload = normalizeTestsSettingsForSave(testsSettings);
      const response = await fetch(`${API_BASE}/api/health-checks/settings`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify(payload),
      });
      if (response.ok) {
        setTestsStatus('saved');
        setTimeout(() => setTestsStatus('idle'), 2000);
        // Mark that test settings changed - event dispatched on drawer close
        testsSettingsChangedRef.current = true;
      } else {
        setTestsStatus('error');
      }
    } catch {
      setTestsStatus('error');
    }
  }, [testsSettings]);

  const saveWifiSettings = useCallback(async () => {
    setWifiStatus('saving');
    try {
      const response = await fetch(`${API_BASE}/api/canopy/wifi/settings`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({ interface: wifiSettings.interface }),
      });
      if (response.ok) {
        setWifiStatus('saved');
        setTimeout(() => setWifiStatus('idle'), 2000);
      } else {
        setWifiStatus('error');
      }
    } catch {
      setWifiStatus('error');
    }
  }, [wifiSettings.interface]);

  // Save link settings to backend (fixes #734)
  const saveLinkSettings = useCallback(async () => {
    setLinkStatus('saving');
    try {
      const response = await fetch(`${API_BASE}/api/settings/link`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({
          mode: linkSettings.mode,
          availableModes: linkSettings.availableModes,
        }),
      });
      if (response.ok) {
        setLinkStatus('saved');
        setTimeout(() => setLinkStatus('idle'), 2000);
      } else {
        setLinkStatus('error');
      }
    } catch {
      setLinkStatus('error');
    }
  }, [linkSettings]);

  // Save cable test settings to backend (fixes #740)
  const saveCableTestSettings = useCallback(async () => {
    setCableTestStatus('saving');
    try {
      const response = await fetch(`${API_BASE}/api/settings/cable`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({
          enabled: cableTestSettings.enabled,
        }),
      });
      if (response.ok) {
        setCableTestStatus('saved');
        setTimeout(() => setCableTestStatus('idle'), 2000);
      } else {
        setCableTestStatus('error');
      }
    } catch {
      setCableTestStatus('error');
    }
  }, [cableTestSettings]);

  // Debounced auto-save effects for every settings group
  useDebouncedAutoSave(saveThresholds, thresholdsInitRef, thresholdsTimerRef);
  useDebouncedAutoSave(saveTestsSettings, testsInitRef, testsTimerRef);
  useDebouncedAutoSave(saveWifiSettings, wifiInitRef, wifiTimerRef);
  useDebouncedAutoSave(saveLinkSettings, linkInitRef, linkTimerRef);
  useDebouncedAutoSave(saveCableTestSettings, cableTestInitRef, cableTestTimerRef);
  useDebouncedAutoSave(
    saveNetworkDiscoverySettings,
    networkDiscoveryInitRef,
    networkDiscoveryTimerRef,
  );
  useDebouncedAutoSave(saveSnmpSettings, snmpInitRef, snmpTimerRef);
  useDebouncedAutoSave(saveVulnSettings, vulnInitRef, vulnTimerRef);

  // Fixes #917: Master cleanup effect for all timer refs on unmount
  // Individual useEffects clean up on re-render, but this ensures cleanup on unmount
  useEffect(
    (): (() => void) => (): void => {
      if (thresholdsTimerRef.current) {
        clearTimeout(thresholdsTimerRef.current);
      }
      if (testsTimerRef.current) {
        clearTimeout(testsTimerRef.current);
      }
      if (wifiTimerRef.current) {
        clearTimeout(wifiTimerRef.current);
      }
      if (linkTimerRef.current) {
        clearTimeout(linkTimerRef.current);
      }
      if (cableTestTimerRef.current) {
        clearTimeout(cableTestTimerRef.current);
      }
      if (networkDiscoveryTimerRef.current) {
        clearTimeout(networkDiscoveryTimerRef.current);
      }
      if (snmpTimerRef.current) {
        clearTimeout(snmpTimerRef.current);
      }
      if (vulnTimerRef.current) {
        clearTimeout(vulnTimerRef.current);
      }
    },
    [],
  );

  // Validate IP address format
  const isValidIp = (ip: string): boolean => {
    if (!ip) {
      return true; // Empty is OK for optional fields
    }
    const parts = ip.split('.');
    if (parts.length !== 4) {
      return false;
    }
    return parts.every((p) => {
      const n = Number.parseInt(p, 10);
      return !Number.isNaN(n) && n >= 0 && n <= 255 && p === String(n);
    });
  };

  const drawerRef = useRef<HTMLDivElement>(null);
  const closeButtonRef = useRef<HTMLButtonElement>(null);

  // Handle ESC key to close drawer
  useEffect(() => {
    if (!isOpen) {
      return;
    }

    const handleKeyDown = (e: globalThis.KeyboardEvent): void => {
      if (e.key === 'Escape') {
        onClose();
      }
    };

    document.addEventListener('keydown', handleKeyDown);

    // Focus the close button when drawer opens
    // Fixes #918: Track timeout for cleanup to prevent stale closure
    const focusTimeout = setTimeout(() => closeButtonRef.current?.focus(), 100);

    return (): void => {
      document.removeEventListener('keydown', handleKeyDown);
      clearTimeout(focusTimeout);
    };
  }, [isOpen, onClose]);

  if (!isOpen) {
    return null;
  }

  return (
    <>
      {/* Backdrop */}
      <div
        class="fixed inset-0 bg-black/50 backdrop-blur-sm z-40"
        onClick={onClose}
        aria-hidden="true"
      />

      {/* Drawer - full width on mobile, 384px on larger screens */}
      <div
        ref={drawerRef}
        role="dialog"
        aria-modal="true"
        aria-labelledby="settings-drawer-title"
        onClick={(e: React.MouseEvent): void => e.stopPropagation()}
        onKeyDown={(e: React.KeyboardEvent): void => e.stopPropagation()}
        class="fixed right-0 top-0 h-full w-full sm:w-96 lg:w-lg bg-surface-raised border-l border-surface-border z-50 overflow-y-auto shadow-xl"
      >
        {/* Header */}
        <div
          class={cn(
            layout.flex.between,
            'pad sm:pad-lg border-b border-surface-border sticky top-0 bg-surface-raised z-10',
          )}
        >
          <div class="stack-xs">
            <h2 id="settings-drawer-title" class="heading-3">
              {t('title')}
            </h2>
            <p class="body-small">{t('subtitle')}</p>
          </div>
          <button
            type="button"
            ref={closeButtonRef}
            onClick={onClose}
            class={cn(
              button.size.md,
              radius.md,
              'hover:bg-surface-hover active:bg-surface-hover text-text-muted touch-manipulation focus:outline-none focus:ring-2 focus:ring-brand-primary focus:ring-offset-2 focus:ring-offset-surface-raised',
            )}
            aria-label={t('network.closeSettings')}
          >
            <svg
              class={iconTokens.size.lg}
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              aria-hidden="true"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M6 18L18 6M6 6l12 12"
              />
            </svg>
          </button>
        </div>

        <div
          class={cn(spacing.drawerPad, 'section-gap body-small leading-relaxed')}
          ref={scrollRef}
        >
          {/* Settings sections ordered to match dashboard card order */}
          {/* Link Settings - always visible for ethernet interface config */}
          <LinkSettings
            linkSettings={linkSettings}
            setLinkSettings={setLinkSettings}
            linkStatus={linkStatus}
            cardSettings={cardSettings}
            updateCardSettings={updateCardSettings}
          />

          {/* Cable Test Settings - always visible for cable diagnostics */}
          <CableTestSettings
            cableTestSettings={cableTestSettings}
            setCableTestSettings={setCableTestSettings}
            cableTestStatus={cableTestStatus}
          />

          {/* Network Section - IP/DHCP config (third) */}
          <SettingsDrawerNetworkSection
            ipSettings={ipSettings}
            setIpSettings={setIpSettings}
            dnsInput={dnsInput}
            setDnsInput={setDnsInput}
            saveIpSettings={saveIpSettings}
            savingIp={savingIp}
            ipMessage={ipMessage}
            displayOptions={displayOptions}
            setDisplayOptions={setDisplayOptions}
            displayStatus={displayStatus}
            isValidIp={isValidIp}
          />

          {/* WiFi Settings - only shown in WiFi mode (#754) */}
          {isWifi ? (
            <WiFiSettings
              wifiSettings={wifiSettings}
              setWifiSettings={setWifiSettings}
              wifiStatus={wifiStatus}
            />
          ) : null}

          {/* DNS Settings - matches DnsCard position */}
          <DnsSettings
            testsSettings={testsSettings}
            setTestsSettings={setTestsSettings}
            testsStatus={testsStatus}
            cardSettings={cardSettings}
            updateCardSettings={updateCardSettings}
          />

          <HealthChecksSettings
            testsSettings={testsSettings}
            setTestsSettings={setTestsSettings}
            testsStatus={testsStatus}
            cardSettings={cardSettings}
            updateCardSettings={updateCardSettings}
          />

          <PerformanceSettings
            testsSettings={testsSettings}
            setTestsSettings={setTestsSettings}
            iperfSettings={iperfSettings}
            setIperfSettings={setIperfSettings}
            iperfStatus={iperfStatus}
            iperfSuggestions={iperfSuggestions}
            iperfSuggestionsStatus={iperfSuggestionsStatus}
            iperfSuggestionsError={iperfSuggestionsError}
            fetchIperfSuggestions={fetchIperfSuggestions}
            cardSettings={cardSettings}
            updateCardSettings={updateCardSettings}
          />

          <DiscoverySettings
            networkDiscoverySettings={networkDiscoverySettings}
            setNetworkDiscoverySettings={setNetworkDiscoverySettings}
            networkDiscoveryStatus={networkDiscoveryStatus}
            subnets={subnets}
            subnetsStatus={subnetsStatus}
            newSubnetCidr={newSubnetCidr}
            setNewSubnetCidr={setNewSubnetCidr}
            newSubnetName={newSubnetName}
            setNewSubnetName={setNewSubnetName}
            subnetError={subnetError}
            setSubnetError={setSubnetError}
            addSubnet={addSubnet}
            toggleSubnet={toggleSubnet}
            deleteSubnet={deleteSubnet}
            snmpSettings={snmpSettings}
            setSnmpSettings={setSnmpSettings}
            snmpStatus={snmpStatus}
            cardSettings={cardSettings}
            updateCardSettings={updateCardSettings}
          />

          <VulnerabilitySettings
            settings={vulnSettings}
            setSettings={setVulnSettings}
            status={vulnStatus}
          />

          <ThresholdsSettings
            thresholds={thresholds}
            setThresholds={setThresholds}
            thresholdsStatus={thresholdsStatus}
          />

          {/* Appearance Section */}
          <AppearanceSettings
            theme={theme}
            setTheme={setTheme}
            isDark={isDark}
            unitSystem={displayOptions.unitSystem || 'sae'}
            setUnitSystem={(unit: 'sae' | 'metric'): void =>
              setDisplayOptions((prev) => ({ ...prev, unitSystem: unit }))
            }
          />

          {/* Config Backups Section (implements #494) */}
          <ConfigBackupsSection />

          {/* Updates Section (implements #862) */}
          <UpdateSettings currentVersion={version} />

          <SettingsDrawerFooter
            version={version}
            fetchLogPreview={fetchLogPreview}
            logLoading={logLoading}
            logError={logError}
            logPreview={logPreview}
          />
        </div>
      </div>
    </>
  );
});
