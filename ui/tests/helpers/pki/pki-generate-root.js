/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

export const SELECTORS = {
  mainSectionTitle: '[data-test-generate-root-title="Root parameters"]',
  urlSectionTitle: '[data-test-generate-root-title="Issuer URLs"]',
  keyParamsGroupToggle: '[data-test-toggle-group="Key parameters"]',
  sanGroupToggle: '[data-test-toggle-group="Subject Alternative Name (SAN) Options"]',
  additionalGroupToggle: '[data-test-toggle-group="Additional subject fields"]',
  toggleGroupDescription: '[data-test-toggle-group-description]',
  formField: '[data-test-field]',
  typeField: '[data-test-input="type"]',
  inputByName: (name) => `[data-test-input="${name}"]`,
  fieldByName: (name) => `[data-test-field="${name}"]`,
  generateRootSave: '[data-test-pki-generate-root-save]',
  generateRootCancel: '[data-test-pki-generate-root-cancel]',
  formInvalidError: '[data-test-pki-generate-root-validation-error]',
  urlsSection: '[data-test-urls-section]',
  urlField: '[data-test-urls-section] [data-test-input]',
};
