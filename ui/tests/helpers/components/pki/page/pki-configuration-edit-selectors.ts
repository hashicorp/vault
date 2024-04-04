/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export const PKI_CONFIG_EDIT = {
  errorBanner: '[data-test-config-edit-error]',
  acmeEditSection: '[data-test-acme-edit-section]',
  configEditSection: '[data-test-cluster-config-edit-section]',
  configInput: (attr: string) => `[data-test-input="${attr}"]`,
  stringListInput: (attr: string) => `[data-test-input="${attr}"] [data-test-string-list-input="0"]`,
  urlsEditSection: '[data-test-urls-edit-section]',
  urlFieldInput: (attr: string) => `[data-test-input="${attr}"] textarea`,
  urlFieldLabel: (attr: string) => `[data-test-input="${attr}"] label`,
  crlEditSection: '[data-test-crl-edit-section]',
  crlToggleInput: (attr: string) => `[data-test-input="${attr}"] input`,
  crlTtlInput: (attr: string) => `[data-test-ttl-value="${attr}"]`,
  crlFieldLabel: (attr: string) => `[data-test-input="${attr}"] label`,
  saveButton: '[data-test-configuration-edit-save]',
  cancelButton: '[data-test-configuration-edit-cancel]',
  validationAlert: '[data-test-configuration-edit-validation-alert]',
  deleteButton: (attr: string) => `[data-test-input="${attr}"] [data-test-string-list-button="delete"]`,
  groupHeader: (group: string) => `[data-test-crl-header="${group}"]`,
  checkboxInput: (attr: string) => `[data-test-input="${attr}"]`,
};
