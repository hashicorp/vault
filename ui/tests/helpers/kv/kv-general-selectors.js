/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

export const SELECTORS = {
  toggleJson: '[data-test-json-view-toggle]',
  toggleMasked: '[data-test-button="toggle-masked"]',
  jsonEditor: '[data-test-component="code-mirror-modifier"]',
  tooltipTrigger: '[data-test-tooltip-trigger]',
  pageTitle: '[data-test-header-title]',
  infoRowValue: (label) => `[data-test-value-div="${label}"]`,
  secretTab: (tab) => `[data-test-secrets-tab="${tab}"]`,
  emptyStateTitle: '[data-test-empty-state-title]',
  emptyStateMessage: '[data-test-empty-state-message]',
  kvRow: '[data-test-kv-row]',
  maxVersions: '[data-test-input="maxVersions"]',
  casRequired: '[data-test-input="casRequired"]',
  automateSecretDeletion: '[data-test-ttl-value="Automate secret deletion"]',
  inlineAlert: '[data-test-inline-alert]',
  metadataCancel: '[data-test-kv-metadata-cancel]',
};

export const parseJsonEditor = (find) => {
  return JSON.parse(find(SELECTORS.jsonEditor).innerText);
};
