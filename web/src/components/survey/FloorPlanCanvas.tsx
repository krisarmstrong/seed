/**
 * FloorPlanCanvas Component (~234 lines)
 *
 * Purpose: HTML5 Canvas-based floor plan visualization for WiFi survey heatmaps.
 * Renders floor plan images with overlaid sample points and heatmap color gradients
 * showing signal strength or performance metrics at each location.
 *
 * Key Features:
 * - Image rendering: Displays floor plan images maintaining aspect ratio
 * - Sample markers: Shows measurement points as circles on the floor plan
 * - Heatmap visualization: Color-coded overlays (red=weak, yellow=medium, green=strong)
 * - Interactive mode: Click to add new sample points
 * - Responsive sizing: Automatically adapts to container width
 * - Metric selection: Display RSSI, throughput, or latency heatmaps
 * - Interpolation: Smooth color gradient between sample points
 * - Canvas optimization: Efficient rendering with minimal redraws
 *
 * Usage:
 * ```typescript
 * <FloorPlanCanvas
 *   floorPlan={floorPlanData}
 *   samples={samplePoints}
 *   onPointClick={handleClick}
 *   interactive={true}
 *   heatmapMetric="rssi"
 * />
 * ```
 *
 * Dependencies: Canvas API, floor plan and sample point types
 * State: Canvas dimensions, rendering state
 */

import { useRef, useEffect, useState } from "react";
import { radius } from "../../styles/theme";
import type { FloorPlan, SamplePoint } from "../../hooks/useSurvey";

export interface CalibrationPoint {
  x: number;
  y: number;
}

interface FloorPlanCanvasProps {
  floorPlan: FloorPlan;
  samples: SamplePoint[];
  onPointClick?: (x: number, y: number) => void;
  interactive?: boolean;
  heatmapMetric?: "rssi" | "throughput" | "latency" | null;
  calibrationMode?: boolean;
  calibrationPoints?: CalibrationPoint[];
  onCalibrationClick?: (x: number, y: number) => void;
}

/**
 * FloorPlanCanvas Component
 * Renders floor plan with interactive sample points and heatmap visualization
 */
