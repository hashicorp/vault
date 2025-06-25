/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, fillIn, visit } from '@ember/test-helpers';
import VAULT_KEYS from 'vault/tests/helpers/vault-keys';
import { AUTH_FORM } from 'vault/tests/helpers/auth/auth-form-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

import type { LoginFields } from 'vault/vault/auth/form';

export const { rootToken } = VAULT_KEYS;

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
  await visit('/vault/auth');

  await fillIn(AUTH_FORM.selectMethod, 'token');
  await fillIn(GENERAL.inputByAttr('token'), token);
  return click(GENERAL.submitButton);
};

export const loginNs = async (ns: string, token = rootToken) => {
  // make sure we're always logged out and logged back in
  await logout();
  await visit('/vault/auth');

  await fillIn(GENERAL.inputByAttr('namespace'), ns);

  await fillIn(AUTH_FORM.selectMethod, 'token');
  await fillIn(GENERAL.inputByAttr('token'), token);
  return click(GENERAL.submitButton);
};

// LOGIN WITH NON-TOKEN METHODS
export const loginMethod = async (
  loginFields: LoginFields,
  options: { authType?: string; toggleOptions?: boolean }
) => {
  // make sure we're always logged out and logged back in
  await logout();
  const type = options?.authType || 'token';

  await fillIn(AUTH_FORM.selectMethod, type);

  await fillInLoginFields(loginFields, options);
  return click(GENERAL.submitButton);
};

export const fillInLoginFields = async (loginFields: LoginFields, { toggleOptions = false } = {}) => {
  if (toggleOptions) await click(AUTH_FORM.advancedSettings);

  for (const [input, value] of Object.entries(loginFields)) {
    if (value) {
      await fillIn(GENERAL.inputByAttr(input), value);
    }
  }
};

const LOGIN_DATA = {
  token: { token: 'mysupersecuretoken' },
  username: { username: 'matilda', password: 'password' },
  role: { role: 'some-dev' },
};
// maps auth type to login input data
export const AUTH_METHOD_LOGIN_DATA = {
  // token methods
  token: LOGIN_DATA.token,
  github: LOGIN_DATA.token,
  // username and password methods
  userpass: LOGIN_DATA.username,
  ldap: LOGIN_DATA.username,
  okta: LOGIN_DATA.username,
  radius: LOGIN_DATA.username,
  // role
  oidc: LOGIN_DATA.role,
  jwt: LOGIN_DATA.role,
  saml: LOGIN_DATA.role,
};

// Mock response for `sys/internal/ui/mounts`
export const SYS_INTERNAL_UI_MOUNTS = {
  'userpass/': {
    description: '',
    options: {},
    type: 'userpass',
  },
  'userpass2/': {
    description: '',
    options: {},
    type: 'userpass',
  },
  // there was a problem with the API service camel-casing mounts that were snake cased
  // so including a snake cased mount for testing
  'my_oidc/': {
    description: '',
    options: {},
    type: 'oidc',
  },
  'ldap/': {
    description: '',
    options: null,
    type: 'ldap',
  },
};
