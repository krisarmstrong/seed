import { twMerge } from 'tailwind-merge';

import { badge, button, card, input, modal } from './themeComponents';

/**
 * =============================================================================
 * THE SEED DESIGN SYSTEM - Mustard Seed Networks
 * =============================================================================
 *
 * Centralized design tokens and utilities for consistent UI across the app.
 *
 * ARCHITECTURE:
 * 1. CSS Variables (index.css) - Core color tokens for light/dark modes
 * 2. This file (theme.ts) - Barrel that re-exports the per-domain token
 *    modules + the class-name composition helpers (cn, buttonClass, etc.)
 * 3. Tailwind Classes - CSS-first configuration using @theme directive
 *
 * Token modules:
 *  - theme_spacing.ts    — margin / padding / gap / stack tokens
 *  - theme_typography.ts — heading / body / size / weight / family / leading
 *  - theme_components.ts — button, input, card, badge, toast, alert, modal, section
 *  - theme_colors.ts     — status, severity, timing, category, moduleColor, brand, gauge
 *  - theme_layout.ts     — sizing, icon, radius, border, layout
 *
 * BRAND COLORS:
 * - Primary: Seed Green (#2d7a3e / #81c784 dark) - Actions, links, focus states
 * - Accent: Lighter Seed Green (#4caf50 / #a5d6a7 dark) - Hover states
 * - Gold: Mustard Gold (#d4a017 / #fbbf24 dark) - Special highlights, premium
 *
 * USAGE:
 * import { spacing, button, cn, moduleColor } from '../styles/theme';
 * <button className={cn(button.base, button.variant.primary)}>Action</button>
 *
 * =============================================================================
 */

// biome-ignore lint/performance/noBarrelFile: theme.ts is the public design-token surface used by ~90 components; per-domain re-exports keep existing call sites working unchanged.
export {
  brand,
  category,
  discoveryMethod,
  gauge,
  moduleColor,
  progressBar,
  severity,
  status,
  timing,
} from './themeColors';
export { alert, badge, button, card, input, modal, section, toast } from './themeComponents';
export { border, icon, layout, radius, sizing } from './themeLayout';
// Re-export domain token modules so existing call sites keep working.
export { spacing } from './themeSpacing';
export { typography } from './themeTypography';

// ============================================================================
// UTILITY FUNCTIONS
// ============================================================================

/**
 * Combine class names with Tailwind class conflict resolution.
 * Uses tailwind-merge to properly handle conflicting Tailwind classes
 * (e.g., z-50 vs z-20, p-4 vs p-2 will resolve to the last value).
 */
export function cn(...classes: (string | boolean | undefined | null)[]): string {
  return twMerge(classes.filter(Boolean).join(' '));
}

// Type-safe Maps for dynamic lookups
const buttonVariantMap: Map<keyof typeof button.variant, string> = new Map<
  keyof typeof button.variant,
  string
>(Object.entries(button.variant) as [keyof typeof button.variant, string][]);
const buttonSizeMap: Map<keyof typeof button.size, string> = new Map<
  keyof typeof button.size,
  string
>(Object.entries(button.size) as [keyof typeof button.size, string][]);

/**
 * Build a button class string
 */
export function buttonClass(
  variant: keyof typeof button.variant = 'primary',
  size: keyof typeof button.size = 'md',
  className?: string,
): string {
  return cn(button.base, buttonVariantMap.get(variant), buttonSizeMap.get(size), className);
}

// Type-safe Maps for input lookups
const inputStateMap: Map<keyof typeof input.state, string> = new Map<
  keyof typeof input.state,
  string
>(Object.entries(input.state) as [keyof typeof input.state, string][]);
const inputSizeMap: Map<keyof typeof input.size, string> = new Map<keyof typeof input.size, string>(
  Object.entries(input.size) as [keyof typeof input.size, string][],
);

/**
 * Build an input class string
 */
export function inputClass(
  state: keyof typeof input.state = 'default',
  size: keyof typeof input.size = 'md',
  className?: string,
): string {
  return cn(input.base, inputStateMap.get(state), inputSizeMap.get(size), className);
}

// Type-safe Maps for card lookups
const cardVariantMap: Map<keyof typeof card.variant, string> = new Map<
  keyof typeof card.variant,
  string
>(Object.entries(card.variant) as [keyof typeof card.variant, string][]);
const cardPaddingMap: Map<keyof typeof card.padding, string> = new Map<
  keyof typeof card.padding,
  string
>(Object.entries(card.padding) as [keyof typeof card.padding, string][]);

/**
 * Build a card class string
 */
export function cardClass(
  variant: keyof typeof card.variant = 'default',
  padding: keyof typeof card.padding = 'md',
  className?: string,
): string {
  return cn(card.base, cardVariantMap.get(variant), cardPaddingMap.get(padding), className);
}

// Type-safe Map for badge lookups
const badgeVariantMap: Map<keyof typeof badge.variant, string> = new Map<
  keyof typeof badge.variant,
  string
>(Object.entries(badge.variant) as [keyof typeof badge.variant, string][]);

/**
 * Build a badge class string
 */
export function badgeClass(
  variant: keyof typeof badge.variant = 'default',
  className?: string,
): string {
  return cn(badge.base, badgeVariantMap.get(variant), className);
}

// Type-safe Maps for modal lookups
const modalSizeMap: Map<keyof typeof modal.size, string> = new Map<keyof typeof modal.size, string>(
  Object.entries(modal.size) as [keyof typeof modal.size, string][],
);
const modalPaddingMap: Map<keyof typeof modal.padding, string> = new Map<
  keyof typeof modal.padding,
  string
>(Object.entries(modal.padding) as [keyof typeof modal.padding, string][]);

/**
 * Build a modal class string
 */
export function modalClass(
  size: keyof typeof modal.size = 'md',
  padding: keyof typeof modal.padding = 'md',
  className?: string,
): string {
  return cn(modal.content, modalSizeMap.get(size), modalPaddingMap.get(padding), className);
}
