/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, fillIn, visit } from '@ember/test-helpers';
import VAULT_KEYS from 'vault/tests/helpers/vault-keys';
import { AUTH_FORM } from 'vault/tests/helpers/auth/auth-form-selectors';

const { rootToken } = VAULT_KEYS;

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

// LoginFields completes the input's test selector and fills it in with the corresponding value
interface LoginFields {
  username?: string;
  password?: string;
  token?: string;
  role?: string;
  'auth-form-mount-path': string; // todo update selectors
  'auth-form-ns-input': string; // todo update selectors
}
interface LoginOptions {
  authType?: string; // TODO not used right now
  toggleOptions?: boolean;
}

export const fillInLoginFields = async (loginFields: LoginFields, options: LoginOptions) => {
  const { toggleOptions = false, authType = 'token' } = options;

  await fillIn(AUTH_FORM.method, authType);

  if (toggleOptions) await click(AUTH_FORM.moreOptions);

  for (const [input, value] of Object.entries(loginFields)) {
    await fillIn(AUTH_FORM.input(input), value);
  }
};

export const loginMethod = async (loginFields: LoginFields, options: LoginOptions) => {
  // make sure we're always logged out and logged back in
  await logout();
  await visit(`/vault/auth?with=${options.authType}`);

  await fillInLoginFields(loginFields, options);
  return click(AUTH_FORM.login);
};

export const logout = async () => {
  // make sure we're always logged out and logged back in
  await visit('/vault/logout');
  // clear session storage to ensure we have a clean state
  window.localStorage.clear();
  return;
};

interface ResponseData {
  authType?: string;
  authMountPath?: string;
  username?: string;
  isMfa?: boolean;
}
export const authRequest = (context: any, options: ResponseData) => {
  const { isMfa = false, authMountPath = '', username } = options;
  let warnings = null;
  let mfa_requirement = null;
  if (isMfa) {
    warnings = [
      'A login request was issued that is subject to MFA validation. Please make sure to validate the login by sending another request to mfa/validate endpoint.',
    ];

    // in the real world more info is returned by this request, only including pertinent data for testing
    mfa_requirement = {
      mfa_request_id: '4c85205e-a946-bb01-be91-b2420e8c0822',
      mfa_constraints: {
        [authMountPath]: {
          any: [
            {
              type: 'totp',
              id: '6a55fcb1-efa9-eb89-ad92-5db8d0fd9c02',
              uses_passcode: true,
            },
          ],
        },
      },
    };
  }
  // in the real world more info is returned by this request, only including pertinent data for testing
  const authResponse = {
    policies: ['default'],
    mfa_requirement,
  };

  return context.server.post(`/auth/${authMountPath}/login/${username}`, () => {
    return {
      request_id: '51b49e25-e980-55a2-76ba-4c690adcc0c3',
      warnings,
      auth: authResponse,
    };
  });
};
