/**
 * SurveyView Component (~539 lines)
 *
 * Purpose: Full-featured WiFi site survey editor and viewer. Allows users to create
 * detailed signal strength and performance heatmaps by recording samples at physical
 * locations on a floor plan, supporting passive scanning, active testing, and iperf3.
 *
 * Key Features:
 * - Floor plan canvas: Interactive image-based coordinate system for sample placement
 * - Floor plan upload: Users can upload custom floor plan images
 * - Passive sampling: Record RSSI (signal strength) measurements from WiFi beacons
 * - Active sampling: Measure latency and packet loss to target server
 * - Throughput sampling: Use iperf3 server for bandwidth measurements
 * - Heatmap visualization: Color-coded overlays showing signal/performance gradients
 * - Sample management: Create, update, delete sample points
 * - Export data: Save survey results for analysis
 * - Real-time updates: Reflect changes immediately in UI
 *
 * Usage:
 * ```typescript
 * <SurveyView
 *   survey={surveyData}
 *   onClose={handleClose}
 *   onUpdate={handleUpdate}
 * />
 * ```
 *
 * Dependencies: FloorPlanCanvas, useAuth, useSurvey hook, API communication
 * State: survey data, sampling status, heatmap metric selection, upload progress
 */

import { useState, useEffect, useCallback } from "react";
import { useTranslation } from "react-i18next";
import { FloorPlanCanvas, type CalibrationPoint } from "./FloorPlanCanvas";
import { getAuthHeaders } from "../../hooks/useAuth";
import type { Survey, PassiveSample, ActiveSample, ThroughputSample } from "../../hooks/useSurvey";
import { X, Upload, Play, Pause, CheckCircle, Loader } from "../ui/Icons";
// Import Ruler directly from lucide-react since it's not in our Icons module yet
import { Ruler } from "lucide-react";
import { radius, layout, spacing, button, icon as iconTokens } from "../../styles/theme";

const API_BASE = import.meta.env.VITE_API_BASE || "";

interface SurveyViewProps {
  survey: Survey;
  onClose: () => void;
  onUpdate: () => void;
}

/**
 * SurveyView Component
 * Main survey interface with floor plan, sampling controls, and heatmap visualization
 */
// WiFi adapter status from /api/wifi/status
interface WiFiStatus {
  status: "unavailable" | "available" | "ready";
  message: string;
  currentInterface: string;
  isWireless: boolean;
  availableAdapters: string[];
  canScan: boolean;
}

/**
 *
 */
