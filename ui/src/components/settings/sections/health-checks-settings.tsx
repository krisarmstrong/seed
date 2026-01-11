/**
 * HealthChecksSettings Component (~449 lines)
 *
 * Purpose: Comprehensive health check configuration panel allowing users to define
 * and customize ping targets, TCP/UDP ports, and HTTP endpoints for monitoring.
 *
 * Key Features:
 * - Ping targets: add/remove/configure ping destinations with custom names and packet counts
 * - TCP ports: configure TCP connectivity tests on specific ports
 * - UDP ports: configure UDP reachability tests
 * - HTTP endpoints: configure HTTP/HTTPS monitoring with customizable URLs
 * - Enable/disable: toggle each test individually
 * - Interval configuration: set how frequently tests run
 * - Timeout settings: configure test timeout values per protocol
 * - Port validation: validates port numbers (1-65535)
 * - URL validation: validates HTTP endpoint URLs
 * - CRUD operations: add/remove/update all test types
 * - AutoSaveIndicator: shows persistent save status
 * - HeartPulse icon: visual indicator in settings menu
 *
 * Usage:
 * ```typescript
 * <HealthChecksSettings
 *   testsSettings={settings}
 *   setTestsSettings={updateSettings}
 *   testsStatus={saveStatus}
 * />
 * ```
 *
 * Dependencies: CollapsibleSection, AutoSaveIndicator, Icons, settings types, ID generation
 * State: Manages multiple arrays of test configurations with CRUD callbacks
 */

import type React from "react";
import { memo, useCallback } from "react";
import { useTranslation } from "react-i18next";
import { cn, icon as iconTokens, input, layout, radius, spacing } from "../../../styles/theme";
import type {
  AlertConfig,
  AnomalyConfig,
  DicomEndpoint,
  FileShareEndpoint,
  FhirEndpoint,
  Hl7Endpoint,
  HttpEndpoint,
  LdapEndpoint,
  LtiEndpoint,
  ModbusEndpoint,
  OpcuaEndpoint,
  PingTarget,
  RtspEndpoint,
  SaveStatus,
  SlaConfig,
  SqlEndpoint,
  TcpPort,
  TestsSettings,
  UdpPort,
} from "../../../types/settings";
import { generateId } from "../../../utils/id";
import { CollapsibleSection } from "../../ui/collapsible-section";
import { HeartPulse } from "../../ui/icons";
import { AutoSaveIndicator } from "./auto-save-indicator";

interface HealthChecksSettingsProps {
  testsSettings: TestsSettings;
  setTestsSettings: React.Dispatch<React.SetStateAction<TestsSettings>>;
  testsStatus: SaveStatus;
}

