/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

export const SHAMIR_FORM = {
  inputLabel: '[data-test-shamir-key-label]',
  flowStep: (step: string) => `[data-test-dr-token-flow-step="${step}"]`,
  otpInfo: '[data-test-otp-info]',
  otpCode: '[data-test-otp]',
  progress: '.shamir-progress',
};
