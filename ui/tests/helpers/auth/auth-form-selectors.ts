/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export const AUTH_FORM = {
  description: '[data-test-description]',
  form: '[data-test-auth-form]',
  linkedBlockAuth: (type: string) => `[data-test-auth-backend-link="${type}"]`,
  login: '[data-test-auth-submit]',
  selectMethod: '[data-test-select="auth type"]',
  tabBtn: (method: string) => `[data-test-auth-tab="${method}"] button`, // method is all lowercased
  tabs: '[data-test-auth-tab]',
  // old form toggle, will eventually be deleted
  moreOptions: '[data-test-auth-form-options-toggle]',
  // new toggle, hds component is a button
  advancedSettings: '[data-test-auth-form-options-toggle] button',
  authForm: (type: string) => `[data-test-auth-form="${type}"]`,
  helpText: '[data-test-auth-helptext]',
  logo: '[data-test-auth-logo]',
  managedNsRoot: '[data-test-managed-namespace-root]',
  otherMethodsBtn: '[data-test-other-methods-button]',
};
