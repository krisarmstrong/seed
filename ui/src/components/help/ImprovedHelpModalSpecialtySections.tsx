/**
 * Help modal specialty sections: profiles, WiFi survey, RTSP, DICOM, how-to, glossary, plus inline helper components.
 */

import type React from 'react';
import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { cn, layout, radius, spacing } from '../../styles/theme';
import { Search } from '../ui/icons';

interface TroubleshootingIssue {
  symptom: string;
  causes: string[];
  solutions: string[];
}

function _troubleshootingCategory({
  title,
  issues,
}: {
  title: string;
  issues: TroubleshootingIssue[];
}): React.JSX.Element {
  return (
    <div class={cn(spacing.margin.top.section)}>
      <h4 class={cn('font-semibold text-text-primary', spacing.margin.bottom.content)}>{title}</h4>
      <div class="stack-lg">
        {issues.map((issue) => (
          <div
            key={issue.symptom}
            class={cn('border border-surface-border', radius.default, spacing.pad.default)}
          >
            <h5 class={cn('font-semibold text-status-warning', spacing.margin.bottom.inline)}>
              {issue.symptom}
            </h5>
            <div class="grid md:grid-cols-2 gap-4 body-small">
              <div>
                <p class="font-semibold text-text-primary mb-1">Possible Causes:</p>
                <ul class={cn('text-text-secondary', spacing.margin.left.comfortable, 'list-disc')}>
                  {issue.causes.map((cause) => (
                    <li key={cause}>{cause}</li>
                  ))}
                </ul>
              </div>
              <div>
                <p class="font-semibold text-text-primary mb-1">Solutions:</p>
                <ul class={cn('text-text-secondary', spacing.margin.left.comfortable, 'list-disc')}>
                  {issue.solutions.map((solution) => (
                    <li key={solution}>{solution}</li>
                  ))}
                </ul>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

// ============================================================================
// HELPER COMPONENTS
// ============================================================================

function _featureCard({
  title,
  description,
}: {
  title: string;
  description: string;
}): React.JSX.Element {
  return (
    <div
      class={cn('bg-surface-hover border border-surface-border', radius.lg, spacing.pad.default)}
    >
      <h4 class={cn('font-semibold text-text-primary', spacing.margin.bottom.inline)}>{title}</h4>
      <p class="body-small text-text-secondary">{description}</p>
    </div>
  );
}

function _stepCard({
  number,
  title,
  description,
}: {
  number: number;
  title: string;
  description: string;
}): React.JSX.Element {
  return (
    <div class={cn('flex', spacing.gap.comfortable)}>
      <div
        class={cn(
          'shrink-0 w-8 h-8',
          radius.full,
          'bg-brand-primary text-text-inverse',
          layout.flex.center,
          'font-semibold',
        )}
      >
        {number}
      </div>
      <div class="flex-1">
        <h4 class={cn('font-semibold', spacing.margin.bottom.inline)}>{title}</h4>
        <p class="body-small">{description}</p>
      </div>
    </div>
  );
}

function _helpContentSection({
  title,
  children,
}: {
  title: string;
  children: ReactNode;
}): React.JSX.Element {
  return (
    <div class="max-w-3xl">
      <h3 class={cn('heading-2', spacing.margin.bottom.content)}>{title}</h3>
      {children}
    </div>
  );
}

function _helpTermList({
  items,
}: {
  items: Array<{ term: string; description: string }>;
}): React.JSX.Element {
  return (
    <dl class="stack-lg">
      {items.map((item) => (
        <div key={item.term} class={cn('border-l-2 border-surface-border', spacing.pad.default)}>
          <dt class={cn('font-semibold text-text-primary', spacing.margin.bottom.inline)}>
            {item.term}
          </dt>
          <dd class="body-small text-text-secondary">{item.description}</dd>
        </div>
      ))}
    </dl>
  );
}

// ============================================================================
// NEW FEATURE SECTIONS
// ============================================================================

function _profilesSection(): React.JSX.Element {
  const { t } = useTranslation('help');
  const capabilities = t('content.profiles.capabilities', { returnObjects: true }) as string[];
  const useCases = t('content.profiles.useCases.items', { returnObjects: true }) as Array<{
    name: string;
    description: string;
  }>;

  return (
    <helpContentSection title={t('sections.profiles')}>
      <p class={cn('body-small text-text-secondary', spacing.margin.bottom.content)}>
        {t('content.profiles.description')}
      </p>

      <div class={spacing.margin.bottom.section}>
        <h4 class={cn('font-semibold text-text-primary', spacing.margin.bottom.content)}>
          {t('content.profiles.overview.title')}
        </h4>
        <p class="body-small text-text-secondary">{t('content.profiles.overview.content')}</p>
      </div>

      <div class={spacing.margin.bottom.section}>
        <h4 class={cn('font-semibold text-text-primary', spacing.margin.bottom.content)}>
          {t('content.profiles.capabilities_title', 'Profile Capabilities')}
        </h4>
        <ul
          class={cn(
            'body-small text-text-secondary stack-sm',
            spacing.margin.left.spacious,
            'list-disc',
          )}
        >
          {capabilities?.map((cap) => (
            <li key={cap}>{cap}</li>
          ))}
        </ul>
      </div>

      <div>
        <h4 class={cn('font-semibold text-text-primary', spacing.margin.bottom.content)}>
          {t('content.profiles.useCases.title')}
        </h4>
        <div class="stack-lg">
          {useCases?.map((useCase) => (
            <div
              key={useCase.name}
              class={cn('border-l-2 border-brand-primary', spacing.pad.default)}
            >
              <dt class={cn('font-semibold text-text-primary', spacing.margin.bottom.inline)}>
                {useCase.name}
              </dt>
              <dd class="body-small text-text-secondary">{useCase.description}</dd>
            </div>
          ))}
        </div>
      </div>
    </helpContentSection>
  );
}

function _wiFiSurveySection(): React.JSX.Element {
  const { t } = useTranslation('help');
  const visualizations = t('content.wifiSurvey.visualizations', { returnObjects: true }) as Array<{
    type: string;
    description: string;
  }>;
  const bestPractices = t('content.wifiSurvey.bestPractices.items', {
    returnObjects: true,
  }) as string[];

  return (
    <helpContentSection title={t('sections.wifiSurvey')}>
      <p class={cn('body-small text-text-secondary', spacing.margin.bottom.content)}>
        {t('content.wifiSurvey.description')}
      </p>

      <helpTermList
        items={[
          {
            term: t('content.wifiSurvey.terms.floorPlan.term'),
            description: t('content.wifiSurvey.terms.floorPlan.description'),
          },
          {
            term: t('content.wifiSurvey.terms.heatmap.term'),
            description: t('content.wifiSurvey.terms.heatmap.description'),
          },
          {
            term: t('content.wifiSurvey.terms.surveyPoint.term'),
            description: t('content.wifiSurvey.terms.surveyPoint.description'),
          },
          {
            term: t('content.wifiSurvey.terms.dataRate.term'),
            description: t('content.wifiSurvey.terms.dataRate.description'),
          },
        ]}
      />

      <div class={spacing.margin.top.section}>
        <h4 class={cn('font-semibold text-text-primary', spacing.margin.bottom.content)}>
          {t('content.wifiSurvey.visualizationsTitle', 'Visualization Modes')}
        </h4>
        <div class="grid md:grid-cols-2 gap-4">
          {visualizations?.map((viz) => (
            <div
              key={viz.type}
              class={cn(
                'bg-surface-hover border border-surface-border',
                radius.default,
                spacing.pad.sm,
              )}
            >
              <h5 class="font-semibold text-text-primary">{viz.type}</h5>
              <p class="body-small text-text-secondary">{viz.description}</p>
            </div>
          ))}
        </div>
      </div>

      <div
        class={cn(
          spacing.margin.top.section,
          'bg-status-info/10 border border-status-info/20',
          radius.default,
          spacing.pad.default,
        )}
      >
        <h4 class={cn('font-semibold text-text-primary', spacing.margin.bottom.inline)}>
          {t('content.wifiSurvey.bestPractices.title')}
        </h4>
        <ul
          class={cn(
            'body-small text-text-secondary stack-sm',
            spacing.margin.left.spacious,
            'list-disc',
          )}
        >
          {bestPractices?.map((practice) => (
            <li key={practice}>{practice}</li>
          ))}
        </ul>
      </div>
    </helpContentSection>
  );
}

function _rtspChecksSection(): React.JSX.Element {
  const { t } = useTranslation('help');
  const configuration = t('content.rtspChecks.configuration', { returnObjects: true }) as Array<{
    field: string;
    description: string;
  }>;

  return (
    <helpContentSection title={t('sections.rtspChecks')}>
      <p class={cn('body-small text-text-secondary', spacing.margin.bottom.content)}>
        {t('content.rtspChecks.description')}
      </p>

      <helpTermList
        items={[
          {
            term: t('content.rtspChecks.terms.rtsp.term'),
            description: t('content.rtspChecks.terms.rtsp.description'),
          },
          {
            term: t('content.rtspChecks.terms.options.term'),
            description: t('content.rtspChecks.terms.options.description'),
          },
          {
            term: t('content.rtspChecks.terms.describe.term'),
            description: t('content.rtspChecks.terms.describe.description'),
          },
          {
            term: t('content.rtspChecks.terms.authentication.term'),
            description: t('content.rtspChecks.terms.authentication.description'),
          },
        ]}
      />

      <div class={spacing.margin.top.section}>
        <h4 class={cn('font-semibold text-text-primary', spacing.margin.bottom.content)}>
          {t('content.rtspChecks.configurationTitle', 'Configuration Options')}
        </h4>
        <div class="stack-sm">
          {configuration?.map((config) => (
            <div key={config.field} class={cn('border-l-2 border-surface-border', spacing.pad.sm)}>
              <span class="font-mono text-brand-primary">{config.field}</span>
              <span class="body-small text-text-secondary ml-2">{config.description}</span>
            </div>
          ))}
        </div>
      </div>
    </helpContentSection>
  );
}

function _dicomChecksSection(): React.JSX.Element {
  const { t } = useTranslation('help');
  const configuration = t('content.dicomChecks.configuration', { returnObjects: true }) as Array<{
    field: string;
    description: string;
  }>;
  const commonIssues = t('content.dicomChecks.commonIssues', { returnObjects: true }) as Array<{
    issue: string;
    solution: string;
  }>;

  return (
    <helpContentSection title={t('sections.dicomChecks')}>
      <p class={cn('body-small text-text-secondary', spacing.margin.bottom.content)}>
        {t('content.dicomChecks.description')}
      </p>

      <helpTermList
        items={[
          {
            term: t('content.dicomChecks.terms.dicom.term'),
            description: t('content.dicomChecks.terms.dicom.description'),
          },
          {
            term: t('content.dicomChecks.terms.cEcho.term'),
            description: t('content.dicomChecks.terms.cEcho.description'),
          },
          {
            term: t('content.dicomChecks.terms.aeTitle.term'),
            description: t('content.dicomChecks.terms.aeTitle.description'),
          },
          {
            term: t('content.dicomChecks.terms.scp.term'),
            description: t('content.dicomChecks.terms.scp.description'),
          },
          {
            term: t('content.dicomChecks.terms.scu.term'),
            description: t('content.dicomChecks.terms.scu.description'),
          },
        ]}
      />

      <div class={spacing.margin.top.section}>
        <h4 class={cn('font-semibold text-text-primary', spacing.margin.bottom.content)}>
          {t('content.dicomChecks.configurationTitle', 'Configuration')}
        </h4>
        <div class="stack-sm">
          {configuration?.map((config) => (
            <div key={config.field} class={cn('border-l-2 border-surface-border', spacing.pad.sm)}>
              <span class="font-mono text-brand-primary">{config.field}</span>
              <span class="body-small text-text-secondary ml-2">{config.description}</span>
            </div>
          ))}
        </div>
      </div>

      <div
        class={cn(
          spacing.margin.top.section,
          'bg-status-warning/10 border border-status-warning/20',
          radius.default,
          spacing.pad.default,
        )}
      >
        <h4 class={cn('font-semibold text-text-primary', spacing.margin.bottom.content)}>
          {t('content.dicomChecks.commonIssuesTitle', 'Common Issues')}
        </h4>
        <div class="stack-lg">
          {commonIssues?.map((item) => (
            <div key={item.issue}>
              <p class="font-semibold text-status-warning">{item.issue}</p>
              <p class="body-small text-text-secondary">{item.solution}</p>
            </div>
          ))}
        </div>
      </div>
    </helpContentSection>
  );
}

function _howToSection(): React.JSX.Element {
  const { t } = useTranslation('help');
  const guides = t('content.howTo.guides', { returnObjects: true }) as Record<
    string,
    {
      title: string;
      description: string;
      steps: string[];
    }
  >;

  return (
    <helpContentSection title={t('sections.howTo')}>
      <p class={cn('body-small text-text-secondary', spacing.margin.bottom.content)}>
        {t('content.howTo.description')}
      </p>

      <div class="stack-xl">
        {guides
          ? Object.entries(guides).map(([key, guide]) => (
              <div
                key={key}
                class={cn('border border-surface-border', radius.lg, spacing.pad.default)}
              >
                <h4 class={cn('font-semibold text-text-primary', spacing.margin.bottom.inline)}>
                  {guide.title}
                </h4>
                <p class={cn('body-small text-text-secondary', spacing.margin.bottom.content)}>
                  {guide.description}
                </p>
                <ol
                  class={cn(
                    'body-small text-text-secondary stack-sm',
                    spacing.margin.left.spacious,
                    'list-decimal',
                  )}
                >
                  {guide.steps.map((step) => (
                    <li key={`${key}-${step.slice(0, 50)}`}>{step}</li>
                  ))}
                </ol>
              </div>
            ))
          : null}
      </div>
    </helpContentSection>
  );
}

function _glossarySection(): React.JSX.Element {
  const { t } = useTranslation('glossary');
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedCategory, setSelectedCategory] = useState<string>('all');

  const categories = t('categories', { returnObjects: true }) as Record<string, string>;
  const terms = t('terms', { returnObjects: true }) as Record<
    string,
    {
      term: string;
      fullName: string;
      definition: string;
      category: string;
    }
  >;

  const filteredTerms = terms
    ? Object.entries(terms).filter(([, termData]) => {
        const matchesSearch =
          searchTerm === '' ||
          termData.term.toLowerCase().includes(searchTerm.toLowerCase()) ||
          termData.fullName.toLowerCase().includes(searchTerm.toLowerCase()) ||
          termData.definition.toLowerCase().includes(searchTerm.toLowerCase());

        const matchesCategory =
          selectedCategory === 'all' || termData.category === selectedCategory;

        return matchesSearch && matchesCategory;
      })
    : [];

  return (
    <div class="max-w-3xl">
      <h3 class={cn('heading-2', spacing.margin.bottom.content)}>{t('title')}</h3>
      <p class={cn('body-small text-text-secondary', spacing.margin.bottom.content)}>
        {t('description')}
      </p>

      {/* Search and Filter */}
      <div class={cn('flex flex-wrap gap-4', spacing.margin.bottom.section)}>
        <div class="flex-1 min-w-[200px]">
          <div class="relative">
            <Search
              class={cn('absolute left-3 top-1/2 -translate-y-1/2', 'w-4 h-4 text-text-muted')}
            />
            <input
              type="text"
              placeholder="Search terms..."
              value={searchTerm}
              onChange={(e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>): void =>
                setSearchTerm(e.target.value)
              }
              class={cn(
                'w-full pl-9 pr-3 py-2',
                'body-small',
                radius.default,
                'border border-surface-border bg-surface-raised text-text-primary placeholder-text-muted',
                'focus:outline-none focus:ring-2 focus:ring-brand-primary',
              )}
            />
          </div>
        </div>
        <select
          value={selectedCategory}
          onChange={(e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>): void =>
            setSelectedCategory(e.target.value)
          }
          class={cn(
            'px-3 py-2',
            'body-small',
            radius.default,
            'border border-surface-border bg-surface-raised text-text-primary',
            'focus:outline-none focus:ring-2 focus:ring-brand-primary',
          )}
        >
          <option value="all">All Categories</option>
          {categories
            ? Object.entries(categories).map(([key, label]) => (
                <option key={key} value={key}>
                  {label}
                </option>
              ))
            : null}
        </select>
      </div>

      {/* Terms List */}
      <div class="stack-lg">
        {filteredTerms.map(([key, termData]) => (
          <div
            key={key}
            class={cn(
              'border border-surface-border',
              radius.default,
              spacing.pad.default,
              'hover:border-brand-primary/50 transition-colors',
            )}
          >
            <div class="flex items-start justify-between gap-4">
              <div class="flex-1">
                <div class="flex items-baseline gap-2 mb-1">
                  <span class="font-bold text-brand-primary">{termData.term}</span>
                  <span class="body-small text-text-muted">({termData.fullName})</span>
                </div>
                <p class="body-small text-text-secondary">{termData.definition}</p>
              </div>
              <span
                class={cn(
                  'px-2 py-0.5 text-xs font-medium',
                  radius.default,
                  'bg-surface-hover text-text-muted capitalize',
                )}
              >
                {categories?.[termData.category] || termData.category}
              </span>
            </div>
          </div>
        ))}

        {filteredTerms.length === 0 && (
          <div class="text-center py-8 text-text-muted">No terms found matching your search.</div>
        )}
      </div>
    </div>
  );
}
