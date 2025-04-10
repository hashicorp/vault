/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, fillIn, visit } from '@ember/test-helpers';
import VAULT_KEYS from 'vault/tests/helpers/vault-keys';
import { AUTH_FORM } from 'vault/tests/helpers/auth/auth-form-selectors';

const { rootToken } = VAULT_KEYS;

// LOGOUT
export const logout = async () => {
  // make sure we're always logged out and logged back in
  await visit('/vault/logout');
  // clear session storage to ensure we have a clean state
  window.localStorage.clear();
  return;
};

// LOGIN WITH TOKEN
export const login = async (token = rootToken) => {
  // make sure we're always logged out and logged back in
  await logout();
  await visit('/vault/auth?with=token');
  await fillIn(AUTH_FORM.input('token'), token);
  return click(AUTH_FORM.login);
};

export const loginNs = async (ns: string, token = rootToken) => {
  // make sure we're always logged out and logged back in
  await logout();
  await visit('/vault/auth?with=token');
  await fillIn(AUTH_FORM.namespaceInput, ns);
  await fillIn(AUTH_FORM.input('token'), token);
  return click(AUTH_FORM.login);
};

// LOGIN WITH NON-TOKEN methods
interface LoginOptions {
  authType?: string;
  toggleOptions?: boolean;
}
export const loginMethod = async (loginFields: LoginFields, options: LoginOptions) => {
  // make sure we're always logged out and logged back in
  await logout();
  await visit(`/vault/auth?with=${options.authType}`);

  await fillInLoginFields(loginFields, options);
  return click(AUTH_FORM.login);
};

// the keys complete the input's test selector and the helper fills the input with the corresponding value
interface LoginFields {
  username?: string;
  password?: string;
  token?: string;
  role?: string;
  'auth-form-mount-path': string; // todo update selectors
  'auth-form-ns-input': string; // todo update selectors
}

export const fillInLoginFields = async (loginFields: LoginFields, { toggleOptions = false } = {}) => {
  if (toggleOptions) await click(AUTH_FORM.moreOptions);

  for (const [input, value] of Object.entries(loginFields)) {
    await fillIn(AUTH_FORM.input(input), value);
  }
};
