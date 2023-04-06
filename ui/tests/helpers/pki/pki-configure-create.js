/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { SELECTORS as GENERATE_ROOT } from './pki-generate-root';

export const SELECTORS = {
  // page::pki-configure-create
  breadcrumbContainer: '[data-test-breadcrumbs]',
  title: '[data-test-pki-engine-page-title]',
  option: '[data-test-pki-config-option]',
  optionByKey: (key) => `[data-test-pki-config-option="${key}"]`,
  cancelButton: '[data-test-pki-config-cancel]',
  saveButton: '[data-test-pki-config-save]',
  doneButton: '[data-test-done]',
  // pki-generate-root
  ...GENERATE_ROOT,
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
