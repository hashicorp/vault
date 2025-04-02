/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { click, visit, fillIn } from '@ember/test-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { AUTH_FORM } from 'vault/tests/helpers/auth/auth-form-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { MFA_SELECTORS } from 'vault/tests/helpers/mfa/mfa-selectors';
import { constraintId, setupTotpMfaResponse } from 'vault/tests/helpers/mfa/mfa-helpers';
import { fillInLoginFields } from 'vault/tests/helpers/auth/auth-helpers';
import { callbackData, windowStub } from 'vault/tests/helpers/oidc-window-stub';

const ENT_ONLY = ['saml'];

// See AUTH_METHOD_TEST_CASES for how request data maps to method types
// authRequest is the request made on submit and what returns mfa_validation requirements (if any)
// additionalRequest are any third party requests the auth method expects
const REQUEST_DATA = {
  username: {
    loginData: { username: 'matilda', password: 'password' },
    stubRequests: (server, path) =>
      server.post(`/auth/${path}/login/matilda`, () => setupTotpMfaResponse(path)),
  },
  github: {
    loginData: { token: 'mysupersecuretoken' },
    stubRequests: (server, path) => server.post(`/auth/${path}/login`, () => setupTotpMfaResponse(path)),
  },
  oidc: {
    loginData: { role: 'some-dev' },
    hasPopupWindow: true,
    stubRequests: (server, path) => {
      server.get(`/auth/${path}/oidc/callback`, () => setupTotpMfaResponse(path));
      server.post(`/auth/${path}/oidc/auth_url`, () => ({
        data: { auth_url: 'http://dev-foo-bar.com' },
      }));
    },
  },
  saml: {
    loginData: { role: 'some-dev' },
    hasPopupWindow: true,
    stubRequests: (server, path) => {
      server.put(`/auth/${path}/token`, () => setupTotpMfaResponse(path));
      server.put(`/auth/${path}/sso_service_url`, () => ({
        data: { sso_service_url: 'http://sso-url.hashicorp.com/service', token_poll_id: '1234' },
      }));
    },
  },
};

// maps auth type to request data (line breaks to help separate and clarify which methods share request paths)
const AUTH_METHOD_TEST_CASES = [
  { authType: 'github', options: REQUEST_DATA.github },

  { authType: 'userpass', options: REQUEST_DATA.username },
  { authType: 'ldap', options: REQUEST_DATA.username },
  { authType: 'okta', options: REQUEST_DATA.username },
  { authType: 'radius', options: REQUEST_DATA.username },

  { authType: 'oidc', options: REQUEST_DATA.oidc },
  { authType: 'jwt', options: REQUEST_DATA.oidc },

  // ENTERPRISE ONLY
  { authType: 'saml', options: REQUEST_DATA.saml },
];

