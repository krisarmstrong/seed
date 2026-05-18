// biome-ignore-all lint/style/noInferrableTypes: useExplicitType requires types on default params
/**
 * Number and Time Formatting Utilities
 *
 * Provides safe formatting functions that guard against NaN, Infinity,
 * and other invalid numeric values to prevent "NaN" from displaying in the UI.
 *
 * #775: Centralized formatting with NaN guards to prevent display issues.
 */

/**
 * Checks if a value is a valid, finite number (not NaN, Infinity, or non-numeric).
 */
export function isValidNumber(value: unknown): value is number {
  return typeof value === 'number' && Number.isFinite(value);
}

/**
 * Safely formats a time value in milliseconds.
 * Returns a dash "-" for invalid/zero values.
 *
 * @param ms - Time in milliseconds
 * @param fallback - Value to return if ms is invalid (default: "-")
 */
export function formatTime(ms: number | undefined | null, fallback: string = '-'): string {
  if (!isValidNumber(ms) || ms <= 0) {
    return fallback;
  }
  if (ms < 1) {
    return '<1ms';
  }
  if (ms >= 1000) {
    return `${(ms / 1000).toFixed(1)}s`;
  }
  return `${Math.round(ms * 10) / 10}ms`;
}

/**
 * Safely formats a latency value (alias for formatTime).
 */
export const formatLatency: typeof formatTime = formatTime;

/**
 * Safely formats a byte size value to human-readable format.
 * Returns fallback for invalid values.
 *
 * @param bytes - Size in bytes
 * @param decimals - Number of decimal places (default: 1)
 * @param fallback - Value to return if invalid (default: "-")
 */
export function formatBytes(
  bytes: number | undefined | null,
  decimals: number = 1,
  fallback: string = '-',
): string {
  if (!isValidNumber(bytes) || bytes < 0) {
    return fallback;
  }
  if (bytes === 0) {
    return '0 B';
  }

  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  const size = i < sizes.length ? sizes[i] : sizes.at(-1);

  return `${Number.parseFloat((bytes / k ** i).toFixed(decimals))} ${size}`;
}
