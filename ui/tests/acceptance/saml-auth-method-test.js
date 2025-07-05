/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { click, fillIn, find, visit, waitUntil } from '@ember/test-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { Response } from 'miragejs';
import { DELAY_IN_MS, windowStub } from 'vault/tests/helpers/oidc-window-stub';
import { AUTH_FORM } from 'vault/tests/helpers/auth/auth-form-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

import { logout } from 'vault/tests/helpers/auth/auth-helpers';

module('Acceptance | enterprise saml auth method', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.openStub = windowStub();
    this.server.put('/auth/saml/sso_service_url', () => ({
      data: {
        sso_service_url: 'test/fake/sso/route', // we aren't actually opening a popup so use a fake url
        token_poll_id: '1234',
      },
    }));
    this.server.put('/auth/saml/token', () => ({
      auth: { client_token: 'root' },
    }));
    // ensure clean state
    await logout(); // clears local storage
  });

  hooks.afterEach(function () {
    this.openStub.restore();
  });

  test('it should login with saml when selected from auth methods dropdown', async function (assert) {
    assert.expect(1);

    this.server.get('/auth/token/lookup-self', (schema, req) => {
      assert.true(true, 'request made to auth/token/lookup-self after saml callback');
      return req.passthrough();
    });
    // select from dropdown or click auth path tab
    await waitUntil(() => find(AUTH_FORM.selectMethod), { timeout: DELAY_IN_MS });
    await fillIn(AUTH_FORM.selectMethod, 'saml');
    await click(GENERAL.submitButton);
  });

  test('it should login with saml from listed auth mount tab', async function (assert) {
    assert.expect(4);
    this.server.get('/sys/internal/ui/mounts', () => ({
      data: {
        auth: {
          'test-path/': { description: '', options: {}, type: 'saml' },
        },
      },
    }));
    this.server.put('/auth/test-path/sso_service_url', () => {
      assert.ok(true, 'role request made to correct non-standard mount path');
      return {
        data: {
          sso_service_url: 'test/fake/sso/route',
          token_poll_id: '1234',
        },
      };
    });
    this.server.put('/auth/test-path/token', () => {
      assert.ok(true, 'login request made to correct non-standard mount path');
      return {
        auth: { client_token: 'root' },
      };
    });
    this.server.get('/auth/token/lookup-self', (schema, req) => {
      assert.ok(true, 'request made to auth/token/lookup-self after oidc callback');
      assert.deepEqual(
        req.requestHeaders,
        { 'X-Vault-Token': 'root' },
        'calls lookup-self with returned client token after login'
      );
      return req.passthrough();
    });

    await logout(); // clear local storage and refresh route so sys/internal/ui/mounts is reliably called
    // click auth path tab
    await waitUntil(() => find(AUTH_FORM.tabBtn('saml')), { timeout: DELAY_IN_MS });
    await click(GENERAL.submitButton);
  });

  test('it should render API errors from sso_service_url', async function (assert) {
    assert.expect(1);
    this.server.put('/auth/saml/sso_service_url', () => {
      return new Response(
        400,
        { 'Content-Type': 'application/json' },
        JSON.stringify({ errors: ["missing required 'role' parameter"] })
      );
    });

    // select saml auth type
    await waitUntil(() => find(AUTH_FORM.selectMethod), { timeout: DELAY_IN_MS });
    await fillIn(AUTH_FORM.selectMethod, 'saml');
    await click(GENERAL.submitButton);
    assert
      .dom('[data-test-message-error-description]')
      .hasText("Authentication failed: missing required 'role' parameter", 'shows API error from role fetch');
  });

  test('it should render API errors from saml token login url', async function (assert) {
    assert.expect(1);
    this.server.put('/auth/saml/token', () => {
      return new Response(
        400,
        { 'Content-Type': 'application/json' },
        JSON.stringify({ errors: ['something went wrong'] })
      );
    });

    // select saml auth type
    await waitUntil(() => find(AUTH_FORM.selectMethod), { timeout: DELAY_IN_MS });
    await fillIn(AUTH_FORM.selectMethod, 'saml');
    await click(GENERAL.submitButton);
    assert
      .dom('[data-test-message-error-description]')
      .hasText('Authentication failed: something went wrong', 'shows API error from login attempt');
  });

  test('it should populate saml auth method on logout', async function (assert) {
    await visit('/vault/logout');
    // select from dropdown
    await waitUntil(() => find(AUTH_FORM.selectMethod), { timeout: DELAY_IN_MS });
    await fillIn(AUTH_FORM.selectMethod, 'saml');
    await click(GENERAL.submitButton);
    await waitUntil(() => find(GENERAL.button('user-menu-trigger')), { timeout: DELAY_IN_MS });
    await click(GENERAL.button('user-menu-trigger'));
    await click('#logout');
    assert.dom(AUTH_FORM.selectMethod).hasValue('saml', 'Previous auth method selected on logout');
  });
});