for (const method of AUTH_METHOD_TEST_CASES) {
  const { authType, options } = method;
  const isEntMethod = ENT_ONLY.includes(authType);
  // adding "enterprise" to the module title filters it out of the test runner for the CE repo
  module(`Acceptance | auth | mfa ${authType}${isEntMethod ? ' enterprise' : ''}`, function (hooks) {
    setupApplicationTest(hooks);
    setupMirage(hooks);

    hooks.beforeEach(async function () {
      if (options?.hasPopupWindow) {
        this.windowStub = windowStub();
      }
      await visit('/vault/auth');
    });

    hooks.afterEach(function () {
      if (options?.hasPopupWindow) {
        this.windowStub.restore();
      }
    });

    test(`${authType}: it displays mfa requirement for default paths`, async function (assert) {
      this.mountPath = authType;
      options.stubRequests(this.server, this.mountPath);

      const loginKeys = Object.keys(options.loginData);
      assert.expect(3 + loginKeys.length);

      // Fill in login form
      await fillIn(AUTH_FORM.method, authType);
      await fillInLoginFields(options.loginData);

      if (options?.hasPopupWindow) {
        // fires "message" event which methods that rely on popup windows wait for
        setTimeout(() => {
          // set path which is used to set :mount param in the callback url => /auth/:mount/oidc/callback
          window.postMessage(callbackData({ path: this.mountPath }), window.origin);
        }, 50);
      }

      await click(AUTH_FORM.login);
      assert
        .dom(MFA_SELECTORS.mfaForm)
        .hasText(
          'Multi-factor authentication is enabled for your account. Enter your authentication code to log in. TOTP passcode Verify'
        );
      await click(GENERAL.backButton);
      assert.dom(AUTH_FORM.form).exists('clicking back returns to auth form');
      assert.dom(GENERAL.selectByAttr('auth-method')).hasValue(authType, 'preserves method type on back');
      for (const field of loginKeys) {
        assert.dom(AUTH_FORM.input(field)).hasValue('', `${field} input clears on back`);
      }
    });

    test(`${authType}: it displays mfa requirement for custom paths`, async function (assert) {
      this.mountPath = `${authType}-custom`;
      options.stubRequests(this.server, this.mountPath);
      const loginKeys = Object.keys(options.loginData);
      assert.expect(3 + loginKeys.length);

      // Fill in login form
      await fillIn(AUTH_FORM.method, authType);
      // Toggle more options to input a custom mount path
      await fillInLoginFields(
        { ...options.loginData, 'auth-form-mount-path': this.mountPath },
        { toggleOptions: true }
      );

      if (options?.hasPopupWindow) {
        // fires "message" event which methods that rely on popup windows wait for
        setTimeout(() => {
          // set path which is used to set :mount param in the callback url => /auth/:mount/oidc/callback
          window.postMessage(callbackData({ path: this.mountPath }), window.origin);
        }, 50);
      }

      await click(AUTH_FORM.login);
      assert
        .dom(MFA_SELECTORS.mfaForm)
        .hasText(
          'Multi-factor authentication is enabled for your account. Enter your authentication code to log in. TOTP passcode Verify'
        );
      await click(GENERAL.backButton);
      assert.dom(AUTH_FORM.form).exists('clicking back returns to auth form');
      assert.dom(GENERAL.selectByAttr('auth-method')).hasValue(authType, 'preserves method type on back');
      for (const field of loginKeys) {
        assert.dom(AUTH_FORM.input(field)).hasValue('', `${field} input clears on back`);
      }
    });

    test(`${authType}: it submits mfa requirement for default paths`, async function (assert) {
      assert.expect(2);
      this.mountPath = authType;
      options.stubRequests(this.server, this.mountPath);

      const expectedOtp = '12345';
      server.post('/sys/mfa/validate', async (_, req) => {
        const [actualOtp] = JSON.parse(req.requestBody).mfa_payload[constraintId];
        assert.true(true, 'it makes request to mfa validate endpoint');
        assert.strictEqual(actualOtp, expectedOtp, 'payload contains otp');
      });

      // Fill in login form
      await fillIn(AUTH_FORM.method, authType);
      await fillInLoginFields(options.loginData);

      if (options?.hasPopupWindow) {
        // fires "message" event which methods that rely on popup windows wait for
        setTimeout(() => {
          // set path which is used to set :mount param in the callback url => /auth/:mount/oidc/callback
          window.postMessage(callbackData({ path: this.mountPath }), window.origin);
        }, 50);
      }

      await click(AUTH_FORM.login);
      await fillIn(MFA_SELECTORS.passcode(0), expectedOtp);
      await click(MFA_SELECTORS.validate);
    });

    test(`${authType}: it submits mfa requirement for custom paths`, async function (assert) {
      assert.expect(2);

      this.mountPath = `${authType}-custom`;
      options.stubRequests(this.server, this.mountPath);

      const expectedOtp = '12345';
      server.post('/sys/mfa/validate', async (_, req) => {
        const [actualOtp] = JSON.parse(req.requestBody).mfa_payload[constraintId];
        assert.true(true, 'it makes request to mfa validate endpoint');
        assert.strictEqual(actualOtp, expectedOtp, 'payload contains otp');
      });

      // Fill in login form
      await fillIn(AUTH_FORM.method, authType);
      // Toggle more options to input a custom mount path
      await fillInLoginFields(
        { ...options.loginData, 'auth-form-mount-path': this.mountPath },
        { toggleOptions: true }
      );

      if (options?.hasPopupWindow) {
        // fires "message" event which methods that rely on popup windows wait for
        setTimeout(() => {
          // set path which is used to set :mount param in the callback url => /auth/:mount/oidc/callback
          window.postMessage(callbackData({ path: this.mountPath }), window.origin);
        }, 50);
      }

      await click(AUTH_FORM.login);
      await fillIn(MFA_SELECTORS.passcode(0), expectedOtp);
      await click(MFA_SELECTORS.validate);
    });
  });
}
