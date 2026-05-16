/**
 * ImprovedHelpModal Component (~681 lines)
 *
 * Purpose: Comprehensive application help modal providing user guidance across multiple topics.
 * Features tabbed navigation, search functionality, and rich content for all major features.
 *
 * Key Features:
 * - Multi-section help: About, Network Discovery, WiFi, Cable/Link, Performance, etc.
 * - Search functionality: Filter help content by keyword
 * - Icon-based navigation: Visual section selector with icons
 * - Rich content: Markdown-like formatting for help text
 * - Modal overlay: Centered help dialog with close button
 * - Responsive design: Adapts to different screen sizes
 * - Keyboard support: ESC key closes modal
 * - Scrollable sections: Long help content in scrollable containers
 *
 * Usage:
 * ```typescript
 * <ImprovedHelpModal isOpen={showHelp} onClose={() => setShowHelp(false)} />
 * ```
 *
 * Dependencies: Icons, theme utilities, useState for tab/search state management
 * State: activeSection (current tab), searchQuery (help search text)
 */

import type React from 'react';
import { type ReactNode, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { cn, icon as iconTokens, layout, modal, radius, spacing } from '../../styles/theme';
import {
  Activity,
  AlertTriangle,
  BookOpen,
  Cable,
  Heart,
  HeartPulse,
  Info,
  LayoutDashboard,
  Lightbulb,
  Monitor,
  Network,
  Search,
  Server,
  Shield,
  Signal,
  SlidersHorizontal,
  Wifi,
  Zap,
} from '../ui/icons';

interface HelpModalProps {
  isOpen: boolean;
  onClose: () => void;
  /** Application version from backend */
  version?: string;
}

interface HelpSection {
  id: string;
  title: string;
  icon: ReactNode;
  content: ReactNode;
}

/**
 * ImprovedHelpModal Component
 * Renders a modal dialog with tabbed help content and search functionality
 */
export function ImprovedHelpModal({
  isOpen,
  onClose,
  version = 'dev',
}: HelpModalProps): React.JSX.Element | null {
  const { t } = useTranslation('help');
  // Track which help section is currently active
  const [activeSection, setActiveSection] = useState<string>('about');
  // Track search query for filtering help content
  const [searchQuery, setSearchQuery] = useState('');

  if (!isOpen) {
    return null;
  }

  const sections: HelpSection[] = [
    {
      id: 'about',
      title: t('sections.about'),
      icon: <Info class={iconTokens.size.sm} />,
      content: <aboutSection version={version} />,
    },
    {
      id: 'getting-started',
      title: t('sections.gettingStarted'),
      icon: <LayoutDashboard class={iconTokens.size.sm} />,
      content: <gettingStartedSection />,
    },
    {
      id: 'link',
      title: t('sections.link'),
      icon: <Activity class={iconTokens.size.sm} />,
      content: <linkStatusSection />,
    },
    {
      id: 'cable',
      title: t('sections.cable'),
      icon: <Cable class={iconTokens.size.sm} />,
      content: <cableTestSection />,
    },
    {
      id: 'wifi',
      title: t('sections.wifi'),
      icon: <Wifi class={iconTokens.size.sm} />,
      content: <wiFiStatusSection />,
    },
    {
      id: 'network',
      title: t('sections.network'),
      icon: <Network class={iconTokens.size.sm} />,
      content: <networkSection />,
    },
    {
      id: 'gateway',
      title: t('sections.gateway'),
      icon: <Server class={iconTokens.size.sm} />,
      content: <gatewaySection />,
    },
    {
      id: 'dns',
      title: t('sections.dns'),
      icon: <Search class={iconTokens.size.sm} />,
      content: <dnsSection />,
    },
    {
      id: 'performance',
      title: t('sections.performance'),
      icon: <Zap class={iconTokens.size.sm} />,
      content: <performanceSection />,
    },
    {
      id: 'discovery',
      title: t('sections.discovery'),
      icon: <Search class={iconTokens.size.sm} />,
      content: <discoverySection />,
    },
    {
      id: 'healthChecks',
      title: t('sections.healthChecks'),
      icon: <Heart class={iconTokens.size.sm} />,
      content: <healthChecksSection />,
    },
    {
      id: 'security',
      title: t('sections.security'),
      icon: <Shield class={iconTokens.size.sm} />,
      content: <securitySection />,
    },
    {
      id: 'troubleshooting',
      title: t('sections.troubleshooting'),
      icon: <AlertTriangle class={iconTokens.size.sm} />,
      content: <troubleshootingSection />,
    },
    {
      id: 'profiles',
      title: t('sections.profiles'),
      icon: <SlidersHorizontal class={iconTokens.size.sm} />,
      content: <profilesSection />,
    },
    {
      id: 'wifiSurvey',
      title: t('sections.wifiSurvey'),
      icon: <Signal class={iconTokens.size.sm} />,
      content: <wiFiSurveySection />,
    },
    {
      id: 'rtspChecks',
      title: t('sections.rtspChecks'),
      icon: <Monitor class={iconTokens.size.sm} />,
      content: <rtspChecksSection />,
    },
    {
      id: 'dicomChecks',
      title: t('sections.dicomChecks'),
      icon: <HeartPulse class={iconTokens.size.sm} />,
      content: <dicomChecksSection />,
    },
    {
      id: 'howTo',
      title: t('sections.howTo'),
      icon: <Lightbulb class={iconTokens.size.sm} />,
      content: <howToSection />,
    },
    {
      id: 'glossary',
      title: t('sections.glossary'),
      icon: <BookOpen class={iconTokens.size.sm} />,
      content: <glossarySection />,
    },
  ];

  const filteredSections = sections.filter(
    (section) =>
      section.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
      section.id.toLowerCase().includes(searchQuery.toLowerCase()),
  );

  const currentSection = sections.find((s) => s.id === activeSection);

  return (
    <div class={modal.overlay}>
      {/* Backdrop */}
      <div class={modal.backdrop} onClick={onClose} aria-hidden="true" />

      {/* Modal */}
      <div
        class={cn(
          'relative',
          modal.content,
          modal.size.xl,
          radius.lg,
          'flex flex-col overflow-hidden',
        )}
        role="dialog"
        aria-modal="true"
        aria-labelledby="help-modal-title"
      >
        {/* Header */}
        <div
          class={cn(
            layout.flex.between,
            spacing.pad.default,
            'border-b border-surface-border shrink-0',
          )}
        >
          <h2 id="help-modal-title" class="heading-3">
            {t('modal.title')}
          </h2>
          <button
            type="button"
            onClick={onClose}
            class={cn(
              spacing.pad.xs,
              'text-text-muted hover:text-text-primary transition-colors',
              radius.default,
              'hover:bg-surface-hover',
            )}
            aria-label={t('modal.closeHelp')}
          >
            <svg
              class={iconTokens.size.md}
              viewBox="0 0 20 20"
              fill="currentColor"
              aria-hidden="true"
            >
              <path
                fillRule="evenodd"
                d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
                clipRule="evenodd"
              />
            </svg>
          </button>
        </div>

        {/* Content area with sidebar */}
        <div class="flex flex-1 overflow-hidden">
          {/* Sidebar / TOC */}
          <aside class="w-64 border-r border-surface-border bg-surface-base overflow-y-auto shrink-0">
            {/* Search */}
            <div class={cn(spacing.pad.sm, 'border-b border-surface-border')}>
              <div class="relative">
                <Search
                  class={cn(
                    'absolute left-3 top-1/2 -translate-y-1/2',
                    iconTokens.size.sm,
                    'text-text-muted',
                  )}
                />
                <input
                  type="text"
                  placeholder={t('modal.searchPlaceholder')}
                  value={searchQuery}
                  onChange={(e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>): void =>
                    setSearchQuery(e.target.value)
                  }
                  class={cn(
                    'w-full pl-9',
                    spacing.chip.lg,
                    'body-small',
                    radius.default,
                    'border border-surface-border bg-surface-raised text-text-primary placeholder-text-muted focus:outline-none focus:ring-2 focus:ring-brand-primary',
                  )}
                />
              </div>
            </div>

            {/* Table of Contents */}
            <nav class={cn(spacing.pad.xs, 'stack-xs')}>
              <p class={cn('caption', spacing.chip.lg, 'uppercase tracking-wider')}>
                {t('modal.contents')}
              </p>
              {filteredSections.map((section) => (
                <button
                  type="button"
                  key={section.id}
                  onClick={(): void => setActiveSection(section.id)}
                  class={cn(
                    'w-full flex items-center',
                    spacing.gap.default,
                    spacing.tab,
                    radius.default,
                    'body-small transition-colors text-left',
                    activeSection === section.id
                      ? 'bg-brand-primary/10 text-brand-primary font-medium'
                      : 'text-text-secondary hover:bg-surface-hover hover:text-text-primary',
                  )}
                >
                  {section.icon}
                  <span>{section.title}</span>
                </button>
              ))}
            </nav>
          </aside>

          {/* Main content */}
          <main class={cn('flex-1 overflow-y-auto', spacing.pad.lg)}>
            {currentSection && <div>{currentSection.content}</div>}
          </main>
        </div>
      </div>
    </div>
  );
}
