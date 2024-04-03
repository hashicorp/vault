/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export const PKI_CROSS_SIGN = {
  objectListInput: (key, row = 0) => `[data-test-object-list-input="${key}-${row}"]`,
  inputByName: (name) => `[data-test-input="${name}"]`, // GENERAL.inputByAttr()
  addRow: '[data-test-object-list-add-button',
  submitButton: '[data-test-save]',
  cancelButton: '[data-test-cancel]',
  statusCount: '[data-test-cross-sign-status-count]',
  signedIssuerRow: (row = 0) => `[data-test-info-table-row="${row}"]`,
  signedIssuerCol: (attr) => `[data-test-info-table-column="${attr}"]`,
  rowValue: (attr) => `[data-test-row-value="${attr}"]`,
  copyButton: (attr) => `[data-test-value-div="${attr}"] [data-test-copy-button]`,
};
