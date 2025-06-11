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
import { AUTH_METHOD_MAP, fillInLoginFields } from 'vault/tests/helpers/auth/auth-helpers';
import {
  callbackData,
  DELAY_IN_MS,
  triggerMessageEvent,
  windowStub,
} from 'vault/tests/helpers/oidc-window-stub';

const ENT_ONLY = ['saml'];

for (const method of AUTH_METHOD_MAP) {
  const { authType, options } = method;
  // token doesn't support MFA
  if (authType === 'token') continue;

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
      options.stubRequests(this.server, this.mountPath, setupTotpMfaResponse(this.mountPath));

      const loginKeys = Object.keys(options.loginData);
      assert.expect(3 + loginKeys.length);

      // Fill in login form
      await fillIn(AUTH_FORM.selectMethod, authType);
      await fillInLoginFields(options.loginData);

      if (options?.hasPopupWindow) {
        // fires "message" event which methods that rely on popup windows wait for
        // pass mount path which is used to set :mount param in the callback url => /auth/:mount/oidc/callback
        triggerMessageEvent(this.mountPath);
      }

      await click(GENERAL.submitButton);
      assert
        .dom(MFA_SELECTORS.mfaForm)
        .hasText(
          'Back Multi-factor authentication is enabled for your account. Enter your authentication code to log in. TOTP passcode Verify'
        );
      await click(GENERAL.backButton);
      assert.dom(AUTH_FORM.form).exists('clicking back returns to auth form');
      assert.dom(AUTH_FORM.selectMethod).hasValue(authType, 'preserves method type on back');
      for (const field of loginKeys) {
        assert.dom(GENERAL.inputByAttr(field)).hasValue('', `${field} input clears on back`);
      }
    });

    test(`${authType}: it displays mfa requirement for custom paths`, async function (assert) {
      this.mountPath = `${authType}-custom`;
      options.stubRequests(this.server, this.mountPath, setupTotpMfaResponse(this.mountPath));
      const loginKeys = Object.keys(options.loginData);
      assert.expect(3 + loginKeys.length);

      // Fill in login form
      await fillIn(AUTH_FORM.selectMethod, authType);
      // Toggle more options to input a custom mount path
      await fillInLoginFields({ ...options.loginData, path: this.mountPath }, { toggleOptions: true });

      if (options?.hasPopupWindow) {
        // fires "message" event which methods that rely on popup windows wait for
        setTimeout(() => {
          // set path which is used to set :mount param in the callback url => /auth/:mount/oidc/callback
          window.postMessage(callbackData({ path: this.mountPath }), window.origin);
        }, DELAY_IN_MS);
      }

      await click(GENERAL.submitButton);
      assert
        .dom(MFA_SELECTORS.mfaForm)
        .hasText(
          'Back Multi-factor authentication is enabled for your account. Enter your authentication code to log in. TOTP passcode Verify'
        );
      await click(GENERAL.backButton);
      assert.dom(AUTH_FORM.form).exists('clicking back returns to auth form');
      assert.dom(AUTH_FORM.selectMethod).hasValue(authType, 'preserves method type on back');
      for (const field of loginKeys) {
        assert.dom(GENERAL.inputByAttr(field)).hasValue('', `${field} input clears on back`);
      }
    });

    test(`${authType}: it submits mfa requirement for default paths`, async function (assert) {
      assert.expect(2);
      this.mountPath = authType;
      options.stubRequests(this.server, this.mountPath, setupTotpMfaResponse(this.mountPath));

      const expectedOtp = '12345';
      server.post('/sys/mfa/validate', async (_, req) => {
        const [actualOtp] = JSON.parse(req.requestBody).mfa_payload[constraintId];
        assert.true(true, 'it makes request to mfa validate endpoint');
        assert.strictEqual(actualOtp, expectedOtp, 'payload contains otp');
      });

      // Fill in login form
      await fillIn(AUTH_FORM.selectMethod, authType);
      await fillInLoginFields(options.loginData);

      if (options?.hasPopupWindow) {
        // fires "message" event which methods that rely on popup windows wait for
        setTimeout(() => {
          // set path which is used to set :mount param in the callback url => /auth/:mount/oidc/callback
          window.postMessage(callbackData({ path: this.mountPath }), window.origin);
        }, DELAY_IN_MS);
      }

      await click(GENERAL.submitButton);
      await fillIn(MFA_SELECTORS.passcode(0), expectedOtp);
      await click(MFA_SELECTORS.validate);
    });

    test(`${authType}: it submits mfa requirement for custom paths`, async function (assert) {
      assert.expect(2);

      this.mountPath = `${authType}-custom`;
      options.stubRequests(this.server, this.mountPath, setupTotpMfaResponse(this.mountPath));

      const expectedOtp = '12345';
      server.post('/sys/mfa/validate', async (_, req) => {
        const [actualOtp] = JSON.parse(req.requestBody).mfa_payload[constraintId];
        assert.true(true, 'it makes request to mfa validate endpoint');
        assert.strictEqual(actualOtp, expectedOtp, 'payload contains otp');
      });

      // Fill in login form
      await fillIn(AUTH_FORM.selectMethod, authType);
      // Toggle more options to input a custom mount path
      await fillInLoginFields({ ...options.loginData, path: this.mountPath }, { toggleOptions: true });

      if (options?.hasPopupWindow) {
        // fires "message" event which methods that rely on popup windows wait for
        setTimeout(() => {
          // set path which is used to set :mount param in the callback url => /auth/:mount/oidc/callback
          window.postMessage(callbackData({ path: this.mountPath }), window.origin);
        }, DELAY_IN_MS);
      }

      await click(GENERAL.submitButton);
      await fillIn(MFA_SELECTORS.passcode(0), expectedOtp);
      await click(MFA_SELECTORS.validate);
    });
  });
}
