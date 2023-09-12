/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

const SELECTORS = {
  getReplicationTitle: (type, name) => `[data-test-${type}-replication] [data-test-title="${name}"]`,
  getStateTooltipTitle: (type, name) => `[data-test-${type}-replication] [data-test-tooltip-title="${name}"]`,
  getStateTooltipIcon: (type, name, icon) =>
    `[data-test-${type}-replication] [data-test-tooltip-title="${name}"] [data-test-icon="${icon}"]`,
  drOnlyStateSubText: '[data-test-dr-replication] [data-test-subtext="state"]',
  knownSecondariesLabel: '[data-test-stat-text="known secondaries"] .stat-label',
  knownSecondariesSubtext: '[data-test-stat-text="known secondaries"] .stat-text',
  knownSecondariesValue: '[data-test-stat-text="known secondaries"] .stat-value',
  replicationEmptyState: '[data-test-card="replication"] [data-test-component="empty-state"]',
  replicationEmptyStateTitle:
    '[data-test-card="replication"] [data-test-component="empty-state"] .empty-state-title',
  replicationEmptyStateMessage:
    '[data-test-card="replication"] [data-test-component="empty-state"] .empty-state-message',
  replicationEmptyStateActions:
    '[data-test-card="replication"] [data-test-component="empty-state"] .empty-state-actions',
  replicationEmptyStateActionsLink:
    '[data-test-card="replication"] [data-test-component="empty-state"] .empty-state-actions a',
};

export default SELECTORS;
