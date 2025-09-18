/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { constraintId, setupTotpMfaResponse } from 'vault/tests/helpers/mfa/mfa-helpers';
import setupTestContext from './setup-test-context';
import { ERROR_JWT_LOGIN } from 'vault/utils/auth-form-helpers';
import { overrideResponse } from 'vault/tests/helpers/stubs';
import sinon from 'sinon';
import { triggerMessageEvent, windowStub } from 'vault/tests/helpers/oidc-window-stub';
import { AUTH_FORM } from 'vault/tests/helpers/auth/auth-form-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { MFA_SELECTORS } from 'vault/tests/helpers/mfa/mfa-selectors';
import { fillInLoginFields } from 'vault/tests/helpers/auth/auth-helpers';
import { click, fillIn, waitFor } from '@ember/test-helpers';

const mfaTests = (test) => {
  test('it displays mfa requirement for default paths', async function (assert) {
    const loginKeys = Object.keys(this.loginData);
    assert.expect(3 + loginKeys.length);
    this.stubRequests();
    await this.renderComponent();

    // Fill in login form
    await fillIn(AUTH_FORM.selectMethod, this.authType);
    await fillInLoginFields(this.loginData);

    if (this.authType === 'oidc') {
      // fires "message" event which methods that rely on popup windows wait for
      // pass mount path which is used to set :mount param in the callback url => /auth/:mount/oidc/callback
      triggerMessageEvent(this.path);
    }

    await click(GENERAL.submitButton);
    await waitFor(MFA_SELECTORS.mfaForm);
    assert
      .dom(MFA_SELECTORS.mfaForm)
      .hasText(
        'Sign in to Vault Multi-factor authentication is enabled for your account. Enter your authentication code to log in. TOTP passcode Verify Cancel'
      );
    await click(GENERAL.cancelButton);
    assert.dom(AUTH_FORM.form).exists('clicking back returns to auth form');
    assert.dom(AUTH_FORM.selectMethod).hasValue(this.authType, 'preserves method type on back');
    for (const field of loginKeys) {
      assert.dom(GENERAL.inputByAttr(field)).hasValue('', `${field} input clears on back`);
    }
  });

  test('it displays mfa requirement for custom paths', async function (assert) {
    this.path = `${this.authType}-custom`;
    const loginKeys = Object.keys(this.loginData);
    assert.expect(3 + loginKeys.length);
    this.stubRequests();
    await this.renderComponent();

    // Fill in login form
    await fillIn(AUTH_FORM.selectMethod, this.authType);
    // Toggle more options to input a custom mount path
    await fillInLoginFields({ ...this.loginData, path: this.path }, { toggleOptions: true });

    if (this.authType === 'oidc') {
      // fires "message" event which methods that rely on popup windows wait for
      triggerMessageEvent(this.path);
    }

    await click(GENERAL.submitButton);
    await waitFor(MFA_SELECTORS.mfaForm);
    assert
      .dom(MFA_SELECTORS.mfaForm)
      .hasText(
        'Sign in to Vault Multi-factor authentication is enabled for your account. Enter your authentication code to log in. TOTP passcode Verify Cancel'
      );
    await click(GENERAL.cancelButton);
    assert.dom(AUTH_FORM.form).exists('clicking back returns to auth form');
    assert.dom(AUTH_FORM.selectMethod).hasValue(this.authType, 'preserves method type on back');
    for (const field of loginKeys) {
      assert.dom(GENERAL.inputByAttr(field)).hasValue('', `${field} input clears on back`);
    }
  });

  test('it submits mfa requirement for default paths', async function (assert) {
    assert.expect(2);
    this.stubRequests();
    await this.renderComponent();

    const expectedOtp = '12345';
    this.server.post('/sys/mfa/validate', async (_, req) => {
      const [actualOtp] = JSON.parse(req.requestBody).mfa_payload[constraintId];
      assert.true(true, 'it makes request to mfa validate endpoint');
      assert.strictEqual(actualOtp, expectedOtp, 'payload contains otp');
    });

    // Fill in login form
    await fillIn(AUTH_FORM.selectMethod, this.authType);
    await fillInLoginFields(this.loginData);

    if (this.authType === 'oidc') {
      // fires "message" event which methods that rely on popup windows wait for
      triggerMessageEvent(this.path);
    }

    await click(GENERAL.submitButton);
    await waitFor(MFA_SELECTORS.mfaForm);
    await fillIn(MFA_SELECTORS.passcode(0), expectedOtp);
    await click(GENERAL.button('Verify'));
  });

  test('it submits mfa requirement for custom paths', async function (assert) {
    assert.expect(2);
    this.path = `${this.authType}-custom`;
    this.stubRequests();
    await this.renderComponent();

    const expectedOtp = '12345';
    this.server.post('/sys/mfa/validate', async (_, req) => {
      const [actualOtp] = JSON.parse(req.requestBody).mfa_payload[constraintId];
      assert.true(true, 'it makes request to mfa validate endpoint');
      assert.strictEqual(actualOtp, expectedOtp, 'payload contains otp');
    });

    // Fill in login form
    await fillIn(AUTH_FORM.selectMethod, this.authType);
    // Toggle more options to input a custom mount path
    await fillInLoginFields({ ...this.loginData, path: this.path }, { toggleOptions: true });

    if (this.authType === 'oidc') {
      triggerMessageEvent(this.path);
    }

    await click(GENERAL.submitButton);
    await waitFor(MFA_SELECTORS.mfaForm);
    await fillIn(MFA_SELECTORS.passcode(0), expectedOtp);
    await click(GENERAL.button('Verify'));
  });
};