export function FloorPlanCanvas({
  floorPlan,
  samples,
  onPointClick,
  interactive = false,
  heatmapMetric = null,
  calibrationMode = false,
  calibrationPoints = [],
  onCalibrationClick,
}: FloorPlanCanvasProps) {
  // Canvas DOM reference for drawing
  const canvasRef = useRef<HTMLCanvasElement>(null);
  // Container reference to measure available space
  const containerRef = useRef<HTMLDivElement>(null);
  // Track canvas dimensions for responsive sizing
  const [dimensions, setDimensions] = useState({ width: 0, height: 0 });

  // Calculate canvas dimensions maintaining aspect ratio
  useEffect(() => {
    if (!containerRef.current || !floorPlan) return;

    const container = containerRef.current;
    const containerWidth = container.clientWidth;
    const aspectRatio = floorPlan.height / floorPlan.width;

    const width = Math.min(containerWidth, 1200);
    const height = width * aspectRatio;

    setDimensions({ width, height });
  }, [floorPlan]);

  // Draw floor plan and samples
  useEffect(() => {
    if (!canvasRef.current || !floorPlan || dimensions.width === 0) return;

    const canvas = canvasRef.current;
    const ctx = canvas.getContext("2d");
    if (!ctx) return;

    // Set canvas size
    canvas.width = dimensions.width;
    canvas.height = dimensions.height;

    // Draw floor plan image
    const img = new Image();
    img.onload = () => {
      ctx.drawImage(img, 0, 0, dimensions.width, dimensions.height);

      // Calculate scale factor
      const scaleX = dimensions.width / floorPlan.width;
      const scaleY = dimensions.height / floorPlan.height;

      // Draw heatmap if requested
      if (heatmapMetric && samples.length > 0) {
        drawHeatmap(ctx, samples, heatmapMetric, scaleX, scaleY);
      }

      // Draw sample points
      samples.forEach((sample) => {
        const x = sample.x * scaleX;
        const y = sample.y * scaleY;

        // Draw point
        // Note: Canvas API requires direct color values - CSS variables don't work
        // These colors are intentionally hardcoded for Canvas compatibility
        // See: https://developer.mozilla.org/en-US/docs/Web/API/CanvasRenderingContext2D/fillStyle
        ctx.beginPath();
        ctx.arc(x, y, 8, 0, 2 * Math.PI);
        ctx.fillStyle = heatmapMetric
          ? "rgba(255, 255, 255, 0.8)" // white for visibility on heatmap
          : "rgba(37, 99, 235, 0.8)"; // brand-primary (#2563eb at 80% opacity)
        ctx.fill();
        ctx.strokeStyle = "rgba(255, 255, 255, 1)"; // white border for visibility
        ctx.lineWidth = 2;
        ctx.stroke();

        // Draw point number - high contrast text for visibility
        // Canvas API limitation: must use direct color values
        ctx.fillStyle = heatmapMetric ? "#1e293b" : "#f8fafc"; // slate-800 / slate-50
        ctx.font = "bold 10px sans-serif";
        ctx.textAlign = "center";
        ctx.textBaseline = "middle";
        const pointNum = samples.indexOf(sample) + 1;
        ctx.fillText(pointNum.toString(), x, y);
      });

      // Draw calibration points and line if in calibration mode
      if (calibrationMode && calibrationPoints.length > 0) {
        // Draw calibration points
        calibrationPoints.forEach((point, index) => {
          const cx = point.x * scaleX;
          const cy = point.y * scaleY;

          // Draw point
          ctx.beginPath();
          ctx.arc(cx, cy, 10, 0, 2 * Math.PI);
          ctx.fillStyle = "rgba(234, 88, 12, 0.9)"; // orange-600
          ctx.fill();
          ctx.strokeStyle = "rgba(255, 255, 255, 1)";
          ctx.lineWidth = 2;
          ctx.stroke();

          // Draw label (A or B)
          ctx.fillStyle = "#ffffff";
          ctx.font = "bold 12px sans-serif";
          ctx.textAlign = "center";
          ctx.textBaseline = "middle";
          ctx.fillText(index === 0 ? "A" : "B", cx, cy);
        });

        // Draw line between points if we have two
        if (calibrationPoints.length === 2) {
          const p1 = calibrationPoints[0];
          const p2 = calibrationPoints[1];

          ctx.beginPath();
          ctx.moveTo(p1.x * scaleX, p1.y * scaleY);
          ctx.lineTo(p2.x * scaleX, p2.y * scaleY);
          ctx.strokeStyle = "rgba(234, 88, 12, 0.9)"; // orange-600
          ctx.lineWidth = 2;
          ctx.setLineDash([5, 5]);
          ctx.stroke();
          ctx.setLineDash([]); // Reset dash

          // Draw pixel distance label at midpoint
          const midX = (p1.x * scaleX + p2.x * scaleX) / 2;
          const midY = (p1.y * scaleY + p2.y * scaleY) / 2;
          const pixelDist = Math.sqrt((p2.x - p1.x) ** 2 + (p2.y - p1.y) ** 2);

          ctx.fillStyle = "rgba(0, 0, 0, 0.7)";
          ctx.fillRect(midX - 40, midY - 10, 80, 20);
          ctx.fillStyle = "#ffffff";
          ctx.font = "bold 11px sans-serif";
          ctx.fillText(`${pixelDist.toFixed(0)} px`, midX, midY);
        }
      }
    };

    img.src = floorPlan.imageData;
  }, [floorPlan, samples, dimensions, heatmapMetric, calibrationMode, calibrationPoints]);

  // Handle canvas click
  const handleCanvasClick = (e: React.MouseEvent<HTMLCanvasElement>) => {
    if (!canvasRef.current || !floorPlan) return;

    const canvas = canvasRef.current;
    const rect = canvas.getBoundingClientRect();
    const x = e.clientX - rect.left;
    const y = e.clientY - rect.top;

    // Convert canvas coordinates to floor plan coordinates
    const scaleX = floorPlan.width / dimensions.width;
    const scaleY = floorPlan.height / dimensions.height;

    const floorPlanX = Math.round(x * scaleX);
    const floorPlanY = Math.round(y * scaleY);

    // Handle calibration click if in calibration mode
    if (calibrationMode && onCalibrationClick) {
      onCalibrationClick(floorPlanX, floorPlanY);
      return;
    }

    // Handle regular point click if interactive
    if (interactive && onPointClick) {
      onPointClick(floorPlanX, floorPlanY);
    }
  };

  return (
    <div ref={containerRef} className="w-full">
      <canvas
        ref={canvasRef}
        onClick={handleCanvasClick}
        className={`border border-surface-border ${radius.md} ${
          interactive || calibrationMode ? "cursor-crosshair" : ""
        }`}
        width={dimensions.width}
        height={dimensions.height}
      />
    </div>
  );
}

