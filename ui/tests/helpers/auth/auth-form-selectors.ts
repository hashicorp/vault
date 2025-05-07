/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export const AUTH_FORM = {
  selectMethod: '[data-test-select="auth type"]',
  form: '[data-test-auth-form]',
  login: '[data-test-auth-submit]',
  // TODO RENAME
  tabs: (method: string) => (method ? `[data-test-auth-method="${method}"]` : '[data-test-auth-method]'),
  method: (method: string) => `[data-test-auth-method="${method}"]`,
  tabBtn: (method: string) => `[data-test-auth-method="${method}"] button`,
  description: '[data-test-description]',
  // old form toggle, will eventually be deleted
  moreOptions: '[data-test-auth-form-options-toggle]',
  // new toggle, hds component is a button
  advancedSettings: '[data-test-auth-form-options-toggle] button',
  managedNsRoot: '[data-test-managed-namespace-root]',
  logo: '[data-test-auth-logo]',
  helpText: '[data-test-auth-helptext]',
  authForm: (type: string) => `[data-test-auth-form="${type}"]`,
  otherMethodsBtn: '[data-test-other-methods-button]',
};
