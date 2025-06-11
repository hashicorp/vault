/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, fillIn } from '@ember/test-helpers';
import { AUTH_FORM } from 'vault/tests/helpers/auth/auth-form-selectors';
import { fillInLoginFields } from 'vault/tests/helpers/auth/auth-helpers';
import { EXPECTED_TOKEN_DATA } from 'vault/tests/helpers/auth/response-stubs';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { setupTotpMfaResponse } from 'vault/tests/helpers/mfa/mfa-helpers';

export const methodAuthenticationTests = (test) => {
  test('it sets token data on login for default path', async function (assert) {
    assert.expect(5);
    // Setup
    this.stubRequests();
    // Render and log in
    await this.renderComponent();
    await fillIn(AUTH_FORM.selectMethod, this.authType);
    await fillInLoginFields(this.loginData);
    await click(AUTH_FORM.login);

    const [tokenName, persistedTokenData] = this.setTokenDataSpy.lastCall.args;

    const expectedData = {
      ...EXPECTED_TOKEN_DATA[this.authType],
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
    await click(AUTH_FORM.login);

    const [actual] = this.onAuthSuccess.lastCall.args;
    const expected = {
      namespace: '',
      token: this.tokenName,
      isRoot: false,
    };
    assert.propEqual(actual, expected, `onAuthSuccess called with: ${JSON.stringify(actual)}`);
  });

  test('it preselects auth type from canceled mfa', async function (assert) {
    assert.expect(1);
    this.response = setupTotpMfaResponse(this.path);
    this.stubRequests();

    await this.renderComponent();
    await fillIn(AUTH_FORM.selectMethod, this.authType);
    await fillInLoginFields(this.loginData);

    await click(AUTH_FORM.login);
    await click(GENERAL.backButton);
    assert.dom(AUTH_FORM.selectMethod).hasValue(this.authType, `${this.authType} is selected in dropdown`);
  });
};
// test('it sets token data on login for default path', async function (assert) {
//   assert.expect(5);
//   const { loginData, stubRequests } = options;
//   stubRequests({ server: this.server, path: authType, response: RESPONSE_STUBS[authType] });

//   await this.renderComponent();
//   await fillIn(AUTH_FORM.selectMethod, authType);
//   await fillInLoginFields(loginData);

//   // For OIDC
//   if (options?.hasPopupWindow) {
//     // fires "message" event which methods that rely on popup windows wait for
//     setTimeout(() => {
//       window.postMessage(callbackData({ path: authType }), window.origin);
//     }, DELAY_IN_MS);
//   }

//   await click(AUTH_FORM.login);

//   const [tokenName, actualTokenData] = this.setTokenDataSpy.lastCall.args;

//   const expectedData = {
//     ...EXPECTED_TOKEN_DATA[authType],
//     // there are other tests that confirm this calculation happens as expected, just copy value from spy
//     tokenExpirationEpoch: actualTokenData.tokenExpirationEpoch,
//   };
//   const expectedName = ['oidc', 'saml'].includes(authType) ? 'token' : authType;
//   assert.strictEqual(tokenName, `vault-${expectedName}☃1`, 'token data has expected token name');
//   assert.propEqual(actualTokenData, expectedData, `setTokenData is called with expected data`);
//   // propEqual failures are challenging to parse in CI so pulling out a couple of important attrs
//   const {
//     token: expectedToken,
//     displayName: expectedDisplayName,
//     entity_id: expectedEntityId,
//   } = expectedData;
//   assert.strictEqual(actualTokenData.token, expectedToken, 'persisted token data has expected token');
//   assert.strictEqual(
//     actualTokenData.displayName,
//     expectedDisplayName,
//     'persisted token data has expected display name'
//   );
//   assert.strictEqual(
//     actualTokenData.entity_id,
//     expectedEntityId,
//     'persisted token data has expected entity_id'
//   );
// });