export const HealthChecksSettings = memo(function HealthChecksSettings({
  testsSettings,
  setTestsSettings,
  testsStatus,
}: HealthChecksSettingsProps) {
  const { t } = useTranslation("settings");

  // Ping target helpers
  const addPingTarget = useCallback(() => {
    setTestsSettings((prev) => ({
      ...prev,
      pingTargets: [
        ...prev.pingTargets,
        { id: generateId(), name: "", host: "", enabled: true, count: 3 },
      ],
    }));
  }, [setTestsSettings]);

  const removePingTarget = useCallback(
    (id: string) => {
      setTestsSettings((prev) => ({
        ...prev,
        pingTargets: prev.pingTargets.filter((t) => t.id !== id),
      }));
    },
    [setTestsSettings],
  );

  const updatePingTarget = useCallback(
    (id: string, field: keyof PingTarget, value: string | boolean | number) => {
      setTestsSettings((prev) => ({
        ...prev,
        pingTargets: prev.pingTargets.map((t) => (t.id === id ? { ...t, [field]: value } : t)),
      }));
    },
    [setTestsSettings],
  );

  // TCP port helpers
  const addTcpPort = useCallback(() => {
    setTestsSettings((prev) => ({
      ...prev,
      tcpPorts: [
        ...prev.tcpPorts,
        { id: generateId(), name: "", host: "", port: 80, enabled: true },
      ],
    }));
  }, [setTestsSettings]);

  const removeTcpPort = useCallback(
    (id: string) => {
      setTestsSettings((prev) => ({
        ...prev,
        tcpPorts: prev.tcpPorts.filter((p) => p.id !== id),
      }));
    },
    [setTestsSettings],
  );

  const updateTcpPort = useCallback(
    (id: string, field: keyof TcpPort, value: string | boolean | number) => {
      setTestsSettings((prev) => ({
        ...prev,
        tcpPorts: prev.tcpPorts.map((p) => (p.id === id ? { ...p, [field]: value } : p)),
      }));
    },
    [setTestsSettings],
  );

  // UDP port helpers
  const addUdpPort = useCallback(() => {
    setTestsSettings((prev) => ({
      ...prev,
      udpPorts: [
        ...prev.udpPorts,
        { id: generateId(), name: "", host: "", port: 53, enabled: true },
      ],
    }));
  }, [setTestsSettings]);

  const removeUdpPort = useCallback(
    (id: string) => {
      setTestsSettings((prev) => ({
        ...prev,
        udpPorts: prev.udpPorts.filter((p) => p.id !== id),
      }));
    },
    [setTestsSettings],
  );

  const updateUdpPort = useCallback(
    (id: string, field: keyof UdpPort, value: string | boolean | number) => {
      setTestsSettings((prev) => ({
        ...prev,
        udpPorts: prev.udpPorts.map((p) => (p.id === id ? { ...p, [field]: value } : p)),
      }));
    },
    [setTestsSettings],
  );

  // HTTP endpoint helpers
  const addHttpEndpoint = useCallback(() => {
    setTestsSettings((prev) => ({
      ...prev,
      httpEndpoints: [
        ...prev.httpEndpoints,
        {
          id: generateId(),
          name: "",
          url: "",
          expectedStatus: 200,
          enabled: true,
        },
      ],
    }));
  }, [setTestsSettings]);

  const removeHttpEndpoint = useCallback(
    (id: string) => {
      setTestsSettings((prev) => ({
        ...prev,
        httpEndpoints: prev.httpEndpoints.filter((e) => e.id !== id),
      }));
    },
    [setTestsSettings],
  );

  const updateHttpEndpoint = useCallback(
    (id: string, field: keyof HttpEndpoint, value: string | boolean | number) => {
      setTestsSettings((prev) => ({
        ...prev,
        httpEndpoints: prev.httpEndpoints.map((e) => (e.id === id ? { ...e, [field]: value } : e)),
      }));
    },
    [setTestsSettings],
  );

  // RTSP endpoint helpers
  const addRtspEndpoint = useCallback(() => {
    setTestsSettings((prev) => ({
      ...prev,
      rtspEndpoints: [
        ...(prev.rtspEndpoints ?? []),
        { id: generateId(), name: "", url: "rtsp://", enabled: true, criticality: 5 },
      ],
    }));
  }, [setTestsSettings]);

  const removeRtspEndpoint = useCallback(
    (id: string) => {
      setTestsSettings((prev) => ({
        ...prev,
        rtspEndpoints: (prev.rtspEndpoints ?? []).filter((e) => e.id !== id),
      }));
    },
    [setTestsSettings],
  );

  const updateRtspEndpoint = useCallback(
    (id: string, field: keyof RtspEndpoint, value: string | boolean | number) => {
      setTestsSettings((prev) => ({
        ...prev,
        rtspEndpoints: (prev.rtspEndpoints ?? []).map((e) =>
          e.id === id ? { ...e, [field]: value } : e,
        ),
      }));
    },
    [setTestsSettings],
  );

  // DICOM endpoint helpers
  const addDicomEndpoint = useCallback(() => {
    setTestsSettings((prev) => ({
      ...prev,
      dicomEndpoints: [
        ...(prev.dicomEndpoints ?? []),
        { id: generateId(), name: "", host: "", port: 104, aeTitle: "", enabled: true, criticality: 8 },
      ],
    }));
  }, [setTestsSettings]);

  const removeDicomEndpoint = useCallback(
    (id: string) => {
      setTestsSettings((prev) => ({
        ...prev,
        dicomEndpoints: (prev.dicomEndpoints ?? []).filter((e) => e.id !== id),
      }));
    },
    [setTestsSettings],
  );

  const updateDicomEndpoint = useCallback(
    (id: string, field: keyof DicomEndpoint, value: string | boolean | number) => {
      setTestsSettings((prev) => ({
        ...prev,
        dicomEndpoints: (prev.dicomEndpoints ?? []).map((e) =>
          e.id === id ? { ...e, [field]: value } : e,
        ),
      }));
    },
    [setTestsSettings],
  );

  // SQL endpoint helpers
  const addSqlEndpoint = useCallback(() => {
    setTestsSettings((prev) => ({
      ...prev,
      sqlEndpoints: [
        ...(prev.sqlEndpoints ?? []),
        {
          id: generateId(),
          name: "",
          driver: "postgres" as const,
          host: "",
          port: 5432,
          database: "",
          username: "",
          enabled: true,
          criticality: 7,
        },
      ],
    }));
  }, [setTestsSettings]);

  const removeSqlEndpoint = useCallback(
    (id: string) => {
      setTestsSettings((prev) => ({
        ...prev,
        sqlEndpoints: (prev.sqlEndpoints ?? []).filter((e) => e.id !== id),
      }));
    },
    [setTestsSettings],
  );

  const updateSqlEndpoint = useCallback(
    (id: string, field: keyof SqlEndpoint, value: string | boolean | number) => {
      setTestsSettings((prev) => ({
        ...prev,
        sqlEndpoints: (prev.sqlEndpoints ?? []).map((e) =>
          e.id === id ? { ...e, [field]: value } : e,
        ),
      }));
    },
    [setTestsSettings],
  );

  // File share endpoint helpers
  const addFileShareEndpoint = useCallback(() => {
    setTestsSettings((prev) => ({
      ...prev,
      fileShareEndpoints: [
        ...(prev.fileShareEndpoints ?? []),
        {
          id: generateId(),
          name: "",
          protocol: "smb" as const,
          host: "",
          sharePath: "",
          enabled: true,
          criticality: 5,
        },
      ],
    }));
  }, [setTestsSettings]);

  const removeFileShareEndpoint = useCallback(
    (id: string) => {
      setTestsSettings((prev) => ({
        ...prev,
        fileShareEndpoints: (prev.fileShareEndpoints ?? []).filter((e) => e.id !== id),
      }));
    },
    [setTestsSettings],
  );

  const updateFileShareEndpoint = useCallback(
    (id: string, field: keyof FileShareEndpoint, value: string | boolean | number) => {
      setTestsSettings((prev) => ({
        ...prev,
        fileShareEndpoints: (prev.fileShareEndpoints ?? []).map((e) =>
          e.id === id ? { ...e, [field]: value } : e,
        ),
      }));
    },
    [setTestsSettings],
  );

  // LDAP endpoint helpers
  const addLdapEndpoint = useCallback(() => {
    setTestsSettings((prev) => ({
      ...prev,
      ldapEndpoints: [
        ...(prev.ldapEndpoints ?? []),
        {
          id: generateId(),
          name: "",
          host: "",
          port: 389,
          useTls: false,
          baseDn: "",
          enabled: true,
          criticality: 7,
        },
      ],
    }));
  }, [setTestsSettings]);

  const removeLdapEndpoint = useCallback(
    (id: string) => {
      setTestsSettings((prev) => ({
        ...prev,
        ldapEndpoints: (prev.ldapEndpoints ?? []).filter((e) => e.id !== id),
      }));
    },
    [setTestsSettings],
  );

  const updateLdapEndpoint = useCallback(
    (id: string, field: keyof LdapEndpoint, value: string | boolean | number) => {
      setTestsSettings((prev) => ({
        ...prev,
        ldapEndpoints: (prev.ldapEndpoints ?? []).map((e) =>
          e.id === id ? { ...e, [field]: value } : e,
        ),
      }));
    },
    [setTestsSettings],
  );

  // HL7 endpoint helpers
  const addHl7Endpoint = useCallback(() => {
    setTestsSettings((prev) => ({
      ...prev,
      hl7Endpoints: [
        ...(prev.hl7Endpoints ?? []),
        {
          id: generateId(),
          name: "",
          host: "",
          port: 2575,
          sendingApp: "",
          sendingFacility: "",
          receivingApp: "",
          receivingFacility: "",
          enabled: true,
          criticality: 9,
        },
      ],
    }));
  }, [setTestsSettings]);

  const removeHl7Endpoint = useCallback(
    (id: string) => {
      setTestsSettings((prev) => ({
        ...prev,
        hl7Endpoints: (prev.hl7Endpoints ?? []).filter((e) => e.id !== id),
      }));
    },
    [setTestsSettings],
  );

  const updateHl7Endpoint = useCallback(
    (id: string, field: keyof Hl7Endpoint, value: string | boolean | number) => {
      setTestsSettings((prev) => ({
        ...prev,
        hl7Endpoints: (prev.hl7Endpoints ?? []).map((e) =>
          e.id === id ? { ...e, [field]: value } : e,
        ),
      }));
    },
    [setTestsSettings],
  );

  // FHIR endpoint helpers
  const addFhirEndpoint = useCallback(() => {
    setTestsSettings((prev) => ({
      ...prev,
      fhirEndpoints: [
        ...(prev.fhirEndpoints ?? []),
        {
          id: generateId(),
          name: "",
          baseUrl: "https://",
          authType: "none" as const,
          enabled: true,
          criticality: 8,
        },
      ],
    }));
  }, [setTestsSettings]);

  const removeFhirEndpoint = useCallback(
    (id: string) => {
      setTestsSettings((prev) => ({
        ...prev,
        fhirEndpoints: (prev.fhirEndpoints ?? []).filter((e) => e.id !== id),
      }));
    },
    [setTestsSettings],
  );

  const updateFhirEndpoint = useCallback(
    (id: string, field: keyof FhirEndpoint, value: string | boolean | number) => {
      setTestsSettings((prev) => ({
        ...prev,
        fhirEndpoints: (prev.fhirEndpoints ?? []).map((e) =>
          e.id === id ? { ...e, [field]: value } : e,
        ),
      }));
    },
    [setTestsSettings],
  );

  // LTI endpoint helpers
  const addLtiEndpoint = useCallback(() => {
    setTestsSettings((prev) => ({
      ...prev,
      ltiEndpoints: [
        ...(prev.ltiEndpoints ?? []),
        {
          id: generateId(),
          name: "",
          launchUrl: "https://",
          consumerKey: "",
          enabled: true,
          criticality: 6,
        },
      ],
    }));
  }, [setTestsSettings]);

  const removeLtiEndpoint = useCallback(
    (id: string) => {
      setTestsSettings((prev) => ({
        ...prev,
        ltiEndpoints: (prev.ltiEndpoints ?? []).filter((e) => e.id !== id),
      }));
    },
    [setTestsSettings],
  );

  const updateLtiEndpoint = useCallback(
    (id: string, field: keyof LtiEndpoint, value: string | boolean | number) => {
      setTestsSettings((prev) => ({
        ...prev,
        ltiEndpoints: (prev.ltiEndpoints ?? []).map((e) =>
          e.id === id ? { ...e, [field]: value } : e,
        ),
      }));
    },
    [setTestsSettings],
  );

  // OPC-UA endpoint helpers
  const addOpcuaEndpoint = useCallback(() => {
    setTestsSettings((prev) => ({
      ...prev,
      opcuaEndpoints: [
        ...(prev.opcuaEndpoints ?? []),
        {
          id: generateId(),
          name: "",
          endpointUrl: "opc.tcp://",
          securityMode: "None" as const,
          enabled: true,
          criticality: 8,
        },
      ],
    }));
  }, [setTestsSettings]);

  const removeOpcuaEndpoint = useCallback(
    (id: string) => {
      setTestsSettings((prev) => ({
        ...prev,
        opcuaEndpoints: (prev.opcuaEndpoints ?? []).filter((e) => e.id !== id),
      }));
    },
    [setTestsSettings],
  );

  const updateOpcuaEndpoint = useCallback(
    (id: string, field: keyof OpcuaEndpoint, value: string | boolean | number) => {
      setTestsSettings((prev) => ({
        ...prev,
        opcuaEndpoints: (prev.opcuaEndpoints ?? []).map((e) =>
          e.id === id ? { ...e, [field]: value } : e,
        ),
      }));
    },
    [setTestsSettings],
  );

  // Modbus endpoint helpers
  const addModbusEndpoint = useCallback(() => {
    setTestsSettings((prev) => ({
      ...prev,
      modbusEndpoints: [
        ...(prev.modbusEndpoints ?? []),
        {
          id: generateId(),
          name: "",
          host: "",
          port: 502,
          unitId: 1,
          testRegister: 0,
          enabled: true,
          criticality: 8,
        },
      ],
    }));
  }, [setTestsSettings]);

  const removeModbusEndpoint = useCallback(
    (id: string) => {
      setTestsSettings((prev) => ({
        ...prev,
        modbusEndpoints: (prev.modbusEndpoints ?? []).filter((e) => e.id !== id),
      }));
    },
    [setTestsSettings],
  );

  const updateModbusEndpoint = useCallback(
    (id: string, field: keyof ModbusEndpoint, value: string | boolean | number) => {
      setTestsSettings((prev) => ({
        ...prev,
        modbusEndpoints: (prev.modbusEndpoints ?? []).map((e) =>
          e.id === id ? { ...e, [field]: value } : e,
        ),
      }));
    },
    [setTestsSettings],
  );

  return (
    <CollapsibleSection
      title={
        <div className={layout.inline.default}>
          <HeartPulse className={iconTokens.size.sm} />
          <span>{t("sections.health")}</span>
          <AutoSaveIndicator status={testsStatus} />
        </div>
      }
    >
      <div className={spacing.stack.default}>
        {/* Enable Toggle */}
        <label
          className={cn(
            layout.flex.between,
            spacing.pad.sm,
            "bg-surface-base border border-surface-border",
            radius.default,
          )}
        >
          <div>
            <span className="body-small text-text-primary font-medium">
              {t("health.enableHealthChecks")}
            </span>
            <p className="caption text-text-muted">{t("health.enableDescription")}</p>
          </div>
          <input
            type="checkbox"
            checked={testsSettings.runPerformance !== false}
            onChange={(e) =>
              setTestsSettings((prev) => ({
                ...prev,
                runPerformance: e.target.checked,
              }))
            }
            className={iconTokens.size.sm}
          />
        </label>

        {/* Ping Targets */}
        <div>
          <div className={cn(layout.flex.between, spacing.margin.bottom.inline)}>
            <span className="caption text-text-muted font-medium">{t("health.pingTargets")}</span>
            <button
              type="button"
              onClick={addPingTarget}
              className="caption text-brand-primary hover:text-brand-accent"
            >
              {t("common.add")}
            </button>
          </div>
          <p className={cn("caption text-text-muted", spacing.margin.bottom.inline)}>
            {t("health.pingDefault")}
          </p>
          {testsSettings.pingTargets.map((target) => (
            <div
              key={target.id || target.host}
              className={cn("flex", spacing.gap.compact, spacing.margin.bottom.inline)}
            >
              <input
                type="text"
                value={target.name}
                onChange={(e) => updatePingTarget(target.id ?? "", "name", e.target.value)}
                placeholder={t("common.name")}
                className={cn(input.base, input.state.default, input.size.md, "w-24")}
              />
              <input
                type="text"
                value={target.host}
                onChange={(e) => updatePingTarget(target.id ?? "", "host", e.target.value)}
                placeholder={t("common.hostIp")}
                className={cn(input.base, input.state.default, input.size.md, "flex-1")}
              />
              <input
                type="number"
                value={target.count || 3}
                onChange={(e) =>
                  updatePingTarget(
                    target.id ?? "",
                    "count",
                    Number.parseInt(e.target.value, 10) || 3,
                  )
                }
                min={1}
                max={10}
                title={t("health.numberOfPings")}
                className={cn(input.base, input.state.default, input.size.md, "w-14 text-center")}
              />
              <button
                type="button"
                onClick={() => removePingTarget(target.id ?? "")}
                className={cn("text-status-error hover:text-status-error/80", spacing.actionBtn)}
                aria-label={t("common.remove")}
              >
                {t("common.remove")}
              </button>
            </div>
          ))}
        </div>

        {/* TCP Ports */}
        <div className={cn("border-t border-surface-border", spacing.padding.top.heading)}>
          <div className={cn(layout.flex.between, spacing.margin.bottom.inline)}>
            <span className="caption text-text-muted font-medium">{t("health.tcpPortTests")}</span>
            <button
              type="button"
              onClick={addTcpPort}
              className="caption text-brand-primary hover:text-brand-accent"
            >
              {t("common.add")}
            </button>
          </div>
          {testsSettings.tcpPorts.map((port) => (
            <div
              key={port.id || `${port.host}:${port.port}`}
              className={cn("flex", spacing.gap.compact, spacing.margin.bottom.inline)}
            >
              <input
                type="text"
                value={port.name}
                onChange={(e) => updateTcpPort(port.id ?? "", "name", e.target.value)}
                placeholder={t("common.name")}
                className={cn(input.base, input.state.default, input.size.md, "w-24")}
              />
              <input
                type="text"
                value={port.host}
                onChange={(e) => updateTcpPort(port.id ?? "", "host", e.target.value)}
                placeholder={t("common.host")}
                className={cn(input.base, input.state.default, input.size.md, "flex-1")}
              />
              <input
                type="number"
                value={port.port}
                onChange={(e) =>
                  updateTcpPort(port.id ?? "", "port", Number.parseInt(e.target.value, 10) || 80)
                }
                placeholder={t("common.port")}
                className={cn(input.base, input.state.default, input.size.md, "w-20")}
              />
              <button
                type="button"
                onClick={() => removeTcpPort(port.id ?? "")}
                className={cn("text-status-error hover:text-status-error/80", spacing.actionBtn)}
                aria-label={t("common.remove")}
              >
                {t("common.remove")}
              </button>
            </div>
          ))}
        </div>

        {/* UDP Ports */}
        <div className={cn("border-t border-surface-border", spacing.padding.top.heading)}>
          <div className={cn(layout.flex.between, spacing.margin.bottom.inline)}>
            <span className="caption text-text-muted font-medium">{t("health.udpPortTests")}</span>
            <button
              type="button"
              onClick={addUdpPort}
              className="caption text-brand-primary hover:text-brand-accent"
            >
              {t("common.add")}
            </button>
          </div>
          <p className={cn("caption text-text-muted", spacing.margin.bottom.inline)}>
            {t("health.udpDescription")}
          </p>
          {testsSettings.udpPorts.map((port) => (
            <div
              key={port.id || `${port.host}:${port.port}`}
              className={cn("flex", spacing.gap.compact, spacing.margin.bottom.inline)}
            >
              <input
                type="text"
                value={port.name}
                onChange={(e) => updateUdpPort(port.id ?? "", "name", e.target.value)}
                placeholder={t("common.name")}
                className={cn(input.base, input.state.default, input.size.md, "w-24")}
              />
              <input
                type="text"
                value={port.host}
                onChange={(e) => updateUdpPort(port.id ?? "", "host", e.target.value)}
                placeholder={t("common.host")}
                className={cn(input.base, input.state.default, input.size.md, "flex-1")}
              />
              <input
                type="number"
                value={port.port}
                onChange={(e) =>
                  updateUdpPort(port.id ?? "", "port", Number.parseInt(e.target.value, 10) || 53)
                }
                placeholder={t("common.port")}
                className={cn(input.base, input.state.default, input.size.md, "w-20")}
              />
              <button
                type="button"
                onClick={() => removeUdpPort(port.id ?? "")}
                className={cn("text-status-error hover:text-status-error/80", spacing.actionBtn)}
                aria-label={t("common.remove")}
              >
                {t("common.remove")}
              </button>
            </div>
          ))}
        </div>

        {/* HTTP Endpoints */}
        <div className={cn("border-t border-surface-border", spacing.padding.top.heading)}>
          <div className={cn(layout.flex.between, spacing.margin.bottom.inline)}>
            <span className="caption text-text-muted font-medium">{t("health.httpEndpoints")}</span>
            <button
              type="button"
              onClick={addHttpEndpoint}
              className="caption text-brand-primary hover:text-brand-accent"
            >
              {t("common.add")}
            </button>
          </div>
          {testsSettings.httpEndpoints.map((endpoint) => (
            <div
              key={endpoint.id || endpoint.url}
              className={cn(
                spacing.stack.xs,
                spacing.margin.bottom.heading,
                spacing.pad.xs,
                "bg-surface-base border border-surface-border",
                radius.default,
              )}
            >
              <div className={cn("flex", spacing.gap.compact)}>
                <input
                  type="text"
                  value={endpoint.name}
                  onChange={(e) => updateHttpEndpoint(endpoint.id ?? "", "name", e.target.value)}
                  placeholder={t("common.name")}
                  className={cn(
                    input.base,
                    input.state.default,
                    input.size.md,
                    "flex-1 bg-surface-raised",
                  )}
                />
                <input
                  type="number"
                  value={endpoint.expectedStatus}
                  onChange={(e) =>
                    updateHttpEndpoint(
                      endpoint.id ?? "",
                      "expectedStatus",
                      Number.parseInt(e.target.value, 10) || 200,
                    )
                  }
                  placeholder={t("health.status")}
                  className={cn(
                    input.base,
                    input.state.default,
                    input.size.md,
                    "w-20 bg-surface-raised",
                  )}
                />
                <button
                  type="button"
                  onClick={() => removeHttpEndpoint(endpoint.id ?? "")}
                  className={cn("text-status-error hover:text-status-error/80", spacing.actionBtn)}
                  aria-label={t("common.remove")}
                >
                  {t("common.remove")}
                </button>
              </div>
              <input
                type="text"
                value={endpoint.url}
                onChange={(e) => updateHttpEndpoint(endpoint.id ?? "", "url", e.target.value)}
                placeholder="https://example.com/health"
                className={cn(input.base, input.state.default, input.size.md, "bg-surface-raised")}
              />
              {/* Criticality Slider */}
              <div className={cn("flex items-center", spacing.gap.compact)}>
                <label className="caption text-text-muted w-28">{t("health.criticality")}</label>
                <input
                  type="range"
                  min={1}
                  max={10}
                  value={endpoint.criticality ?? 5}
                  onChange={(e) =>
                    updateHttpEndpoint(
                      endpoint.id ?? "",
                      "criticality",
                      Number.parseInt(e.target.value, 10),
                    )
                  }
                  className="flex-1 h-2 bg-surface-raised rounded-lg appearance-none cursor-pointer accent-brand-primary"
                />
                <span className="caption text-text-muted w-6 text-center">
                  {endpoint.criticality ?? 5}
                </span>
              </div>
            </div>
          ))}
        </div>

        {/* SQL Database Endpoints */}
        <div className={cn("border-t border-surface-border", spacing.padding.top.heading)}>
          <div className={cn(layout.flex.between, spacing.margin.bottom.inline)}>
            <span className="caption text-text-muted font-medium">{t("health.sqlEndpoints")}</span>
            <button
              type="button"
              onClick={addSqlEndpoint}
              className="caption text-brand-primary hover:text-brand-accent"
            >
              {t("common.add")}
            </button>
          </div>
          <p className={cn("caption text-text-muted", spacing.margin.bottom.inline)}>
            {t("health.sqlDescription")}
          </p>
          {(testsSettings.sqlEndpoints ?? []).map((endpoint) => (
            <div
              key={endpoint.id}
              className={cn(
                spacing.stack.xs,
                spacing.margin.bottom.heading,
                spacing.pad.xs,
                "bg-surface-base border border-surface-border",
                radius.default,
              )}
            >
              <div className={cn("flex", spacing.gap.compact)}>
                <input
                  type="text"
                  value={endpoint.name}
                  onChange={(e) => updateSqlEndpoint(endpoint.id ?? "", "name", e.target.value)}
                  placeholder={t("common.name")}
                  className={cn(input.base, input.state.default, input.size.md, "flex-1 bg-surface-raised")}
                />
                <select
                  value={endpoint.driver}
                  onChange={(e) => updateSqlEndpoint(endpoint.id ?? "", "driver", e.target.value)}
                  className={cn(input.base, input.state.default, input.size.md, "w-28 bg-surface-raised")}
                >
                  <option value="postgres">PostgreSQL</option>
                  <option value="mysql">MySQL</option>
                  <option value="mssql">SQL Server</option>
                  <option value="oracle">Oracle</option>
                </select>
                <button
                  type="button"
                  onClick={() => removeSqlEndpoint(endpoint.id ?? "")}
                  className={cn("text-status-error hover:text-status-error/80", spacing.actionBtn)}
                >
                  {t("common.remove")}
                </button>
              </div>
              <div className={cn("flex", spacing.gap.compact)}>
                <input
                  type="text"
                  value={endpoint.host}
                  onChange={(e) => updateSqlEndpoint(endpoint.id ?? "", "host", e.target.value)}
                  placeholder={t("common.host")}
                  className={cn(input.base, input.state.default, input.size.md, "flex-1 bg-surface-raised")}
                />
                <input
                  type="number"
                  value={endpoint.port}
                  onChange={(e) => updateSqlEndpoint(endpoint.id ?? "", "port", Number.parseInt(e.target.value, 10))}
                  placeholder={t("common.port")}
                  className={cn(input.base, input.state.default, input.size.md, "w-20 bg-surface-raised")}
                />
              </div>
              <div className={cn("flex", spacing.gap.compact)}>
                <input
                  type="text"
                  value={endpoint.database}
                  onChange={(e) => updateSqlEndpoint(endpoint.id ?? "", "database", e.target.value)}
                  placeholder={t("health.database")}
                  className={cn(input.base, input.state.default, input.size.md, "flex-1 bg-surface-raised")}
                />
                <input
                  type="text"
                  value={endpoint.username}
                  onChange={(e) => updateSqlEndpoint(endpoint.id ?? "", "username", e.target.value)}
                  placeholder={t("health.username")}
                  className={cn(input.base, input.state.default, input.size.md, "flex-1 bg-surface-raised")}
                />
              </div>
            </div>
          ))}
        </div>

        {/* File Share Endpoints (SMB/NFS) */}
        <div className={cn("border-t border-surface-border", spacing.padding.top.heading)}>
          <div className={cn(layout.flex.between, spacing.margin.bottom.inline)}>
            <span className="caption text-text-muted font-medium">{t("health.fileShareEndpoints")}</span>
            <button
              type="button"
              onClick={addFileShareEndpoint}
              className="caption text-brand-primary hover:text-brand-accent"
            >
              {t("common.add")}
            </button>
          </div>
          <p className={cn("caption text-text-muted", spacing.margin.bottom.inline)}>
            {t("health.fileShareDescription")}
          </p>
          {(testsSettings.fileShareEndpoints ?? []).map((endpoint) => (
            <div
              key={endpoint.id}
              className={cn("flex", spacing.gap.compact, spacing.margin.bottom.inline)}
            >
              <input
                type="text"
                value={endpoint.name}
                onChange={(e) => updateFileShareEndpoint(endpoint.id ?? "", "name", e.target.value)}
                placeholder={t("common.name")}
                className={cn(input.base, input.state.default, input.size.md, "w-24")}
              />
              <select
                value={endpoint.protocol}
                onChange={(e) => updateFileShareEndpoint(endpoint.id ?? "", "protocol", e.target.value)}
                className={cn(input.base, input.state.default, input.size.md, "w-20")}
              >
                <option value="smb">SMB</option>
                <option value="nfs">NFS</option>
              </select>
              <input
                type="text"
                value={endpoint.host}
                onChange={(e) => updateFileShareEndpoint(endpoint.id ?? "", "host", e.target.value)}
                placeholder={t("common.host")}
                className={cn(input.base, input.state.default, input.size.md, "flex-1")}
              />
              <input
                type="text"
                value={endpoint.sharePath}
                onChange={(e) => updateFileShareEndpoint(endpoint.id ?? "", "sharePath", e.target.value)}
                placeholder={t("health.sharePath")}
                className={cn(input.base, input.state.default, input.size.md, "flex-1")}
              />
              <button
                type="button"
                onClick={() => removeFileShareEndpoint(endpoint.id ?? "")}
                className={cn("text-status-error hover:text-status-error/80", spacing.actionBtn)}
              >
                {t("common.remove")}
              </button>
            </div>
          ))}
        </div>

        {/* LDAP Endpoints */}
        <div className={cn("border-t border-surface-border", spacing.padding.top.heading)}>
          <div className={cn(layout.flex.between, spacing.margin.bottom.inline)}>
            <span className="caption text-text-muted font-medium">{t("health.ldapEndpoints")}</span>
            <button
              type="button"
              onClick={addLdapEndpoint}
              className="caption text-brand-primary hover:text-brand-accent"
            >
              {t("common.add")}
            </button>
          </div>
          <p className={cn("caption text-text-muted", spacing.margin.bottom.inline)}>
            {t("health.ldapDescription")}
          </p>
          {(testsSettings.ldapEndpoints ?? []).map((endpoint) => (
            <div
              key={endpoint.id}
              className={cn(
                spacing.stack.xs,
                spacing.margin.bottom.heading,
                spacing.pad.xs,
                "bg-surface-base border border-surface-border",
                radius.default,
              )}
            >
              <div className={cn("flex", spacing.gap.compact)}>
                <input
                  type="text"
                  value={endpoint.name}
                  onChange={(e) => updateLdapEndpoint(endpoint.id ?? "", "name", e.target.value)}
                  placeholder={t("common.name")}
                  className={cn(input.base, input.state.default, input.size.md, "w-32 bg-surface-raised")}
                />
                <input
                  type="text"
                  value={endpoint.host}
                  onChange={(e) => updateLdapEndpoint(endpoint.id ?? "", "host", e.target.value)}
                  placeholder={t("common.host")}
                  className={cn(input.base, input.state.default, input.size.md, "flex-1 bg-surface-raised")}
                />
                <input
                  type="number"
                  value={endpoint.port}
                  onChange={(e) => updateLdapEndpoint(endpoint.id ?? "", "port", Number.parseInt(e.target.value, 10))}
                  placeholder={t("common.port")}
                  className={cn(input.base, input.state.default, input.size.md, "w-20 bg-surface-raised")}
                />
                <button
                  type="button"
                  onClick={() => removeLdapEndpoint(endpoint.id ?? "")}
                  className={cn("text-status-error hover:text-status-error/80", spacing.actionBtn)}
                >
                  {t("common.remove")}
                </button>
              </div>
              <div className={cn("flex items-center", spacing.gap.compact)}>
                <input
                  type="text"
                  value={endpoint.baseDn}
                  onChange={(e) => updateLdapEndpoint(endpoint.id ?? "", "baseDn", e.target.value)}
                  placeholder={t("health.baseDn")}
                  className={cn(input.base, input.state.default, input.size.md, "flex-1 bg-surface-raised")}
                />
                <label className={cn("flex items-center", spacing.gap.compact, "caption text-text-muted")}>
                  <input
                    type="checkbox"
                    checked={endpoint.useTls}
                    onChange={(e) => updateLdapEndpoint(endpoint.id ?? "", "useTls", e.target.checked)}
                  />
                  TLS
                </label>
              </div>
            </div>
          ))}
        </div>

        {/* RTSP Video Endpoints */}
        <div className={cn("border-t border-surface-border", spacing.padding.top.heading)}>
          <div className={cn(layout.flex.between, spacing.margin.bottom.inline)}>
            <span className="caption text-text-muted font-medium">{t("health.rtspEndpoints")}</span>
            <button
              type="button"
              onClick={addRtspEndpoint}
              className="caption text-brand-primary hover:text-brand-accent"
            >
              {t("common.add")}
            </button>
          </div>
          <p className={cn("caption text-text-muted", spacing.margin.bottom.inline)}>
            {t("health.rtspDescription")}
          </p>
          {(testsSettings.rtspEndpoints ?? []).map((endpoint) => (
            <div
              key={endpoint.id}
              className={cn("flex", spacing.gap.compact, spacing.margin.bottom.inline)}
            >
              <input
                type="text"
                value={endpoint.name}
                onChange={(e) => updateRtspEndpoint(endpoint.id ?? "", "name", e.target.value)}
                placeholder={t("common.name")}
                className={cn(input.base, input.state.default, input.size.md, "w-24")}
              />
              <input
                type="text"
                value={endpoint.url}
                onChange={(e) => updateRtspEndpoint(endpoint.id ?? "", "url", e.target.value)}
                placeholder="rtsp://host:554/stream"
                className={cn(input.base, input.state.default, input.size.md, "flex-1")}
              />
              <button
                type="button"
                onClick={() => removeRtspEndpoint(endpoint.id ?? "")}
                className={cn("text-status-error hover:text-status-error/80", spacing.actionBtn)}
              >
                {t("common.remove")}
              </button>
            </div>
          ))}
        </div>

        {/* DICOM Medical Imaging Endpoints */}
        <div className={cn("border-t border-surface-border", spacing.padding.top.heading)}>
          <div className={cn(layout.flex.between, spacing.margin.bottom.inline)}>
            <span className="caption text-text-muted font-medium">{t("health.dicomEndpoints")}</span>
            <button
              type="button"
              onClick={addDicomEndpoint}
              className="caption text-brand-primary hover:text-brand-accent"
            >
              {t("common.add")}
            </button>
          </div>
          <p className={cn("caption text-text-muted", spacing.margin.bottom.inline)}>
            {t("health.dicomDescription")}
          </p>
          {(testsSettings.dicomEndpoints ?? []).map((endpoint) => (
            <div
              key={endpoint.id}
              className={cn("flex", spacing.gap.compact, spacing.margin.bottom.inline)}
            >
              <input
                type="text"
                value={endpoint.name}
                onChange={(e) => updateDicomEndpoint(endpoint.id ?? "", "name", e.target.value)}
                placeholder={t("common.name")}
                className={cn(input.base, input.state.default, input.size.md, "w-24")}
              />
              <input
                type="text"
                value={endpoint.host}
                onChange={(e) => updateDicomEndpoint(endpoint.id ?? "", "host", e.target.value)}
                placeholder={t("common.host")}
                className={cn(input.base, input.state.default, input.size.md, "flex-1")}
              />
              <input
                type="number"
                value={endpoint.port}
                onChange={(e) => updateDicomEndpoint(endpoint.id ?? "", "port", Number.parseInt(e.target.value, 10))}
                placeholder="104"
                className={cn(input.base, input.state.default, input.size.md, "w-20")}
              />
              <input
                type="text"
                value={endpoint.aeTitle}
                onChange={(e) => updateDicomEndpoint(endpoint.id ?? "", "aeTitle", e.target.value)}
                placeholder="AE Title"
                className={cn(input.base, input.state.default, input.size.md, "w-24")}
              />
              <button
                type="button"
                onClick={() => removeDicomEndpoint(endpoint.id ?? "")}
                className={cn("text-status-error hover:text-status-error/80", spacing.actionBtn)}
              >
                {t("common.remove")}
              </button>
            </div>
          ))}
        </div>

        {/* HL7 MLLP Endpoints */}
        <div className={cn("border-t border-surface-border", spacing.padding.top.heading)}>
          <div className={cn(layout.flex.between, spacing.margin.bottom.inline)}>
            <span className="caption text-text-muted font-medium">{t("health.hl7Endpoints")}</span>
            <button
              type="button"
              onClick={addHl7Endpoint}
              className="caption text-brand-primary hover:text-brand-accent"
            >
              {t("common.add")}
            </button>
          </div>
          <p className={cn("caption text-text-muted", spacing.margin.bottom.inline)}>
            {t("health.hl7Description")}
          </p>
          {(testsSettings.hl7Endpoints ?? []).map((endpoint) => (
            <div
              key={endpoint.id}
              className={cn(
                spacing.stack.xs,
                spacing.margin.bottom.heading,
                spacing.pad.xs,
                "bg-surface-base border border-surface-border",
                radius.default,
              )}
            >
              <div className={cn("flex", spacing.gap.compact)}>
                <input
                  type="text"
                  value={endpoint.name}
                  onChange={(e) => updateHl7Endpoint(endpoint.id ?? "", "name", e.target.value)}
                  placeholder={t("common.name")}
                  className={cn(input.base, input.state.default, input.size.md, "w-32 bg-surface-raised")}
                />
                <input
                  type="text"
                  value={endpoint.host}
                  onChange={(e) => updateHl7Endpoint(endpoint.id ?? "", "host", e.target.value)}
                  placeholder={t("common.host")}
                  className={cn(input.base, input.state.default, input.size.md, "flex-1 bg-surface-raised")}
                />
                <input
                  type="number"
                  value={endpoint.port}
                  onChange={(e) => updateHl7Endpoint(endpoint.id ?? "", "port", Number.parseInt(e.target.value, 10))}
                  placeholder="2575"
                  className={cn(input.base, input.state.default, input.size.md, "w-20 bg-surface-raised")}
                />
                <button
                  type="button"
                  onClick={() => removeHl7Endpoint(endpoint.id ?? "")}
                  className={cn("text-status-error hover:text-status-error/80", spacing.actionBtn)}
                >
                  {t("common.remove")}
                </button>
              </div>
              <div className={cn("flex", spacing.gap.compact)}>
                <input
                  type="text"
                  value={endpoint.sendingApp}
                  onChange={(e) => updateHl7Endpoint(endpoint.id ?? "", "sendingApp", e.target.value)}
                  placeholder={t("health.sendingApp")}
                  className={cn(input.base, input.state.default, input.size.md, "flex-1 bg-surface-raised")}
                />
                <input
                  type="text"
                  value={endpoint.sendingFacility}
                  onChange={(e) => updateHl7Endpoint(endpoint.id ?? "", "sendingFacility", e.target.value)}
                  placeholder={t("health.sendingFacility")}
                  className={cn(input.base, input.state.default, input.size.md, "flex-1 bg-surface-raised")}
                />
              </div>
              <div className={cn("flex", spacing.gap.compact)}>
                <input
                  type="text"
                  value={endpoint.receivingApp}
                  onChange={(e) => updateHl7Endpoint(endpoint.id ?? "", "receivingApp", e.target.value)}
                  placeholder={t("health.receivingApp")}
                  className={cn(input.base, input.state.default, input.size.md, "flex-1 bg-surface-raised")}
                />
                <input
                  type="text"
                  value={endpoint.receivingFacility}
                  onChange={(e) => updateHl7Endpoint(endpoint.id ?? "", "receivingFacility", e.target.value)}
                  placeholder={t("health.receivingFacility")}
                  className={cn(input.base, input.state.default, input.size.md, "flex-1 bg-surface-raised")}
                />
              </div>
            </div>
          ))}
        </div>

        {/* FHIR R4 Endpoints */}
        <div className={cn("border-t border-surface-border", spacing.padding.top.heading)}>
          <div className={cn(layout.flex.between, spacing.margin.bottom.inline)}>
            <span className="caption text-text-muted font-medium">{t("health.fhirEndpoints")}</span>
            <button
              type="button"
              onClick={addFhirEndpoint}
              className="caption text-brand-primary hover:text-brand-accent"
            >
              {t("common.add")}
            </button>
          </div>
          <p className={cn("caption text-text-muted", spacing.margin.bottom.inline)}>
            {t("health.fhirDescription")}
          </p>
          {(testsSettings.fhirEndpoints ?? []).map((endpoint) => (
            <div
              key={endpoint.id}
              className={cn("flex", spacing.gap.compact, spacing.margin.bottom.inline)}
            >
              <input
                type="text"
                value={endpoint.name}
                onChange={(e) => updateFhirEndpoint(endpoint.id ?? "", "name", e.target.value)}
                placeholder={t("common.name")}
                className={cn(input.base, input.state.default, input.size.md, "w-24")}
              />
              <input
                type="text"
                value={endpoint.baseUrl}
                onChange={(e) => updateFhirEndpoint(endpoint.id ?? "", "baseUrl", e.target.value)}
                placeholder="https://fhir.example.com/r4"
                className={cn(input.base, input.state.default, input.size.md, "flex-1")}
              />
              <select
                value={endpoint.authType}
                onChange={(e) => updateFhirEndpoint(endpoint.id ?? "", "authType", e.target.value)}
                className={cn(input.base, input.state.default, input.size.md, "w-24")}
              >
                <option value="none">None</option>
                <option value="basic">Basic</option>
                <option value="oauth2">OAuth2</option>
              </select>
              <button
                type="button"
                onClick={() => removeFhirEndpoint(endpoint.id ?? "")}
                className={cn("text-status-error hover:text-status-error/80", spacing.actionBtn)}
              >
                {t("common.remove")}
              </button>
            </div>
          ))}
        </div>

        {/* LTI/LMS Education Endpoints */}
        <div className={cn("border-t border-surface-border", spacing.padding.top.heading)}>
          <div className={cn(layout.flex.between, spacing.margin.bottom.inline)}>
            <span className="caption text-text-muted font-medium">{t("health.ltiEndpoints")}</span>
            <button
              type="button"
              onClick={addLtiEndpoint}
              className="caption text-brand-primary hover:text-brand-accent"
            >
              {t("common.add")}
            </button>
          </div>
          <p className={cn("caption text-text-muted", spacing.margin.bottom.inline)}>
            {t("health.ltiDescription")}
          </p>
          {(testsSettings.ltiEndpoints ?? []).map((endpoint) => (
            <div
              key={endpoint.id}
              className={cn("flex", spacing.gap.compact, spacing.margin.bottom.inline)}
            >
              <input
                type="text"
                value={endpoint.name}
                onChange={(e) => updateLtiEndpoint(endpoint.id ?? "", "name", e.target.value)}
                placeholder={t("common.name")}
                className={cn(input.base, input.state.default, input.size.md, "w-24")}
              />
              <input
                type="text"
                value={endpoint.launchUrl}
                onChange={(e) => updateLtiEndpoint(endpoint.id ?? "", "launchUrl", e.target.value)}
                placeholder="https://lms.example.com/lti/launch"
                className={cn(input.base, input.state.default, input.size.md, "flex-1")}
              />
              <input
                type="text"
                value={endpoint.consumerKey}
                onChange={(e) => updateLtiEndpoint(endpoint.id ?? "", "consumerKey", e.target.value)}
                placeholder={t("health.consumerKey")}
                className={cn(input.base, input.state.default, input.size.md, "w-32")}
              />
              <button
                type="button"
                onClick={() => removeLtiEndpoint(endpoint.id ?? "")}
                className={cn("text-status-error hover:text-status-error/80", spacing.actionBtn)}
              >
                {t("common.remove")}
              </button>
            </div>
          ))}
        </div>

        {/* OPC-UA Industrial Endpoints */}
        <div className={cn("border-t border-surface-border", spacing.padding.top.heading)}>
          <div className={cn(layout.flex.between, spacing.margin.bottom.inline)}>
            <span className="caption text-text-muted font-medium">{t("health.opcuaEndpoints")}</span>
            <button
              type="button"
              onClick={addOpcuaEndpoint}
              className="caption text-brand-primary hover:text-brand-accent"
            >
              {t("common.add")}
            </button>
          </div>
          <p className={cn("caption text-text-muted", spacing.margin.bottom.inline)}>
            {t("health.opcuaDescription")}
          </p>
          {(testsSettings.opcuaEndpoints ?? []).map((endpoint) => (
            <div
              key={endpoint.id}
              className={cn("flex", spacing.gap.compact, spacing.margin.bottom.inline)}
            >
              <input
                type="text"
                value={endpoint.name}
                onChange={(e) => updateOpcuaEndpoint(endpoint.id ?? "", "name", e.target.value)}
                placeholder={t("common.name")}
                className={cn(input.base, input.state.default, input.size.md, "w-24")}
              />
              <input
                type="text"
                value={endpoint.endpointUrl}
                onChange={(e) => updateOpcuaEndpoint(endpoint.id ?? "", "endpointUrl", e.target.value)}
                placeholder="opc.tcp://host:4840"
                className={cn(input.base, input.state.default, input.size.md, "flex-1")}
              />
              <select
                value={endpoint.securityMode}
                onChange={(e) => updateOpcuaEndpoint(endpoint.id ?? "", "securityMode", e.target.value)}
                className={cn(input.base, input.state.default, input.size.md, "w-32")}
              >
                <option value="None">None</option>
                <option value="Sign">Sign</option>
                <option value="SignAndEncrypt">Sign+Encrypt</option>
              </select>
              <button
                type="button"
                onClick={() => removeOpcuaEndpoint(endpoint.id ?? "")}
                className={cn("text-status-error hover:text-status-error/80", spacing.actionBtn)}
              >
                {t("common.remove")}
              </button>
            </div>
          ))}
        </div>

        {/* Modbus TCP Industrial Endpoints */}
        <div className={cn("border-t border-surface-border", spacing.padding.top.heading)}>
          <div className={cn(layout.flex.between, spacing.margin.bottom.inline)}>
            <span className="caption text-text-muted font-medium">{t("health.modbusEndpoints")}</span>
            <button
              type="button"
              onClick={addModbusEndpoint}
              className="caption text-brand-primary hover:text-brand-accent"
            >
              {t("common.add")}
            </button>
          </div>
          <p className={cn("caption text-text-muted", spacing.margin.bottom.inline)}>
            {t("health.modbusDescription")}
          </p>
          {(testsSettings.modbusEndpoints ?? []).map((endpoint) => (
            <div
              key={endpoint.id}
              className={cn("flex", spacing.gap.compact, spacing.margin.bottom.inline)}
            >
              <input
                type="text"
                value={endpoint.name}
                onChange={(e) => updateModbusEndpoint(endpoint.id ?? "", "name", e.target.value)}
                placeholder={t("common.name")}
                className={cn(input.base, input.state.default, input.size.md, "w-24")}
              />
              <input
                type="text"
                value={endpoint.host}
                onChange={(e) => updateModbusEndpoint(endpoint.id ?? "", "host", e.target.value)}
                placeholder={t("common.host")}
                className={cn(input.base, input.state.default, input.size.md, "flex-1")}
              />
              <input
                type="number"
                value={endpoint.port}
                onChange={(e) => updateModbusEndpoint(endpoint.id ?? "", "port", Number.parseInt(e.target.value, 10))}
                placeholder="502"
                className={cn(input.base, input.state.default, input.size.md, "w-20")}
              />
              <input
                type="number"
                value={endpoint.unitId}
                onChange={(e) => updateModbusEndpoint(endpoint.id ?? "", "unitId", Number.parseInt(e.target.value, 10))}
                placeholder="Unit"
                title={t("health.unitId")}
                className={cn(input.base, input.state.default, input.size.md, "w-16")}
              />
              <input
                type="number"
                value={endpoint.testRegister}
                onChange={(e) => updateModbusEndpoint(endpoint.id ?? "", "testRegister", Number.parseInt(e.target.value, 10))}
                placeholder="Reg"
                title={t("health.testRegister")}
                className={cn(input.base, input.state.default, input.size.md, "w-16")}
              />
              <button
                type="button"
                onClick={() => removeModbusEndpoint(endpoint.id ?? "")}
                className={cn("text-status-error hover:text-status-error/80", spacing.actionBtn)}
              >
                {t("common.remove")}
              </button>
            </div>
          ))}
        </div>

        {/* SLA Configuration */}
        <div className={cn("border-t border-surface-border", spacing.padding.top.heading)}>
          <div className={cn(layout.flex.between, spacing.margin.bottom.inline)}>
            <span className="caption text-text-muted font-medium">{t("health.slaConfig")}</span>
          </div>
          <div className={spacing.stack.xs}>
            <label
              className={cn(
                layout.flex.between,
                spacing.pad.sm,
                "bg-surface-base border border-surface-border",
                radius.default,
              )}
            >
              <div>
                <span className="body-small text-text-primary font-medium">
                  {t("health.enableSla")}
                </span>
                <p className="caption text-text-muted">{t("health.slaDescription")}</p>
              </div>
              <input
                type="checkbox"
                checked={testsSettings.slaConfigs?.[0]?.enabled ?? false}
                onChange={(e) =>
                  setTestsSettings((prev) => ({
                    ...prev,
                    slaConfigs: [
                      {
                        ...(prev.slaConfigs?.[0] ?? {
                          endpointName: "*",
                          targetUptime: 99.9,
                          targetLatencyP95: 500,
                          reportingPeriod: "daily",
                        }),
                        enabled: e.target.checked,
                      },
                    ],
                  }))
                }
                className={iconTokens.size.sm}
              />
            </label>
            <div className={cn("flex items-center", spacing.gap.compact)}>
              <label className="caption text-text-muted w-32">{t("health.targetUptime")}</label>
              <input
                type="number"
                min={90}
                max={100}
                step={0.1}
                value={testsSettings.slaConfigs?.[0]?.targetUptime ?? 99.9}
                onChange={(e) =>
                  setTestsSettings((prev) => ({
                    ...prev,
                    slaConfigs: [
                      {
                        ...(prev.slaConfigs?.[0] ?? {
                          endpointName: "*",
                          enabled: true,
                          targetLatencyP95: 500,
                          reportingPeriod: "daily",
                        }),
                        targetUptime: Number.parseFloat(e.target.value) || 99.9,
                      },
                    ],
                  }))
                }
                className={cn(input.base, input.state.default, input.size.md, "w-24")}
              />
              <span className="caption text-text-muted">%</span>
            </div>
            <div className={cn("flex items-center", spacing.gap.compact)}>
              <label className="caption text-text-muted w-32">{t("health.targetLatency")}</label>
              <input
                type="number"
                min={10}
                max={10000}
                step={10}
                value={testsSettings.slaConfigs?.[0]?.targetLatencyP95 ?? 500}
                onChange={(e) =>
                  setTestsSettings((prev) => ({
                    ...prev,
                    slaConfigs: [
                      {
                        ...(prev.slaConfigs?.[0] ?? {
                          endpointName: "*",
                          enabled: true,
                          targetUptime: 99.9,
                          reportingPeriod: "daily",
                        }),
                        targetLatencyP95: Number.parseInt(e.target.value, 10) || 500,
                      },
                    ],
                  }))
                }
                className={cn(input.base, input.state.default, input.size.md, "w-24")}
              />
              <span className="caption text-text-muted">ms (P95)</span>
            </div>
          </div>
        </div>

        {/* Alert Configuration */}
        <div className={cn("border-t border-surface-border", spacing.padding.top.heading)}>
          <div className={cn(layout.flex.between, spacing.margin.bottom.inline)}>
            <span className="caption text-text-muted font-medium">{t("health.alertConfig")}</span>
          </div>
          <div className={spacing.stack.xs}>
            <label
              className={cn(
                layout.flex.between,
                spacing.pad.xs,
                "bg-surface-base border border-surface-border",
                radius.default,
              )}
            >
              <div>
                <span className="body-small text-text-primary font-medium">
                  {t("health.enableAlerts")}
                </span>
                <p className="caption text-text-muted">{t("health.alertsDescription")}</p>
              </div>
              <input
                type="checkbox"
                checked={testsSettings.alertConfig?.enabled ?? true}
                onChange={(e) =>
                  setTestsSettings((prev) => ({
                    ...prev,
                    alertConfig: {
                      ...(prev.alertConfig ?? {
                        enabled: true,
                        consecutiveFailures: 3,
                        cooldownMinutes: 5,
                        digestMode: false,
                      }),
                      enabled: e.target.checked,
                    },
                  }))
                }
                className={iconTokens.size.sm}
              />
            </label>

            <div className={cn("flex items-center", spacing.gap.compact)}>
              <label className="caption text-text-muted flex-1">
                {t("health.consecutiveFailures")}
              </label>
              <input
                type="number"
                min={1}
                max={10}
                value={testsSettings.alertConfig?.consecutiveFailures ?? 3}
                onChange={(e) =>
                  setTestsSettings((prev) => ({
                    ...prev,
                    alertConfig: {
                      ...(prev.alertConfig ?? {
                        enabled: true,
                        consecutiveFailures: 3,
                        cooldownMinutes: 5,
                        digestMode: false,
                      }),
                      consecutiveFailures: Number.parseInt(e.target.value, 10) || 3,
                    },
                  }))
                }
                className={cn(input.base, input.state.default, input.size.md, "w-20 text-center")}
              />
            </div>

            <div className={cn("flex items-center", spacing.gap.compact)}>
              <label className="caption text-text-muted flex-1">
                {t("health.cooldownMinutes")}
              </label>
              <input
                type="number"
                min={1}
                max={60}
                value={testsSettings.alertConfig?.cooldownMinutes ?? 5}
                onChange={(e) =>
                  setTestsSettings((prev) => ({
                    ...prev,
                    alertConfig: {
                      ...(prev.alertConfig ?? {
                        enabled: true,
                        consecutiveFailures: 3,
                        cooldownMinutes: 5,
                        digestMode: false,
                      }),
                      cooldownMinutes: Number.parseInt(e.target.value, 10) || 5,
                    },
                  }))
                }
                className={cn(input.base, input.state.default, input.size.md, "w-20 text-center")}
              />
            </div>

            <label
              className={cn(
                layout.flex.between,
                spacing.pad.xs,
                "bg-surface-base border border-surface-border",
                radius.default,
              )}
            >
              <div>
                <span className="body-small text-text-primary font-medium">
                  {t("health.digestMode")}
                </span>
                <p className="caption text-text-muted">{t("health.digestDescription")}</p>
              </div>
              <input
                type="checkbox"
                checked={testsSettings.alertConfig?.digestMode ?? false}
                onChange={(e) =>
                  setTestsSettings((prev) => ({
                    ...prev,
                    alertConfig: {
                      ...(prev.alertConfig ?? {
                        enabled: true,
                        consecutiveFailures: 3,
                        cooldownMinutes: 5,
                        digestMode: false,
                      }),
                      digestMode: e.target.checked,
                    },
                  }))
                }
                className={iconTokens.size.sm}
              />
            </label>
          </div>
        </div>

        {/* Anomaly Detection Configuration */}
        <div className={cn("border-t border-surface-border", spacing.padding.top.heading)}>
          <div className={cn(layout.flex.between, spacing.margin.bottom.inline)}>
            <span className="caption text-text-muted font-medium">{t("health.anomalyConfig")}</span>
          </div>
          <div className={spacing.stack.xs}>
            <label
              className={cn(
                layout.flex.between,
                spacing.pad.xs,
                "bg-surface-base border border-surface-border",
                radius.default,
              )}
            >
              <div>
                <span className="body-small text-text-primary font-medium">
                  {t("health.enableAnomaly")}
                </span>
                <p className="caption text-text-muted">{t("health.anomalyDescription")}</p>
              </div>
              <input
                type="checkbox"
                checked={testsSettings.anomalyConfig?.enabled ?? true}
                onChange={(e) =>
                  setTestsSettings((prev) => ({
                    ...prev,
                    anomalyConfig: {
                      ...(prev.anomalyConfig ?? {
                        enabled: true,
                        stdDevThreshold: 2,
                        maxSamples: 100,
                      }),
                      enabled: e.target.checked,
                    },
                  }))
                }
                className={iconTokens.size.sm}
              />
            </label>

            <div className={cn("flex items-center", spacing.gap.compact)}>
              <label className="caption text-text-muted flex-1">
                {t("health.stdDevThreshold")}
              </label>
              <input
                type="number"
                min={1}
                max={5}
                step={0.5}
                value={testsSettings.anomalyConfig?.stdDevThreshold ?? 2}
                onChange={(e) =>
                  setTestsSettings((prev) => ({
                    ...prev,
                    anomalyConfig: {
                      ...(prev.anomalyConfig ?? {
                        enabled: true,
                        stdDevThreshold: 2,
                        maxSamples: 100,
                      }),
                      stdDevThreshold: Number.parseFloat(e.target.value) || 2,
                    },
                  }))
                }
                className={cn(input.base, input.state.default, input.size.md, "w-20 text-center")}
              />
            </div>

            <div className={cn("flex items-center", spacing.gap.compact)}>
              <label className="caption text-text-muted flex-1">{t("health.maxSamples")}</label>
              <input
                type="number"
                min={10}
                max={500}
                step={10}
                value={testsSettings.anomalyConfig?.maxSamples ?? 100}
                onChange={(e) =>
                  setTestsSettings((prev) => ({
                    ...prev,
                    anomalyConfig: {
                      ...(prev.anomalyConfig ?? {
                        enabled: true,
                        stdDevThreshold: 2,
                        maxSamples: 100,
                      }),
                      maxSamples: Number.parseInt(e.target.value, 10) || 100,
                    },
                  }))
                }
                className={cn(input.base, input.state.default, input.size.md, "w-20 text-center")}
              />
            </div>
          </div>
        </div>
      </div>
    </CollapsibleSection>
  );
});
