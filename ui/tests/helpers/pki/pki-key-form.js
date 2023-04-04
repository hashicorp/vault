/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

export const SELECTORS = {
  keyCreateButton: '[data-test-pki-key-save]',
  keyCancelButton: '[data-test-pki-key-cancel]',
  keyNameInput: '[data-test-input="keyName"]',
  typeInput: '[data-test-input="type"]',
  keyTypeInput: '[data-test-input="keyType"]',
  keyBitsInput: '[data-test-input="keyBits"]',
  validationError: '[data-test-pki-key-validation-error]',
  fieldErrorByName: (name) => `[data-test-field-validation="${name}"]`,
};
