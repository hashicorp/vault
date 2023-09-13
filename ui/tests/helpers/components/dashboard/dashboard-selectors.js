export const SELECTORS = {
  cardName: (name) => `[data-test-card="${name}"]`,
  emptyState: (name) => `[data-test-empty-state="${name}"]`,
  emptyStateTitle: (name) => `[data-test-empty-state="${name}"] [data-test-empty-state-title]`,
  emptyStateMessage: (name) => `[data-test-empty-state="${name}"] [data-test-empty-state-message]`,
  emptyStateActions: (name) => `[data-test-empty-state="${name}"] [data-test-empty-state-actions]`,
  cardHeader: (name) => `[data-test-dashboard-card-header="${name}"]`,
  tableRow: (name) => `[data-test-dashboard-table="${name}"] tr`,
  replicationCard: {
    getReplicationTitle: (type, name) => `[data-test-${type}-replication] [data-test-title="${name}"]`,
    getStateTooltipTitle: (type, name) =>
      `[data-test-${type}-replication] [data-test-tooltip-title="${name}"]`,
    getStateTooltipIcon: (type, name, icon) =>
      `[data-test-${type}-replication] [data-test-tooltip-title="${name}"] [data-test-icon="${icon}"]`,
    drOnlyStateSubText: '[data-test-dr-replication] [data-test-subtext="state"]',
    knownSecondariesLabel: '[data-test-stat-text="known secondaries"] .stat-label',
    knownSecondariesSubtext: '[data-test-stat-text="known secondaries"] .stat-text',
    knownSecondariesValue: '[data-test-stat-text="known secondaries"] .stat-value',
  },
};
