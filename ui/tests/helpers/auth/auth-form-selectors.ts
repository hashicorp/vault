/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export const AUTH_FORM = {
  method: '[data-test-select=auth-method]',
  form: '[data-test-auth-form]',
  login: '[data-test-auth-submit]',
  tabs: (method: string) => (method ? `[data-test-auth-method="${method}"]` : '[data-test-auth-method]'),
  tabBtn: (method: string) => `[data-test-auth-method="${method}"] button`,
  description: '[data-test-description]',
  roleInput: '[data-test-role]',
  input: (item: string) => `[data-test-${item}]`, // i.e. jwt, role, token, password or username
  mountPathInput: '[data-test-auth-form-mount-path]',
  moreOptions: '[data-test-auth-form-options-toggle]',
  advancedSettings: '[data-test-auth-form-options-toggle] button',
  namespaceInput: '[data-test-auth-form-ns-input]',
  managedNsRoot: '[data-test-managed-namespace-root]',
  logo: '[data-test-auth-logo]',
  helpText: '[data-test-auth-helptext]',
  authForm: (type: string) => `[data-test-auth-form="${type}"]`,
  otherMethodsBtn: '[data-test-other-methods-button]',
};
