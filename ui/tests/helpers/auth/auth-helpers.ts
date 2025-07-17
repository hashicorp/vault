/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, fillIn, visit } from '@ember/test-helpers';
import VAULT_KEYS from 'vault/tests/helpers/vault-keys';
import { AUTH_FORM } from 'vault/tests/helpers/auth/auth-form-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { Server } from 'miragejs';

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
  return click(AUTH_FORM.login);
};

export const loginNs = async (ns: string, token = rootToken) => {
  // make sure we're always logged out and logged back in
  await logout();
  await visit('/vault/auth');

  await fillIn(GENERAL.inputByAttr('namespace'), ns);

  await fillIn(AUTH_FORM.selectMethod, 'token');
  await fillIn(GENERAL.inputByAttr('token'), token);
  return click(AUTH_FORM.login);
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
  return click(AUTH_FORM.login);
};

export const fillInLoginFields = async (loginFields: LoginFields, { toggleOptions = false } = {}) => {
  if (toggleOptions) await click(AUTH_FORM.advancedSettings);

  for (const [input, value] of Object.entries(loginFields)) {
    if (value) {
      await fillIn(GENERAL.inputByAttr(input), value);
    }
  }
};

// See AUTH_METHOD_MAP for how login data maps to method types,
// stubRequests are the requests made on submit for that method type
export const LOGIN_DATA = {
  token: {
    loginData: { token: 'mytoken' },
    stubRequests: (server: Server, response: object) => server.get('/auth/token/lookup-self', () => response),
  },
  username: {
    loginData: { username: 'matilda', password: 'password' },
    stubRequests: (server: Server, path: string, response: object) =>
      server.post(`/auth/${path}/login/matilda`, () => response),
  },
  github: {
    loginData: { token: 'mysupersecuretoken' },
    stubRequests: (server: Server, path: string, response: object) =>
      server.post(`/auth/${path}/login`, () => response),
  },
  oidc: {
    loginData: { role: 'some-dev' },
    hasPopupWindow: true,
    stubRequests: (server: Server, path: string, response: object) => {
      server.get(`/auth/${path}/oidc/callback`, () => response);
      server.post(`/auth/${path}/oidc/auth_url`, () => {
        return { data: { auth_url: 'http://dev-foo-bar.com' } };
      });
    },
  },
  saml: {
    loginData: { role: 'some-dev' },
    hasPopupWindow: true,
    stubRequests: (server: Server, path: string, response: object) => {
      server.put(`/auth/${path}/token`, () => response);
      server.put(`/auth/${path}/sso_service_url`, () => {
        return { data: { sso_service_url: 'http://sso-url.hashicorp.com/service', token_poll_id: '1234' } };
      });
    },
  },
};

// maps auth type to request data
export const AUTH_METHOD_MAP = [
  { authType: 'token', options: LOGIN_DATA.token },
  { authType: 'github', options: LOGIN_DATA.github },

  // username and password methods
  { authType: 'userpass', options: LOGIN_DATA.username },
  { authType: 'ldap', options: LOGIN_DATA.username },
  { authType: 'okta', options: LOGIN_DATA.username },
  { authType: 'radius', options: LOGIN_DATA.username },

  // oidc
  { authType: 'oidc', options: LOGIN_DATA.oidc },
  { authType: 'jwt', options: LOGIN_DATA.oidc },

  // ENTERPRISE ONLY
  { authType: 'saml', options: LOGIN_DATA.saml },
];

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
