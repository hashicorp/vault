/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export const DASHBOARD = {
  cardName: (name) => `[data-test-card="${name}"]`,
  emptyState: (name) => `[data-test-empty-state="${name}"]`,
  emptyStateTitle: (name) => `[data-test-empty-state="${name}"] [data-test-empty-state-title]`,
  emptyStateMessage: (name) => `[data-test-empty-state="${name}"] [data-test-empty-state-message]`,
  emptyStateActions: (name) => `[data-test-empty-state="${name}"] [data-test-empty-state-actions]`,
  cardHeader: (name) => `[data-test-dashboard-card-header="${name}"]`,
  tableRow: (name) => `[data-test-dashboard-table="${name}"] tr`,
  searchSelect: (name) => `[data-test-search-select="${name}"]`,
  kvSearchSelect: `[data-test-kv-suggestion-input]`,
  actionButton: (action) => `[data-test-button="${action}"]`,
  title: (name) => `[data-test-title="${name}"]`,
  subtitle: (name) => `[data-test-card-subtitle="${name}"]`,
  subtext: (name) => `[data-test-subtext="${name}"]`,
  tooltipTitle: (name) => `[data-test-tooltip-title="${name}"]`,
  tooltipIcon: (type, name, icon) =>
    `[data-test-type="${type}"] [data-test-tooltip-title="${name}"] [data-test-icon="${icon}"]`,
  statLabel: (name) => `[data-test-stat-text="${name}"] .stat-label`,
  statText: (name) => `[data-test-stat-text="${name}"] .stat-text`,
  statValue: (name) => `[data-test-stat-text="${name}"] .stat-value`,
  selectEl: 'select',
  secretsEnginesCard: {
    secretEngineAccessorRow: (engineId) => `[data-test-secrets-engines-row=${engineId}] [data-test-accessor]`,
    secretEngineDescription: (engineId) =>
      `[data-test-secrets-engines-row=${engineId}] [data-test-description]`,
  },
  vaultConfigurationCard: {
    configDetailsField: (name) => `[data-test-vault-config-details="${name}"]`,
  },
};
