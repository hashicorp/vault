/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export const AUTH_FORM = {
  method: '[data-test-select=auth-method]',
  form: '[data-test-auth-form]',
  login: '[data-test-auth-submit]',
  tabs: (method: string) => (method ? `[data-test-auth-method="${method}"]` : '[data-test-auth-method]'),
  description: '[data-test-description]',
  roleInput: '[data-test-role]',
  input: (item: string) => `[data-test-${item}]`, // i.e. jwt, role, token, password or username
  mountPathInput: '[data-test-auth-form-mount-path]',
  moreOptions: '[data-test-auth-form-options-toggle]',
  namespaceInput: '[data-test-auth-form-ns-input]',
};
