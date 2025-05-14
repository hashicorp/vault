/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export const AUTH_FORM = {
  selectMethod: '[data-test-select="auth type"]',
  form: '[data-test-auth-form]',
  login: '[data-test-auth-submit]',
  preferredMethod: (method: string) => `p[data-test-auth-method="${method}"]`,
  tabs: '[data-test-auth-tab]',
  tabBtn: (method: string) => `[data-test-auth-tab="${method}"] button`, // method is all lowercased
  description: '[data-test-description]',
  advancedSettings: '[data-test-advanced-settings] button',
  managedNsRoot: '[data-test-managed-namespace-root]',
  logo: '[data-test-auth-logo]',
  helpText: '[data-test-auth-helptext]',
  authForm: (type: string) => `[data-test-auth-form="${type}"]`,
  otherMethodsBtn: '[data-test-other-methods-button]',
};