module('Integration | Component | auth | page | mfa', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    setupTestContext(this);
    // additional setup for oidc-jwt component
    this.routerStub = sinon.stub(this.owner.lookup('service:router'), 'urlFor').returns('123-example.com');
  });

  hooks.afterEach(function () {
    this.routerStub.restore();
  });

  module('github', function (hooks) {
    hooks.beforeEach(async function () {
      this.authType = 'github';
      this.loginData = { token: 'mysupersecuretoken' };
      this.path = this.authType;
      this.stubRequests = () => {
        this.server.post(`/auth/${this.path}/login`, () => setupTotpMfaResponse(this.path));
      };
    });

    mfaTests(test);
  });

  module('jwt', function (hooks) {
    hooks.beforeEach(async function () {
      this.authType = 'jwt';
      this.loginData = { role: 'some-dev', jwt: 'jwttoken' };
      this.path = this.authType;

      this.stubRequests = () => {
        this.server.post('/auth/:path/oidc/auth_url', () =>
          overrideResponse(400, { errors: [ERROR_JWT_LOGIN] })
        );
        this.server.post(`/auth/${this.path}/login`, () => setupTotpMfaResponse(this.path));
      };
    });

    mfaTests(test);
  });

  module('oidc', function (hooks) {
    hooks.beforeEach(async function () {
      this.authType = 'oidc';
      this.loginData = { role: 'some-dev' };
      this.path = this.authType;
      // Requests are stubbed in the order they are hit
      this.stubRequests = () => {
        this.server.post(`/auth/${this.path}/oidc/auth_url`, () => {
          return { data: { auth_url: 'http://dev-foo-bar.com' } };
        });
        this.server.get(`/auth/${this.path}/oidc/callback`, () => setupTotpMfaResponse(this.path));
      };

      this.windowStub = windowStub();
    });

    hooks.afterEach(function () {
      this.windowStub.restore();
    });

    mfaTests(test);
  });

  module('username and password methods', function (hooks) {
    hooks.beforeEach(async function () {
      this.loginData = { username: 'matilda', password: 'password' };
      this.stubRequests = () => {
        this.server.post(`/auth/${this.path}/login/matilda`, () => setupTotpMfaResponse(this.path));
      };
    });

    module('ldap', function (hooks) {
      hooks.beforeEach(async function () {
        this.authType = 'ldap';
        this.path = this.authType;
      });

      mfaTests(test);
    });

    module('okta', function (hooks) {
      hooks.beforeEach(async function () {
        this.authType = 'okta';
        this.path = this.authType;
      });

      mfaTests(test);
    });

    module('radius', function (hooks) {
      hooks.beforeEach(async function () {
        this.authType = 'radius';
        this.path = this.authType;
      });

      mfaTests(test);
    });

    module('userpass', function (hooks) {
      hooks.beforeEach(async function () {
        this.authType = 'userpass';
        this.path = this.authType;
      });

      mfaTests(test);
    });
  });

  // ENTERPRISE METHODS
  module('saml', function (hooks) {
    hooks.beforeEach(async function () {
      this.version.type = 'enterprise';
      this.authType = 'saml';
      this.path = this.authType;
      this.loginData = { role: 'some-dev' };
      // Requests are stubbed in the order they are hit
      this.stubRequests = () => {
        this.server.post(`/auth/${this.path}/sso_service_url`, () => ({
          data: {
            sso_service_url: 'test/fake/sso/route',
            token_poll_id: '1234',
          },
        }));
        this.server.post(`/auth/${this.path}/token`, () => setupTotpMfaResponse(this.authType));
      };
      this.windowStub = windowStub();
    });

    hooks.afterEach(function () {
      this.windowStub.restore();
    });

    mfaTests(test);
  });
});
