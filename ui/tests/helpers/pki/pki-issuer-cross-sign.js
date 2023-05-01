/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { SELECTORS as CONFIGURE } from './pki-configure-create';

export const SELECTORS = {
  input: (key, row = 0) => `[data-test-object-list-input="${key}-${row}"]`,
  addRow: '[data-test-object-list-add-button',
  submitButton: '[data-test-cross-sign-submit]',
  cancelButton: '[data-test-cross-sign-cancel]',
  statusCount: '[data-test-cross-sign-status-count]',
  signedIssuerRow: (row = 0) => `[data-test-info-table-row="${row}"]`,
  signedIssuerCol: (attr) => `[data-test-info-table-column="${attr}"]`,
  // for cross signing acceptance tests
  ...CONFIGURE,
  rowValue: (attr) => `[data-test-row-value="${attr}"]`,
};
