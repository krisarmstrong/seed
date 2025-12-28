/**
 * LogViewerCard - Minimal dashboard card for system logs.
 *
 * Shows a summary view with:
 * - Live/Paused streaming status toggle
 * - Total logs, error count, warning count
 * - Maximize icon in header to open full-screen modal
 *
 * Full log viewing, filtering, and searching is done in the LogViewerModal.
 *
 * Usage:
 * ```tsx
 * <LogViewerCard className="my-custom-class" />
 * ```
 */

import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useLogs } from "../../hooks/useLogs";
import { cn, spacing, radius, layout, icon as iconTokens } from "../../styles/theme";
import { LogViewerModal } from "./LogViewerModal";
import { Maximize2 } from "../ui/Icons";

/** Props for the LogViewerCard component. */
export interface LogViewerCardProps {
  /** Additional CSS classes. */
  className?: string;
}

/**
 * LogViewerCard - Minimal dashboard card for system logs.
 * Shows summary stats (total, errors, warnings) and streaming status.
 * Full log viewing is done in the modal (click "View Logs").
 */
export function LogViewerCard({ className = "" }: LogViewerCardProps) {
  const { t } = useTranslation("common");
  const { stats, isStreaming, setIsStreaming, isLoading, error } = useLogs({
    maxLogs: 1000,
  });

  const [isModalOpen, setIsModalOpen] = useState(false);

  // Calculate error and warning counts
  const errorCount =
    stats?.by_level && "ERROR" in stats.by_level ? stats.by_level.ERROR : 0;
  const warnCount =
    stats?.by_level && "WARN" in stats.by_level ? stats.by_level.WARN : 0;

  return (
    <div
      className={cn(
        "bg-surface-raised",
        radius.lg,
        "border border-surface-border overflow-hidden",
        className
      )}
    >
      {/* Header */}
      <div
        className={cn(
          spacing.pad.md,
          "border-b border-surface-border",
          layout.flex.between,
          "items-center"
        )}
      >
        <div>
          <h2 className="heading-3 text-text-primary">
            {t("logs.title", "System Logs")}
          </h2>
        </div>
        <div className="flex items-center gap-2">
          {/* Full Screen button */}
          <button
            type="button"
            onClick={() => setIsModalOpen(true)}
            className={cn(
              "p-1.5",
              "bg-surface-hover text-text-secondary",
              radius.md,
              "hover:bg-surface-border hover:text-text-primary transition-colors flex items-center justify-center cursor-pointer"
            )}
            aria-label={t("logs.fullScreen", "Full Screen")}
            title={t("logs.fullScreen", "Full Screen")}
          >
            <Maximize2 className={iconTokens.size.sm} aria-hidden="true" />
          </button>
          {/* Streaming toggle */}
          <button
            type="button"
            onClick={() => setIsStreaming(!isStreaming)}
            className={cn(
              spacing.chip.sm,
              radius.md,
              "text-xs font-medium transition-colors",
              isStreaming
                ? "bg-status-success text-text-inverse hover:brightness-90"
                : "bg-surface-base text-text-primary hover:bg-surface-hover"
            )}
          >
            {isStreaming ? t("logs.streaming", "Live") : t("logs.paused", "Paused")}
          </button>
        </div>
      </div>

      {/* Minimal stats summary */}
      <div className={cn(spacing.pad.md, "space-y-3")}>
        {isLoading ? (
          <div className="text-text-secondary text-sm">
            {t("logs.loading", "Loading logs...")}
          </div>
        ) : error ? (
          <div className="text-status-error text-sm">{error}</div>
        ) : (
          <>
            {/* Stats row */}
            <div
              className={cn(
                layout.flex.between,
                "bg-surface-hover",
                radius.md,
                spacing.pad.sm
              )}
            >
              <div className={cn(layout.inline.comfortable, "text-sm")}>
                <span>
                  <strong className="text-text-primary">
                    {stats?.total_count ?? 0}
                  </strong>{" "}
                  <span className="text-text-secondary">
                    {t("logs.totalLogs", "logs")}
                  </span>
                </span>
                <span className="text-status-error">
                  <strong>{errorCount}</strong> {t("logs.errors", "errors")}
                </span>
                <span className="text-status-warning">
                  <strong>{warnCount}</strong> {t("logs.warnings", "warnings")}
                </span>
              </div>
              {stats?.errors_last_hour && stats.errors_last_hour > 0 && (
                <span className="text-status-error text-xs">
                  {stats.errors_last_hour} {t("logs.lastHour", "last hour")}
                </span>
              )}
            </div>

          </>
        )}
      </div>

      {/* Full Screen Modal */}
      <LogViewerModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
      />
    </div>
  );
}

export default LogViewerCard;
