/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export const MFA_SELECTORS = {
  countdown: '[data-test-mfa-countdown]',
  description: '[data-test-mfa-description]',
  label: '[data-test-mfa-label]',
  mfaForm: '[data-test-mfa-form]',
  passcode: (idx: number) => (idx ? `[data-test-mfa-passcode="${idx}"]` : '[data-test-mfa-passcode]'),
  push: '[data-test-mfa-push-instruction]',
  verifyBadge: (label: string) => `[data-test-mfa-verified="${label}"]`,
  qrCode: '[data-test-qrcode]',
  select: (idx: number) => (idx ? `[data-test-mfa-select="${idx}"]` : '[data-test-mfa-select]'),
  subheader: '[data-test-mfa-subheader]',
  subtitle: '[data-test-mfa-subtitle]',
  verifyForm: 'form#mfa-form',
};
