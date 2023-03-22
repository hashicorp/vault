/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { SELECTORS as GENERATE_ROOT } from './pki-generate-root';

export const SELECTORS = {
  // pki-configure-form
  option: '[data-test-pki-config-option]',
  optionByKey: (key) => `[data-test-pki-config-option="${key}"]`,
  cancelButton: '[data-test-pki-config-cancel]',
  saveButton: '[data-test-pki-config-save]',
  // pki-generate-root
  ...GENERATE_ROOT,
  // pki-ca-cert-import
  importForm: '[data-test-pki-import-pem-bundle-form]',
  // generate-intermediate
  csrDetails: '[data-test-generate-csr-result]',
};
