/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export const SELECTORS = {
  errorBanner: '[data-test-config-edit-error]',
  acmeEditSection: '[data-test-acme-edit-section]',
  configEditSection: '[data-test-cluster-config-edit-section]',
  configInput: (attr) => `[data-test-input="${attr}"]`,
  stringListInput: (attr) => `[data-test-input="${attr}"] [data-test-string-list-input="0"]`,
  urlsEditSection: '[data-test-urls-edit-section]',
  urlFieldInput: (attr) => `[data-test-input="${attr}"] textarea`,
  urlFieldLabel: (attr) => `[data-test-input="${attr}"] label`,
  crlEditSection: '[data-test-crl-edit-section]',
  crlToggleInput: (attr) => `[data-test-input="${attr}"] input`,
  crlTtlInput: (attr) => `[data-test-ttl-value="${attr}"]`,
  crlFieldLabel: (attr) => `[data-test-input="${attr}"] label`,
  saveButton: '[data-test-configuration-edit-save]',
  cancelButton: '[data-test-configuration-edit-cancel]',
  validationAlert: '[data-test-configuration-edit-validation-alert]',
  deleteButton: (attr) => `[data-test-input="${attr}"] [data-test-string-list-button="delete"]`,
  groupHeader: (group) => `[data-test-crl-header="${group}"]`,
  checkboxInput: (attr) => `[data-test-input="${attr}"]`,
};
