/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export const PKI_GENERATE_ROOT = {
  mainSectionTitle: '[data-test-generate-root-title="Root parameters"]',
  urlSectionTitle: '[data-test-generate-root-title="Issuer URLs"]',
  keyParamsGroupToggle: '[data-test-toggle-group="Key parameters"]',
  sanGroupToggle: '[data-test-toggle-group="Subject Alternative Name (SAN) Options"]',
  additionalGroupToggle: '[data-test-toggle-group="Additional subject fields"]',
  toggleGroupDescription: '[data-test-toggle-group-description]',
  groupFields: (group) => `[data-test-group="${group}"] [data-test-field]`,
  formInvalidError: '[data-test-pki-generate-root-validation-error]',
  urlsSection: '[data-test-urls-section]',
  urlField: '[data-test-urls-section] [data-test-input]',
  // Shown values after save
  saved: {
    certificate: '[data-test-value-div="Certificate"] [data-test-certificate-card]',
    commonName: '[data-test-row-value="Common name"]',
    issuerName: '[data-test-row-value="Issuer name"]',
    issuerLink: '[data-test-value-div="Issuer ID"] a',
    keyName: '[data-test-row-value="Key name"]',
    keyLink: '[data-test-value-div="Key ID"] a',
    privateKey: '[data-test-value-div="Private key"] [data-test-certificate-card]',
    serialNumber: '[data-test-row-value="Serial number"]',
  },
};
