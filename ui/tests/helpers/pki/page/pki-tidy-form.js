/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export const SELECTORS = {
  tidyFormName: (attr) => `[data-test-tidy-form="${attr}"]`,
  inputByAttr: (attr) => `[data-test-input="${attr}"]`,
  toggleInput: (attr) => `[data-test-input="${attr}"] input`,
  intervalDuration: '[data-test-ttl-value="Automatic tidy enabled"]',
  acmeAccountSafetyBuffer: '[data-test-ttl-value="Tidy ACME enabled"]',
  toggleLabel: (label) => `[data-test-toggle-label="${label}"]`,
  tidySectionHeader: (header) => `[data-test-tidy-header="${header}"]`,
  tidySave: '[data-test-pki-tidy-button]',
  tidyCancel: '[data-test-pki-tidy-cancel]',
  tidyPauseDuration: '[data-test-ttl-value="Pause duration"]',
  editAutoTidyButton: '[data-test-pki-edit-tidy-auto-link]',
};
