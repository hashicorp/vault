/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export const MFA_SELECTORS = {
  mfaForm: '[data-test-mfa-form]',
  passcode: (idx: number) => `[data-test-mfa-passcode="${idx}"]`,
  validate: '[data-test-mfa-validate]',
};