// Helper function to draw heatmap
function drawHeatmap(
  ctx: CanvasRenderingContext2D,
  samples: SamplePoint[],
  metric: "rssi" | "throughput" | "latency",
  scaleX: number,
  scaleY: number
) {
  if (samples.length === 0) return;

  // Extract metric values
  const values = samples.map((s) => {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any -- Polymorphic sample data with dynamic property access
    const data = s.sampleData as any;
    switch (metric) {
      case "rssi":
        return data.rssi || data.networks?.[0]?.rssi || -100;
      case "throughput":
        return data.downloadMbps || 0;
      case "latency":
        return data.latency || 0;
      default:
        return 0;
    }
  });

  const minValue = Math.min(...values);
  const maxValue = Math.max(...values);

  // Create gradient overlay
  const canvas = ctx.canvas;
  const imageData = ctx.createImageData(canvas.width, canvas.height);
  const data = imageData.data;

  // For each pixel, calculate interpolated value
  for (let y = 0; y < canvas.height; y++) {
    for (let x = 0; x < canvas.width; x++) {
      // Find nearest samples and interpolate
      let totalWeight = 0;
      let weightedValue = 0;

      // Use iterator pattern to avoid dynamic array indexing
      const valuesIterator = values[Symbol.iterator]();
      for (const sample of samples) {
        const sx = sample.x * scaleX;
        const sy = sample.y * scaleY;
        const distance = Math.sqrt((x - sx) ** 2 + (y - sy) ** 2);

        // Inverse distance weighting (IDW)
        const weight = distance === 0 ? 1000 : 1 / distance ** 2;
        totalWeight += weight;
        const valueResult = valuesIterator.next();
        const sampleValue = valueResult.done ? 0 : valueResult.value;
        weightedValue += weight * sampleValue;
      }

      const value = weightedValue / totalWeight;

      // Normalize to 0-1
      const normalized = (value - minValue) / (maxValue - minValue || 1);

      // Get color for this value
      const color = getHeatmapColor(normalized, metric);

      const pixelIndex = (y * canvas.width + x) * 4;
      // Bounds-checked pixel data assignment
      if (pixelIndex >= 0 && pixelIndex + 3 < data.length) {
        data.set([color.r, color.g, color.b, color.a], pixelIndex);
      }
    }
  }

  // Apply heatmap overlay
  ctx.globalAlpha = 0.5;
  ctx.putImageData(imageData, 0, 0);
  ctx.globalAlpha = 1.0;
}

// Get heatmap color based on normalized value (0-1)
function getHeatmapColor(
  value: number,
  metric: "rssi" | "throughput" | "latency"
): { r: number; g: number; b: number; a: number } {
  // For RSSI, lower is worse (invert)
  if (metric === "rssi") {
    value = 1 - value;
  }
  // For latency, lower is better (invert)
  if (metric === "latency") {
    value = 1 - value;
  }

  // Color gradient: red (0) -> yellow (0.5) -> green (1)
  let r, g, b;

  if (value < 0.5) {
    // Red to yellow
    r = 255;
    g = Math.round(255 * (value * 2));
    b = 0;
  } else {
    // Yellow to green
    r = Math.round(255 * (1 - (value - 0.5) * 2));
    g = 255;
    b = 0;
  }

  return { r, g, b, a: 200 }; // Semi-transparent
}