export function SurveyView({ survey: initialSurvey, onClose, onUpdate }: SurveyViewProps) {
  const { t } = useTranslation("survey");
  // Current survey being edited
  const [survey, setSurvey] = useState(initialSurvey);
  // Indicates if a sampling operation is in progress
  const [sampling, setSampling] = useState(false);
  // Indicates if floor plan upload is in progress
  const [uploadingFloorPlan, setUploadingFloorPlan] = useState(false);
  // Selected metric for heatmap visualization (rssi, throughput, latency)
  const [heatmapMetric, setHeatmapMetric] = useState<"rssi" | "throughput" | "latency" | null>(
    null
  );
  const [error, setError] = useState<string | null>(null);
  // WiFi adapter status
  const [wifiStatus, setWifiStatus] = useState<WiFiStatus | null>(null);
  // Calibration mode state
  const [calibrationMode, setCalibrationMode] = useState(false);
  const [calibrationPoints, setCalibrationPoints] = useState<CalibrationPoint[]>([]);
  const [calibrationDistance, setCalibrationDistance] = useState<string>("");
  // Survey settings edit state
  const [editSurveyType, setEditSurveyType] = useState(initialSurvey.surveyType);
  const [editIperfServer, setEditIperfServer] = useState(initialSurvey.iperfServer || "");
  const [editTestDuration, setEditTestDuration] = useState(initialSurvey.testDuration || 3);
  const [savingSettings, setSavingSettings] = useState(false);

  // Check WiFi adapter status on mount
  useEffect(() => {
    const checkWifiStatus = async () => {
      try {
        const res = await fetch(`${API_BASE}/api/wifi/status`, {
          headers: getAuthHeaders(),
        });
        if (res.ok) {
          const status = await res.json();
          setWifiStatus(status);
        }
      } catch (err) {
        console.error("Failed to check WiFi status:", err);
      }
    };
    checkWifiStatus();
  }, []);

  // Poll for survey updates when in progress
  useEffect(() => {
    if (survey.status !== "in_progress") return;

    const interval = setInterval(async () => {
      try {
        const res = await fetch(`${API_BASE}/api/survey?id=${survey.id}`, {
          headers: getAuthHeaders(),
        });
        if (res.ok) {
          const updated = await res.json();
          setSurvey(updated);
        }
      } catch (err) {
        console.error("Failed to refresh survey:", err);
      }
    }, 3000);

    return () => clearInterval(interval);
  }, [survey.id, survey.status]);

  // Handle floor plan upload
  const handleFloorPlanUpload = useCallback(
    async (file: File) => {
      setUploadingFloorPlan(true);
      setError(null);

      try {
        // Read file as base64 using Promise wrapper
        const imageData = await new Promise<string>((resolve, reject) => {
          const reader = new FileReader();
          reader.onload = (e) => {
            const result = e.target?.result;
            if (typeof result === "string") {
              resolve(result);
            } else {
              reject(new Error("Failed to read file as base64"));
            }
          };
          reader.onerror = () => {
            reject(new Error("Failed to read file"));
          };
          reader.readAsDataURL(file);
        });

        // Get image dimensions using Promise wrapper
        const { width, height } = await new Promise<{ width: number; height: number }>(
          (resolve, reject) => {
            const img = new Image();
            img.onload = () => {
              resolve({ width: img.width, height: img.height });
            };
            img.onerror = () => {
              reject(new Error("Failed to load image - file may be corrupted"));
            };
            img.src = imageData;
          }
        );

        const floorPlan = {
          imageData,
          width,
          height,
          scaleM: 0.1, // Default: 10cm per pixel (adjust in settings)
        };

        // Upload to server
        const res = await fetch(`${API_BASE}/api/survey/floorplan?id=${survey.id}`, {
          method: "POST",
          headers: {
            ...getAuthHeaders(),
            "Content-Type": "application/json",
          },
          body: JSON.stringify(floorPlan),
        });

        if (!res.ok) {
          const errorText = await res.text();
          throw new Error(errorText || "Failed to upload floor plan");
        }

        // Refresh survey
        const updated = await res.json();
        setSurvey(updated);
        onUpdate();
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to upload floor plan");
      } finally {
        setUploadingFloorPlan(false);
      }
    },
    [survey.id, onUpdate]
  );

  // Handle taking a sample at clicked location
  const handlePointClick = useCallback(
    async (x: number, y: number) => {
      if (survey.status !== "in_progress") return;

      // Check WiFi availability before sampling
      if (!wifiStatus?.canScan) {
        setError(wifiStatus?.message || "No WiFi adapter available for scanning");
        return;
      }

      setSampling(true);
      setError(null);

      try {
        // Collect sample data based on survey type
        let sampleData: PassiveSample | ActiveSample | ThroughputSample;

        switch (survey.surveyType) {
          case "passive": {
            // Fetch WiFi scan
            const scanRes = await fetch(`${API_BASE}/api/wifi/scan`, {
              headers: getAuthHeaders(),
            });
            if (!scanRes.ok) throw new Error("WiFi scan failed");
            const scanData = await scanRes.json();
            // Check if scan was successful
            if (!scanData.available) {
              throw new Error(scanData.error || "WiFi scan not available");
            }
            sampleData = { networks: scanData.networks || [] };
            break;
          }

          case "active": {
            // Fetch current WiFi status
            const wifiRes = await fetch(`${API_BASE}/api/wifi`, {
              headers: getAuthHeaders(),
            });
            if (!wifiRes.ok) throw new Error("WiFi status fetch failed");
            const wifiData = await wifiRes.json();

            // Check if BSSID changed (roaming)
            const samples = survey.samples ?? [];
            const lastSample = samples[samples.length - 1];
            const lastBssid = lastSample ? (lastSample.sampleData as ActiveSample).bssid : null;
            const roamingEvent = lastBssid !== null && lastBssid !== wifiData.bssid;

            sampleData = {
              ssid: wifiData.ssid || "",
              bssid: wifiData.bssid || "",
              rssi: wifiData.signal || 0,
              dataRate: wifiData.bitrate || 0,
              roamingEvent,
            };
            break;
          }

          case "throughput": {
            // Fetch WiFi status first
            const wifiRes2 = await fetch(`${API_BASE}/api/wifi`, {
              headers: getAuthHeaders(),
            });
            if (!wifiRes2.ok) throw new Error("WiFi status fetch failed");
            const wifiData2 = await wifiRes2.json();

            // Run iperf3 test
            if (!survey.iperfServer) {
              throw new Error("iperf3 server not configured for this survey");
            }

            const [host, port] = survey.iperfServer.split(":");
            const iperfRes = await fetch(`${API_BASE}/api/iperf/client`, {
              method: "POST",
              headers: {
                ...getAuthHeaders(),
                "Content-Type": "application/json",
              },
              body: JSON.stringify({
                host,
                port: port ? parseInt(port) : 5201,
                duration: survey.testDuration || 3,
                reverse: false,
              }),
            });

            if (!iperfRes.ok) throw new Error("iperf3 test failed");
            const iperfData = await iperfRes.json();

            sampleData = {
              ssid: wifiData2.ssid || "",
              bssid: wifiData2.bssid || "",
              rssi: wifiData2.signal || 0,
              downloadMbps: iperfData.summary?.sum_received?.bits_per_second
                ? iperfData.summary.sum_received.bits_per_second / 1_000_000
                : 0,
              uploadMbps: iperfData.summary?.sum_sent?.bits_per_second
                ? iperfData.summary.sum_sent.bits_per_second / 1_000_000
                : 0,
              latency: 0, // Not available in standard iperf3
              jitter: iperfData.summary?.sum_received?.jitter_ms || 0,
              packetLoss: iperfData.summary?.sum_received?.lost_percent || 0,
            };
            break;
          }

          default:
            throw new Error("Unknown survey type");
        }

        // Submit sample to server
        const res = await fetch(`${API_BASE}/api/survey/sample?id=${survey.id}`, {
          method: "POST",
          headers: {
            ...getAuthHeaders(),
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ x, y, sampleData }),
        });

        if (!res.ok) throw new Error("Failed to save sample");

        // Refresh survey to get updated samples
        const refreshRes = await fetch(`${API_BASE}/api/survey?id=${survey.id}`, {
          headers: getAuthHeaders(),
        });
        if (refreshRes.ok) {
          const updated = await refreshRes.json();
          setSurvey(updated);
          onUpdate();
        }
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to take sample");
      } finally {
        setSampling(false);
      }
    },
    [survey, onUpdate, wifiStatus]
  );

  // Handle calibration click - collect two points
  const handleCalibrationClick = useCallback((x: number, y: number) => {
    setCalibrationPoints((prev) => {
      if (prev.length >= 2) {
        // Reset if we already have 2 points
        return [{ x, y }];
      }
      return [...prev, { x, y }];
    });
  }, []);

  // Calculate and save scale from calibration
  const handleSaveCalibration = async () => {
    if (calibrationPoints.length !== 2 || !calibrationDistance) {
      setError("Please select two points and enter the distance");
      return;
    }

    const distance = parseFloat(calibrationDistance);
    if (isNaN(distance) || distance <= 0) {
      setError("Please enter a valid positive distance");
      return;
    }

    // Calculate pixel distance
    const p1 = calibrationPoints[0];
    const p2 = calibrationPoints[1];
    const pixelDist = Math.sqrt((p2.x - p1.x) ** 2 + (p2.y - p1.y) ** 2);

    if (pixelDist === 0) {
      setError("Please select two different points");
      return;
    }

    // Calculate scale (meters per pixel)
    const scaleM = distance / pixelDist;

    try {
      // Update floor plan scale on server
      const res = await fetch(`${API_BASE}/api/survey/floorplan?id=${survey.id}`, {
        method: "POST",
        headers: {
          ...getAuthHeaders(),
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          ...survey.floorPlan,
          scaleM,
        }),
      });

      if (!res.ok) {
        throw new Error("Failed to update floor plan scale");
      }

      const updated = await res.json();
      setSurvey(updated);
      onUpdate();

      // Exit calibration mode
      setCalibrationMode(false);
      setCalibrationPoints([]);
      setCalibrationDistance("");
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to save scale");
    }
  };

  // Cancel calibration
  const handleCancelCalibration = () => {
    setCalibrationMode(false);
    setCalibrationPoints([]);
    setCalibrationDistance("");
  };

  // Save survey settings
  const handleSaveSettings = async () => {
    setSavingSettings(true);
    setError(null);

    try {
      const res = await fetch(`${API_BASE}/api/survey/settings?id=${survey.id}`, {
        method: "PUT",
        headers: {
          ...getAuthHeaders(),
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          surveyType: editSurveyType,
          iperfServer: editIperfServer,
          testDuration: editTestDuration,
        }),
      });

      if (!res.ok) {
        const errorText = await res.text();
        throw new Error(errorText || "Failed to save settings");
      }

      const updated = await res.json();
      setSurvey(updated);
      onUpdate();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to save settings");
    } finally {
      setSavingSettings(false);
    }
  };

  // Handle status changes
  const handleStatusChange = async (action: "start" | "pause" | "complete") => {
    try {
      const res = await fetch(`${API_BASE}/api/survey/${action}?id=${survey.id}`, {
        method: "POST",
        headers: getAuthHeaders(),
      });

      if (res.ok) {
        const updated = await res.json();
        setSurvey(updated);
        onUpdate();
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : `Failed to ${action} survey`);
    }
  };

  return (
    <div className="fixed inset-0 bg-surface-base z-50 overflow-auto">
      {/* Header */}
      <div className="sticky top-0 bg-surface-raised border-b border-surface-border z-10">
        <div className={`max-w-7xl mx-auto pad ${layout.flex.between}`}>
          <div>
            <h1 className="heading-1">{survey.name}</h1>
            <p className={`body-small ${spacing.margin.top.tight}`}>
              {(survey.surveyType ?? "wifi").charAt(0).toUpperCase() +
                (survey.surveyType ?? "wifi").slice(1)}{" "}
              {t("status.survey")} • {(survey.samples ?? []).length} {t("status.samples")} •{" "}
              {survey.status ?? "unknown"}
            </p>
          </div>

          <div className={`${layout.inline.default}`}>
            {/* Status controls */}
            {survey.status === "created" && (
              <button
                onClick={() => handleStatusChange("start")}
                className={`${button.size.md} bg-brand-primary text-text-inverse ${radius.md} hover:bg-brand-primary/90 ${layout.inline.default}`}
              >
                <Play className={iconTokens.size.sm} />
                {t("buttons.startSurvey")}
              </button>
            )}

            {survey.status === "in_progress" && (
              <>
                <button
                  onClick={() => handleStatusChange("pause")}
                  className={`${button.size.md} border border-surface-border ${radius.md} hover:bg-surface-hover ${layout.inline.default}`}
                >
                  <Pause className={iconTokens.size.sm} />
                  {t("buttons.pause")}
                </button>
                <button
                  onClick={() => handleStatusChange("complete")}
                  className={`${button.size.md} bg-status-success text-text-inverse ${radius.md} hover:bg-status-success/90 ${layout.inline.default}`}
                >
                  <CheckCircle className={iconTokens.size.sm} />
                  {t("buttons.complete")}
                </button>
              </>
            )}

            {survey.status === "paused" && (
              <>
                <button
                  onClick={() => handleStatusChange("start")}
                  className={`${button.size.md} bg-brand-primary text-text-inverse ${radius.md} hover:bg-brand-primary/90 ${layout.inline.default}`}
                >
                  <Play className={iconTokens.size.sm} />
                  {t("buttons.resume")}
                </button>
                <button
                  onClick={() => handleStatusChange("complete")}
                  className={`${button.size.md} bg-status-success text-text-inverse ${radius.md} hover:bg-status-success/90 ${layout.inline.default}`}
                >
                  <CheckCircle className={iconTokens.size.sm} />
                  {t("buttons.complete")}
                </button>
              </>
            )}

            <button
              onClick={onClose}
              className={`${button.size.md} border border-surface-border ${radius.md} hover:bg-surface-hover ${layout.inline.default}`}
            >
              <X className={iconTokens.size.sm} />
              {t("buttons.close")}
            </button>
          </div>
        </div>
      </div>

      {/* Main content */}
      <div className={`max-w-7xl mx-auto ${spacing.pad.default} ${spacing.pad.lg}`}>
        {error && (
          <div
            className={`bg-status-error/10 border border-status-error/20 text-status-error ${spacing.pad.sm} ${radius.md} ${spacing.margin.bottom.content}`}
          >
            {error}
          </div>
        )}

        {/* WiFi adapter status banner */}
        {wifiStatus && wifiStatus.status !== "ready" && (
          <div
            className={`${
              wifiStatus.status === "unavailable"
                ? "bg-status-warning/10 border-status-warning/20 text-status-warning"
                : "bg-status-info/10 border-status-info/20 text-status-info"
            } border ${spacing.pad.sm} ${radius.md} ${spacing.margin.bottom.content}`}
          >
            <div className="font-medium">
              {wifiStatus.status === "unavailable" ? "⚠️ " : "ℹ️ "}
              {wifiStatus.status === "unavailable"
                ? t("wifi.noAdapter")
                : t("wifi.adapterAvailable")}
            </div>
            {wifiStatus.availableAdapters.length > 0 && (
              <div className={`caption ${spacing.margin.top.tight}`}>
                {t("wifi.availableAdapters")}: {wifiStatus.availableAdapters.join(", ")}
              </div>
            )}
            <div className={`caption ${spacing.margin.top.tight}`}>{t("wifi.canStillUpload")}</div>
          </div>
        )}

        {sampling && (
          <div
            className={`bg-status-info/10 border border-status-info/20 text-status-info ${spacing.pad.sm} ${radius.md} ${spacing.margin.bottom.content} ${layout.inline.default}`}
          >
            <Loader className={`${iconTokens.size.sm} animate-spin`} />
            {t("progress.takingMeasurement")}
          </div>
        )}

        <div className={`grid grid-cols-1 lg:grid-cols-3 ${spacing.gap.spacious}`}>
          {/* Floor plan */}
          <div className="lg:col-span-2">
            <div className={`bg-surface-raised ${radius.md} border border-surface-border pad`}>
              <div className={`${layout.flex.between} ${spacing.margin.bottom.content}`}>
                <h2 className="heading-3">{t("floorPlan.title")}</h2>
                {heatmapMetric === null && (survey.samples ?? []).length > 0 && (
                  <div className={layout.inline.default}>
                    <button
                      onClick={() => setHeatmapMetric("rssi")}
                      className={`${button.size.sm} body-small border border-surface-border ${radius.md} hover:bg-surface-hover`}
                    >
                      {t("buttons.rssiHeatmap")}
                    </button>
                    {survey.surveyType === "throughput" && (
                      <>
                        <button
                          onClick={() => setHeatmapMetric("throughput")}
                          className={`${button.size.sm} body-small border border-surface-border ${radius.md} hover:bg-surface-hover`}
                        >
                          {t("buttons.throughput")}
                        </button>
                        <button
                          onClick={() => setHeatmapMetric("latency")}
                          className={`${button.size.sm} body-small border border-surface-border ${radius.md} hover:bg-surface-hover`}
                        >
                          {t("buttons.latency")}
                        </button>
                      </>
                    )}
                  </div>
                )}
                {heatmapMetric !== null && (
                  <button
                    onClick={() => setHeatmapMetric(null)}
                    className={`${button.size.sm} body-small bg-brand-primary text-text-inverse ${radius.md} hover:bg-brand-primary/90`}
                  >
                    {t("buttons.hideHeatmap")}
                  </button>
                )}
              </div>

              {!survey.floorPlan ? (
                <div
                  className={`border-2 border-dashed border-surface-border ${radius.md} pad-lg text-center`}
                >
                  <Upload
                    className={`${iconTokens.size.xl} mx-auto text-text-muted ${spacing.margin.bottom.content}`}
                  />
                  <p className={`text-text-muted ${spacing.margin.bottom.content}`}>
                    {t("floorPlan.uploadPrompt")}
                  </p>
                  <label
                    className={`inline-block ${button.size.md} bg-brand-primary text-text-inverse ${radius.md} cursor-pointer hover:bg-brand-primary/90`}
                  >
                    {uploadingFloorPlan ? t("floorPlan.uploading") : t("floorPlan.chooseFile")}
                    <input
                      type="file"
                      accept="image/png,image/jpeg,image/gif,image/webp,image/svg+xml"
                      className="hidden"
                      onChange={(e) => {
                        const file = e.target.files?.[0];
                        if (file) {
                          handleFloorPlanUpload(file);
                        }
                        // Reset input so same file can be selected again if needed
                        e.target.value = "";
                      }}
                      disabled={uploadingFloorPlan}
                    />
                  </label>
                  <p className={`caption text-text-muted ${spacing.margin.top.inline}`}>
                    {t("floorPlan.supportedFormats")}
                  </p>
                </div>
              ) : (
                <div>
                  {/* Calibration panel */}
                  {calibrationMode && (
                    <div
                      className={`bg-status-warning/10 border border-status-warning/20 ${spacing.pad.sm} ${radius.md} ${spacing.margin.bottom.content}`}
                    >
                      <div
                        className={`font-medium text-status-warning ${spacing.margin.bottom.inline}`}
                      >
                        📐 {t("calibration.title")}
                      </div>
                      <p
                        className={`body-small text-text-secondary ${spacing.margin.bottom.content}`}
                      >
                        {t("calibration.instructions")}
                      </p>
                      <div className="stack-sm">
                        <div className={`${layout.inline.default}`}>
                          <span className="body-small text-text-muted w-20">
                            {t("calibration.pointA")}:
                          </span>
                          {calibrationPoints[0] ? (
                            <span className="body-small font-medium">
                              ({calibrationPoints[0].x}, {calibrationPoints[0].y})
                            </span>
                          ) : (
                            <span className="body-small text-text-muted italic">
                              {t("calibration.clickFloorPlan")}
                            </span>
                          )}
                        </div>
                        <div className={`${layout.inline.default}`}>
                          <span className="body-small text-text-muted w-20">
                            {t("calibration.pointB")}:
                          </span>
                          {calibrationPoints[1] ? (
                            <span className="body-small font-medium">
                              ({calibrationPoints[1].x}, {calibrationPoints[1].y})
                            </span>
                          ) : (
                            <span className="body-small text-text-muted italic">
                              {t("calibration.clickFloorPlan")}
                            </span>
                          )}
                        </div>
                        {calibrationPoints.length === 2 && (
                          <div className={`${layout.inline.default}`}>
                            <span className="body-small text-text-muted w-20">
                              {t("calibration.pixelDistance")}:
                            </span>
                            <span className="body-small font-medium">
                              {Math.sqrt(
                                (calibrationPoints[1].x - calibrationPoints[0].x) ** 2 +
                                  (calibrationPoints[1].y - calibrationPoints[0].y) ** 2
                              ).toFixed(0)}{" "}
                              px
                            </span>
                          </div>
                        )}
                        <div className={`${layout.inline.default} ${spacing.margin.top.inline}`}>
                          <label className="body-small text-text-muted w-20">
                            {t("calibration.distance")}:
                          </label>
                          <input
                            type="number"
                            step="0.1"
                            min="0"
                            value={calibrationDistance}
                            onChange={(e) => setCalibrationDistance(e.target.value)}
                            placeholder={t("calibration.enterMeters")}
                            className={`flex-1 ${button.size.sm} border border-surface-border ${radius.md} bg-surface-base text-text-primary`}
                          />
                          <span className="body-small text-text-muted">
                            {t("calibration.meters")}
                          </span>
                        </div>
                        <div className={`${layout.inline.default} ${spacing.margin.top.inline}`}>
                          <button
                            onClick={handleSaveCalibration}
                            disabled={calibrationPoints.length !== 2 || !calibrationDistance}
                            className={`${button.size.sm} bg-brand-primary text-text-inverse ${radius.md} hover:bg-brand-primary/90 disabled:opacity-50 disabled:cursor-not-allowed`}
                          >
                            {t("buttons.saveScale")}
                          </button>
                          <button
                            onClick={handleCancelCalibration}
                            className={`${button.size.sm} border border-surface-border ${radius.md} hover:bg-surface-hover`}
                          >
                            {t("buttons.cancel")}
                          </button>
                          <button
                            onClick={() => setCalibrationPoints([])}
                            className={`${button.size.sm} border border-surface-border ${radius.md} hover:bg-surface-hover`}
                          >
                            {t("buttons.resetPoints")}
                          </button>
                        </div>
                      </div>
                    </div>
                  )}

                  {/* Calibrate button and current scale info */}
                  {!calibrationMode && survey.floorPlan && (
                    <div className={`${layout.flex.between} ${spacing.margin.bottom.inline}`}>
                      <div className="body-small text-text-muted">
                        {t("floorPlan.scale")}: {survey.floorPlan.scaleM.toFixed(3)} m/px
                        {survey.status === "in_progress" && ` • ${t("floorPlan.clickToMeasure")}`}
                      </div>
                      <button
                        onClick={() => setCalibrationMode(true)}
                        className={`${button.size.sm} body-small border border-surface-border ${radius.md} hover:bg-surface-hover ${layout.inline.tight}`}
                      >
                        <Ruler className={iconTokens.size.sm} />
                        {t("buttons.calibrateScale")}
                      </button>
                    </div>
                  )}

                  <FloorPlanCanvas
                    floorPlan={survey.floorPlan}
                    samples={survey.samples ?? []}
                    onPointClick={handlePointClick}
                    interactive={survey.status === "in_progress" && !sampling && !calibrationMode}
                    heatmapMetric={heatmapMetric}
                    calibrationMode={calibrationMode}
                    calibrationPoints={calibrationPoints}
                    onCalibrationClick={handleCalibrationClick}
                  />
                </div>
              )}
            </div>
          </div>

          {/* Settings panel (shown when survey is in created status) and Sample list */}
          <div className={`lg:col-span-1 ${spacing.stack.default}`}>
            {/* Survey Settings Panel - only show when survey hasn't started */}
            {survey.status === "created" && (
              <div className={`bg-surface-raised ${radius.md} border border-surface-border pad`}>
                <h2 className={`heading-3 ${spacing.margin.bottom.content}`}>
                  {t("settings.title")}
                </h2>
                <div className="stack">
                  {/* Survey Type */}
                  <div>
                    <label
                      className={`body-small text-text-muted block ${spacing.margin.bottom.tight}`}
                    >
                      {t("settings.surveyType")}
                    </label>
                    <select
                      value={editSurveyType}
                      onChange={(e) =>
                        setEditSurveyType(e.target.value as "passive" | "active" | "throughput")
                      }
                      className={`w-full ${button.size.md} border border-surface-border ${radius.md} bg-surface-base text-text-primary`}
                    >
                      <option value="passive">{t("settings.types.passive")}</option>
                      <option value="active">{t("settings.types.active")}</option>
                      <option value="throughput">{t("settings.types.throughput")}</option>
                    </select>
                  </div>

                  {/* iperf Server - only show for throughput surveys */}
                  {editSurveyType === "throughput" && (
                    <>
                      <div>
                        <label
                          className={`body-small text-text-muted block ${spacing.margin.bottom.tight}`}
                        >
                          {t("settings.iperfServer")}
                        </label>
                        <input
                          type="text"
                          value={editIperfServer}
                          onChange={(e) => setEditIperfServer(e.target.value)}
                          placeholder="hostname:5201"
                          className={`w-full ${button.size.md} border border-surface-border ${radius.md} bg-surface-base text-text-primary`}
                        />
                        <p className={`caption text-text-muted ${spacing.margin.top.tight}`}>
                          {t("settings.iperfServerHint")}
                        </p>
                      </div>

                      <div>
                        <label
                          className={`body-small text-text-muted block ${spacing.margin.bottom.tight}`}
                        >
                          {t("settings.testDuration")}
                        </label>
                        <input
                          type="number"
                          min="1"
                          max="60"
                          value={editTestDuration}
                          onChange={(e) => setEditTestDuration(parseInt(e.target.value) || 3)}
                          className={`w-full ${button.size.md} border border-surface-border ${radius.md} bg-surface-base text-text-primary`}
                        />
                      </div>
                    </>
                  )}

                  {/* Save button */}
                  <button
                    onClick={handleSaveSettings}
                    disabled={savingSettings}
                    className={`w-full ${button.size.md} bg-brand-primary text-text-inverse ${radius.md} hover:bg-brand-primary/90 disabled:opacity-50`}
                  >
                    {savingSettings ? t("buttons.saving") : t("buttons.saveSettings")}
                  </button>

                  {/* Survey type descriptions */}
                  <div
                    className={`caption text-text-muted border-t border-surface-border ${spacing.padding.top.section} ${spacing.margin.top.inline}`}
                  >
                    <p className={`font-medium ${spacing.margin.bottom.inline}`}>
                      {t("settings.typesDescription")}
                    </p>
                    <ul className={`list-disc list-inside ${spacing.stack.xs}`}>
                      <li>
                        <strong>Passive:</strong> {t("settings.passiveDesc")}
                      </li>
                      <li>
                        <strong>Active:</strong> {t("settings.activeDesc")}
                      </li>
                      <li>
                        <strong>Throughput:</strong> {t("settings.throughputDesc")}
                      </li>
                    </ul>
                  </div>
                </div>
              </div>
            )}

            {/* Sample list */}
            <div className={`bg-surface-raised ${radius.md} border border-surface-border pad`}>
              <h2 className={`heading-3 ${spacing.margin.bottom.content}`}>
                {t("samples.title")} ({(survey.samples ?? []).length})
              </h2>
              <div className="stack-sm max-h-[70vh] overflow-y-auto">
                {(survey.samples ?? []).length === 0 ? (
                  <p className={`body-small text-center ${spacing.pad.lg}`}>
                    {t("samples.noSamples")}{" "}
                    {survey.status === "in_progress"
                      ? t("samples.clickToStart")
                      : t("samples.startToBegin")}
                  </p>
                ) : (
                  (survey.samples ?? []).map((sample, idx) => (
                    <div
                      key={idx}
                      className={`border border-surface-border ${radius.md} pad-sm body-small`}
                    >
                      <div
                        className={`flex items-center justify-between ${spacing.margin.bottom.inline}`}
                      >
                        <span className="font-semibold">#{idx + 1}</span>
                        <span className="caption">
                          {new Date(sample.timestamp).toLocaleTimeString()}
                        </span>
                      </div>
                      <div className="caption stack-xs">
                        {renderSampleData(sample.sampleData, survey.surveyType)}
                      </div>
                    </div>
                  ))
                )}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

// Helper to render sample data
function renderSampleData(
  data: PassiveSample | ActiveSample | ThroughputSample,
  surveyType: string
) {
  if (surveyType === "passive") {
    const passiveData = data as PassiveSample;
    return (
      <>
        <div>Networks: {passiveData.networks?.length || 0}</div>
        {passiveData.networks?.[0] && (
          <>
            <div>Strongest: {passiveData.networks[0].ssid}</div>
            <div>RSSI: {passiveData.networks[0].rssi} dBm</div>
          </>
        )}
      </>
    );
  }

  if (surveyType === "active") {
    const activeData = data as ActiveSample;
    return (
      <>
        <div>SSID: {activeData.ssid}</div>
        <div>RSSI: {activeData.rssi} dBm</div>
        <div>Rate: {activeData.dataRate} Mbps</div>
        {activeData.roamingEvent && (
          <div className="text-status-warning font-semibold">⚠ Roaming</div>
        )}
      </>
    );
  }

  if (surveyType === "throughput") {
    const throughputData = data as ThroughputSample;
    return (
      <>
        <div>RSSI: {throughputData.rssi} dBm</div>
        <div>↓ {throughputData.downloadMbps.toFixed(1)} Mbps</div>
        <div>↑ {throughputData.uploadMbps.toFixed(1)} Mbps</div>
        <div>Jitter: {throughputData.jitter.toFixed(1)} ms</div>
        {throughputData.packetLoss > 0 && (
          <div className="text-status-error">Loss: {throughputData.packetLoss.toFixed(1)}%</div>
        )}
      </>
    );
  }

  return null;
}
