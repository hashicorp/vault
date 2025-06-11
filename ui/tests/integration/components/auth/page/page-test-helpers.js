/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, fillIn, render, waitUntil } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { AUTH_FORM } from 'vault/tests/helpers/auth/auth-form-selectors';
import { fillInLoginFields } from 'vault/tests/helpers/auth/auth-helpers';
import { TOKEN_DATA } from 'vault/tests/helpers/auth/response-stubs';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
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
