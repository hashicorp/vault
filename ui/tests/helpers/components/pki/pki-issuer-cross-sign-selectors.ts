/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export const PKI_CROSS_SIGN = {
  objectListInput: (key: string, row = 0) => `[data-test-object-list-input="${key}-${row}"]`,
  addRow: '[data-test-object-list-add-button',
  statusCount: '[data-test-cross-sign-status-count]',
  signedIssuerRow: (row = 0) => `[data-test-info-table-row="${row}"]`,
  signedIssuerCol: (attr: string) => `[data-test-info-table-column="${attr}"]`,
  rowValue: (attr: string) => `[data-test-row-value="${attr}"]`,
  copyButton: (attr: string) => `[data-test-value-div="${attr}"] [data-test-copy-button]`,
};
