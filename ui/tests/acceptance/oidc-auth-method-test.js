/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { click, fillIn, waitFor } from '@ember/test-helpers';
import { logout } from 'vault/tests/helpers/auth/auth-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import {
  callbackData,
  DELAY_IN_MS,
  triggerMessageEvent,
  windowStub,
} from 'vault/tests/helpers/oidc-window-stub';
import { AUTH_FORM } from 'vault/tests/helpers/auth/auth-form-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { ERROR_MISSING_PARAMS, ERROR_POPUP_FAILED, ERROR_WINDOW_CLOSED } from 'vault/utils/auth-form-helpers';
import { getErrorResponse } from 'vault/tests/helpers/api/error-response';
import sinon from 'sinon';
import { RESPONSE_STUBS } from '../helpers/auth/response-stubs';

module('Acceptance | oidc auth method', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.openStub = windowStub();

    this.setupMocks = () => {
      this.server.post('/auth/oidc/oidc/auth_url', () => ({
        data: { auth_url: 'http://example.com' },
      }));
      this.server.get(`/auth/oidc/oidc/callback`, () => RESPONSE_STUBS.oidc['oidc/callback']);
      this.server.get(`/auth/token/lookup-self`, () => RESPONSE_STUBS.oidc['lookup-self']);
    };
  });

  hooks.afterEach(async function () {
    this.openStub.restore();
    // ensure clean state
    // Cannot use logout() here because it will hit the internal mount request before the mocks can interrupt it
    localStorage.clear();
  });

  // coverage for bug where token was selected as auth method for oidc and jwt
  test('it should populate oidc auth method on logout', async function (assert) {
    this.setupMocks();
    await logout();
    await fillIn(AUTH_FORM.selectMethod, 'oidc');

    triggerMessageEvent('oidc');

    await click(GENERAL.submitButton);
    await waitFor('[data-test-dashboard-card-header="Vault version"]');
    assert
      .dom('[data-test-dashboard-card-header="Vault version"]')
      .exists('Render the dashboard landing page.');

    await logout();
    assert.dom(AUTH_FORM.selectMethod).hasValue('oidc', 'Previous auth method selected on logout');
  });

  // test case for https://github.com/hashicorp/vault/issues/12436
  test('it should ignore messages sent from outside the app while waiting for oidc callback', async function (assert) {
    assert.expect(3); // one for both message events (2) and one for callback request
    this.setupMocks();
    this.server.get('/auth/foo/oidc/callback', () => {
      // third assertion
      assert.true(true, 'request is made to callback url');
      return { auth: { client_token: 'root' } };
    });

    let count = 0;
    const assertEvent = (event) => {
      count++;
      // we have to use the same event method, but need to update what it checks for depending on when it's called
      const source = count === 1 ? 'miscellaneous-source' : 'oidc-callback';
      assert.strictEqual(event.data.source, source, `message event fires with source: ${event.data.source}`);
    };
    window.addEventListener('message', assertEvent);
    await logout();
    await fillIn(AUTH_FORM.selectMethod, 'oidc');

    setTimeout(() => {
      // first assertion
      window.postMessage(callbackData({ source: 'miscellaneous-source' }), window.origin);
      // second assertion
      window.postMessage(callbackData({ source: 'oidc-callback' }), window.origin);
    }, DELAY_IN_MS);

    await click(GENERAL.submitButton);
    // cleanup
    window.removeEventListener('message', assertEvent);
  });

  test('it shows error when message posted with state key, wrong params', async function (assert) {
    this.setupMocks();
    await logout();
    await fillIn(AUTH_FORM.selectMethod, 'oidc');
    setTimeout(() => {
      // callback params are missing "code"
      window.postMessage({ source: 'oidc-callback', state: 'state', foo: 'bar' }, window.origin);
    }, DELAY_IN_MS);
    await click(GENERAL.submitButton);
    assert
      .dom(GENERAL.messageError)
      .hasText(`Error Authentication failed: ${ERROR_MISSING_PARAMS}`, 'displays error when missing params');
  });

  test('it shows error when popup is prematurely closed ', async function (assert) {
    windowStub({ stub: this.openStub, popup: { closed: true, close: () => {} } });

    this.setupMocks();
    await logout();
    await fillIn(AUTH_FORM.selectMethod, 'oidc');
    await click(GENERAL.submitButton);
    assert.dom(GENERAL.messageError).hasText(`Error Authentication failed: ${ERROR_WINDOW_CLOSED}`);
  });

  test('it renders error when window fails to open', async function (assert) {
    this.openStub.returns(null);
    this.setupMocks();
    await logout();
    await fillIn(AUTH_FORM.selectMethod, 'oidc');
    await click(GENERAL.submitButton);
    assert
      .dom(GENERAL.messageError)
      .hasText(`Error Authentication failed: Failed to open OIDC popup window. ${ERROR_POPUP_FAILED}`);
  });

  test('it renders api errors if oidc callback request fails', async function (assert) {
    await logout();
    this.server.post('/auth/oidc/oidc/auth_url', () => ({
      data: { auth_url: 'http://example.com' },
    }));
    const api = this.owner.lookup('service:api');
    const oidcCallbackStub = sinon.stub(api.auth, 'jwtOidcCallback');
    oidcCallbackStub.rejects(getErrorResponse({ errors: ['something went terribly wrong!'] }, 500));
    await fillIn(AUTH_FORM.selectMethod, 'oidc');
    triggerMessageEvent('oidc');
    await click(GENERAL.submitButton);
    assert.dom(GENERAL.messageError).hasText('Error Authentication failed: something went terribly wrong!');
    oidcCallbackStub.restore();
  });
});
