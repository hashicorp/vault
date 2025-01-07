/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import sinon from 'sinon';
import { click, fillIn, find, waitUntil } from '@ember/test-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { Response } from 'miragejs';
import authPage from 'vault/tests/pages/auth';
import { fakeWindow } from 'vault/tests/helpers/oidc-window-stub';
import { setupTotpMfaResponse } from 'vault/tests/helpers/auth/mfa-helpers';

module('Acceptance | enterprise saml auth method', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.openStub = sinon.stub(window, 'open').callsFake(() => fakeWindow.create());
    this.server.put('/auth/saml/sso_service_url', () => ({
      data: {
        sso_service_url: 'http://sso-url.hashicorp.com/service',
        token_poll_id: '1234',
      },
    }));
    this.server.put('/auth/saml/token', () => ({
      auth: { client_token: 'root' },
    }));
    // ensure clean state
    localStorage.removeItem('selectedAuth');
    authPage.logout();
  });

  hooks.afterEach(function () {
    this.openStub.restore();
  });

  test('it should login with saml when selected from auth methods dropdown', async function (assert) {
    assert.expect(1);

    this.server.get('/auth/token/lookup-self', (schema, req) => {
      assert.ok(true, 'request made to auth/token/lookup-self after saml callback');
      req.passthrough();
    });
    // select from dropdown or click auth path tab
    await waitUntil(() => find('[data-test-select="auth-method"]'));
    await fillIn('[data-test-select="auth-method"]', 'saml');
    await click('[data-test-auth-submit]');
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
          sso_service_url: 'http://sso-url.hashicorp.com/service',
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
      req.passthrough();
    });

    // click auth path tab
    await waitUntil(() => find('[data-test-auth-method="test-path"]'));
    await click('[data-test-auth-method="test-path"]');
    await click('[data-test-auth-submit]');
  });

  test('it should render API errors from both endpoints', async function (assert) {
    assert.expect(3);

    this.server.put('/auth/saml/sso_service_url', (schema, { requestBody }) => {
      const { role } = JSON.parse(requestBody);
      if (!role) {
        return new Response(
          400,
          { 'Content-Type': 'application/json' },
          JSON.stringify({ errors: ["missing required 'role' parameter"] })
        );
      }
      return {
        data: {
          sso_service_url: 'http://sso-url.hashicorp.com/service',
          token_poll_id: '1234',
        },
      };
    });
    this.server.put('/auth/saml/token', (schema, { requestHeaders }) => {
      if (requestHeaders['X-Vault-Namespace']) {
        return new Response(
          400,
          { 'Content-Type': 'application/json' },
          JSON.stringify({ errors: ['something went wrong'] })
        );
      }
      return {
        auth: { client_token: 'root' },
      };
    });
    this.server.get('/auth/token/lookup-self', (schema, req) => {
      assert.ok(true, 'request made to auth/token/lookup-self after saml callback');
      req.passthrough();
    });

    // select saml auth type
    await waitUntil(() => find('[data-test-select="auth-method"]'));
    await fillIn('[data-test-select="auth-method"]', 'saml');
    await fillIn('[data-test-auth-form-ns-input]', 'some-ns');
    await click('[data-test-auth-submit]');
    assert
      .dom('[data-test-message-error-description]')
      .hasText("missing required 'role' parameter", 'shows API error from role fetch');

    await fillIn('[data-test-role]', 'my-role');
    await click('[data-test-auth-submit]');
    assert
      .dom('[data-test-message-error-description]')
      .hasText('something went wrong', 'shows API error from login attempt');

    await fillIn('[data-test-auth-form-ns-input]', '');
    await click('[data-test-auth-submit]');
  });

  test('it should populate saml auth method on logout', async function (assert) {
    authPage.logout();
    // select from dropdown
    await waitUntil(() => find('[data-test-select="auth-method"]'));
    await fillIn('[data-test-select="auth-method"]', 'saml');
    await click('[data-test-auth-submit]');
    await waitUntil(() => find('[data-test-user-menu-trigger]'));
    await click('[data-test-user-menu-trigger]');
    await click('#logout');
    assert
      .dom('[data-test-select="auth-method"]')
      .hasValue('saml', 'Previous auth method selected on logout');
  });

  test('it prompts mfa if configured', async function (assert) {
    assert.expect(1);
    this.server.put('/auth/saml/token', () => setupTotpMfaResponse('saml'));

    await waitUntil(() => find('[data-test-select="auth-method"]'));
    await fillIn('[data-test-select="auth-method"]', 'saml');
    await click('[data-test-auth-submit]');
    await waitUntil(() => find('[data-test-mfa-form]'));
    assert.dom('[data-test-mfa-form]').exists('it renders TOTP MFA form');
  });
});
