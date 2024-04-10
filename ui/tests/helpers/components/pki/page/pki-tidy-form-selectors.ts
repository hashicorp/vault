/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export const PKI_TIDY_FORM = {
  tidyFormName: (attr: string) => `[data-test-tidy-form="${attr}"]`,
  inputByAttr: (attr: string) => `[data-test-input="${attr}"]`,
  toggleInput: (attr: string) => `[data-test-input="${attr}"] input`,
  intervalDuration: '[data-test-ttl-value="Automatic tidy enabled"]',
  acmeAccountSafetyBuffer: '[data-test-ttl-value="Tidy ACME enabled"]',
  toggleLabel: (label: string) => `[data-test-toggle-label="${label}"]`,
  tidySectionHeader: (header: string) => `[data-test-tidy-header="${header}"]`,
  tidySave: '[data-test-pki-tidy-button]',
  tidyCancel: '[data-test-pki-tidy-cancel]',
  tidyPauseDuration: '[data-test-ttl-value="Pause duration"]',
  editAutoTidyButton: '[data-test-pki-edit-tidy-auto-link]',
};
