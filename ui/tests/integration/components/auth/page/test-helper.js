/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, fillIn, render, waitFor, waitUntil } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { AUTH_FORM } from 'vault/tests/helpers/auth/auth-form-selectors';
import { fillInLoginFields } from 'vault/tests/helpers/auth/auth-helpers';
import { TOKEN_DATA } from 'vault/tests/helpers/auth/response-stubs';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { constraintId } from 'vault/tests/helpers/mfa/mfa-helpers';
import { MFA_SELECTORS } from 'vault/tests/helpers/mfa/mfa-selectors';
import { triggerMessageEvent } from 'vault/tests/helpers/oidc-window-stub';

export const setupTestContext = (context) => {
  context.version = context.owner.lookup('service:version');
  context.cluster = { id: '1' };
  context.directLinkData = null;
  context.loginSettings = null;
  context.namespaceQueryParam = '';
  context.oidcProviderQueryParam = '';
  context.onAuthSuccess = sinon.spy();
  context.onNamespaceUpdate = sinon.spy();
  context.visibleAuthMounts = false;

  context.renderComponent = () => {
    return render(hbs`<Auth::Page
  @cluster={{this.cluster}}
  @directLinkData={{this.directLinkData}}
  @loginSettings={{this.loginSettings}}
  @namespaceQueryParam={{this.namespaceQueryParam}}
  @oidcProviderQueryParam={{this.oidcProviderQueryParam}}
  @onAuthSuccess={{this.onAuthSuccess}}
  @onNamespaceUpdate={{this.onNamespaceUpdate}}
  @visibleAuthMounts={{this.visibleAuthMounts}}
/>`);
  };
};

export const methodAuthenticationTests = (test) => {
  test('it sets token data on login for default path', async function (assert) {
    assert.expect(5);
    // Setup
    this.stubRequests();
    // Render and log in
    await this.renderComponent();
    await fillIn(AUTH_FORM.selectMethod, this.authType);
    await fillInLoginFields(this.loginData);
    if (this.authType === 'oidc') {
      triggerMessageEvent(this.path);
    }
    await click(GENERAL.submitButton);
    await waitUntil(() => this.setTokenDataSpy.calledOnce);
    const [tokenName, persistedTokenData] = this.setTokenDataSpy.lastCall.args;

    const expectedData = {
      ...TOKEN_DATA[this.authType],
      // there are other tests that confirm this calculation happens as expected, just copy value from spy
      tokenExpirationEpoch: persistedTokenData.tokenExpirationEpoch,
    };

    assert.strictEqual(tokenName, this.tokenName, 'setTokenData is called with expected token name');
    assert.propEqual(persistedTokenData, expectedData, 'setTokenData is called with expected data');

    // propEqual failures are challenging to parse in CI so pulling out a couple of important attrs
    const { token, displayName, entity_id } = expectedData;
    assert.strictEqual(persistedTokenData.token, token, 'setTokenData has expected token');
    assert.strictEqual(persistedTokenData.displayName, displayName, 'setTokenData has expected display name');
    assert.strictEqual(persistedTokenData.entity_id, entity_id, 'setTokenData has expected entity_id');
  });

  test('it calls onAuthSuccess on submit for custom path', async function (assert) {
    assert.expect(1);
    // Setup
    this.path = `${this.authType}-custom`;
    this.loginData = { ...this.loginData, path: this.path };
    this.stubRequests();
    // Render and log in
    await this.renderComponent();
    await fillIn(AUTH_FORM.selectMethod, this.authType);
    // toggle mount path input to specify custom path
    await fillInLoginFields(this.loginData, { toggleOptions: true });
    if (this.authType === 'oidc') {
      triggerMessageEvent(this.path);
    }
    await click(GENERAL.submitButton);

    await waitUntil(() => this.onAuthSuccess.calledOnce);
    const [actual] = this.onAuthSuccess.lastCall.args;
    const expected = { namespace: '', token: this.tokenName, isRoot: false };
    assert.propEqual(actual, expected, `onAuthSuccess called with: ${JSON.stringify(actual)}`);
  });
};

export const mfaTests = (test) => {
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
        'Back Multi-factor authentication is enabled for your account. Enter your authentication code to log in. TOTP passcode Verify'
      );
    await click(GENERAL.backButton);
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
        'Back Multi-factor authentication is enabled for your account. Enter your authentication code to log in. TOTP passcode Verify'
      );
    await click(GENERAL.backButton);
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
    await click(MFA_SELECTORS.validate);
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
    await click(MFA_SELECTORS.validate);
  });
};
