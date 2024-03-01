/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { SELECTORS as CONFIGURE } from './pki-configure-create';
import { SELECTORS as DETAILS } from './pki-issuer-details';

export const SELECTORS = {
  objectListInput: (key, row = 0) => `[data-test-object-list-input="${key}-${row}"]`,
  inputByName: (name) => `[data-test-input="${name}"]`,
  addRow: '[data-test-object-list-add-button',
  submitButton: '[data-test-cross-sign-submit]',
  cancelButton: '[data-test-cross-sign-cancel]',
  statusCount: '[data-test-cross-sign-status-count]',
  signedIssuerRow: (row = 0) => `[data-test-info-table-row="${row}"]`,
  signedIssuerCol: (attr) => `[data-test-info-table-column="${attr}"]`,
  // for cross signing acceptance tests
  configure: { ...CONFIGURE },
  details: { ...DETAILS },
  rowValue: (attr) => `[data-test-row-value="${attr}"]`,
  copyButton: (attr) => `[data-test-value-div="${attr}"] [data-test-copy-button]`,
};
