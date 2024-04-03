/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export const PKI_CONFIGURE_CREATE = {
  // page::pki-configure-create
  nextStepsBanner: '[data-test-config-next-steps]',
  option: '[data-test-pki-config-option]',
  optionByKey: (key) => `[data-test-pki-config-option="${key}"]`,
  doneButton: '[data-test-done]',
  configureButton: '[data-test-configure-pki-button]',
  // pki-generate-root
  generateRootOption: '[data-test-pki-config-option="generate-root"]',
  // pki-ca-cert-import
  importForm: '[data-test-pki-import-pem-bundle-form]',
  importSubmit: '[data-test-pki-import-pem-bundle]',
  importSectionLabel: '[data-test-import-section-label]',
  importMapping: '[data-test-imported-bundle-mapping]',
  importedIssuer: '[data-test-imported-issuer]',
  importedKey: '[data-test-imported-key]',
  // generate-intermediate
  csrDetails: '[data-test-generate-csr-result]',
};
