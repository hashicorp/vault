/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, fillIn, visit } from '@ember/test-helpers';
import VAULT_KEYS from 'vault/tests/helpers/vault-keys';
import { AUTH_FORM } from 'vault/tests/helpers/auth/auth-form-selectors';

const { rootToken } = VAULT_KEYS;

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

export const logout = async () => {
  // make sure we're always logged out and logged back in
  await visit('/vault/logout');
  // clear session storage to ensure we have a clean state
  window.localStorage.clear();
  return;
};
