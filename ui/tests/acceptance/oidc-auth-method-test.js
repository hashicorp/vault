/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { click, fillIn, find, waitUntil } from '@ember/test-helpers';
import authPage from 'vault/tests/pages/auth';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { fakeWindow, buildMessage } from '../helpers/oidc-window-stub';
import sinon from 'sinon';
import { later, _cancelTimers as cancelTimers } from '@ember/runloop';
import { Response } from 'miragejs';
import { setupTotpMfaResponse } from 'vault/tests/helpers/auth/mfa-helpers';

module('Acceptance | oidc auth method', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.openStub = sinon.stub(window, 'open').callsFake(() => fakeWindow.create());

    this.setupMocks = (assert) => {
      this.server.post('/auth/oidc/oidc/auth_url', () => ({
        data: { auth_url: 'http://example.com' },
      }));
      // there was a bug that would result in the /auth/:path/login endpoint hit with an empty payload rather than lookup-self
      // ensure that the correct endpoint is hit after the oidc callback
      if (assert) {
        this.server.get('/auth/token/lookup-self', (schema, req) => {
          assert.ok(true, 'request made to auth/token/lookup-self after oidc callback');
          req.passthrough();
        });
      }
    };

    this.server.get('/auth/foo/oidc/callback', () => ({
      auth: { client_token: 'root' },
    }));

    // select method from dropdown or click auth path tab
    this.selectMethod = async (method, useLink) => {
      const methodSelector = useLink
        ? `[data-test-auth-method-link="${method}"]`
        : '[data-test-select="auth-method"]';
      await waitUntil(() => find(methodSelector));
      if (useLink) {
        await click(`[data-test-auth-method-link="${method}"]`);
      } else {
        await fillIn('[data-test-select="auth-method"]', method);
      }
    };

    // ensure clean state
    localStorage.removeItem('selectedAuth');
    authPage.logout();
  });

  hooks.afterEach(function () {
    this.openStub.restore();
  });

  test('it should login with oidc when selected from auth methods dropdown', async function (assert) {
    assert.expect(1);

    this.setupMocks(assert);

    await this.selectMethod('oidc');
    later(() => {
      window.postMessage(buildMessage().data, window.origin);
      cancelTimers();
    }, 100);

    await click('[data-test-auth-submit]');
  });

  test('it should login with oidc from listed auth mount tab', async function (assert) {
    assert.expect(3);

    this.setupMocks(assert);

    this.server.get('/sys/internal/ui/mounts', () => ({
      data: {
        auth: {
          'test-path/': { description: '', options: {}, type: 'oidc' },
        },
      },
    }));
    // this request is fired twice -- total assertion count should be 3 rather than 2
    // JLR TODO - auth-jwt: verify whether additional request is necessary, especially when glimmerizing component
    // look into whether didReceiveAttrs is necessary to trigger this request
    this.server.post('/auth/test-path/oidc/auth_url', () => {
      assert.ok(true, 'auth_url request made to correct non-standard mount path');
      return { data: { auth_url: 'http://example.com' } };
    });

    await this.selectMethod('oidc', true);
    later(() => {
      window.postMessage(buildMessage().data, window.origin);
      cancelTimers();
    }, 50);

    await click('[data-test-auth-submit]');
  });

  // coverage for bug where token was selected as auth method for oidc and jwt
  test('it should populate oidc auth method on logout', async function (assert) {
    this.setupMocks();
    await this.selectMethod('oidc');
    later(() => {
      window.postMessage(buildMessage().data, window.origin);
      cancelTimers();
    }, 50);

    await click('[data-test-auth-submit]');
    await waitUntil(() => find('[data-test-user-menu-trigger]'));
    await click('[data-test-user-menu-trigger]');
    await click('#logout');
    assert
      .dom('[data-test-select="auth-method"]')
      .hasValue('oidc', 'Previous auth method selected on logout');
  });

  test('it should fetch role when switching between oidc/jwt auth methods and changing the mount path', async function (assert) {
    let reqCount = 0;
    this.server.post('/auth/:method/oidc/auth_url', (schema, req) => {
      reqCount++;
      const errors =
        req.params.method === 'jwt' ? ['OIDC login is not configured for this mount'] : ['missing role'];
      return new Response(400, {}, { errors });
    });

    await this.selectMethod('oidc');
    assert.dom('[data-test-jwt]').doesNotExist('JWT Token input hidden for OIDC');
    await this.selectMethod('jwt');
    assert.dom('[data-test-jwt]').exists('JWT Token input renders for JWT configured method');
    await click('[data-test-auth-form-options-toggle]');
    await fillIn('[data-test-auth-form-mount-path]', 'foo');
    assert.strictEqual(reqCount, 3, 'Role is fetched when dependant values are changed');
  });

  test('it should display role fetch errors when signing in with OIDC', async function (assert) {
    this.server.post('/auth/:method/oidc/auth_url', (schema, req) => {
      const { role } = JSON.parse(req.requestBody);
      const status = role ? 403 : 400;
      const errors = role ? ['permission denied'] : ['missing role'];
      return new Response(status, {}, { errors });
    });

    await this.selectMethod('oidc');
    await click('[data-test-auth-submit]');
    assert.dom('[data-test-message-error-description]').hasText('Invalid role. Please try again.');

    await fillIn('[data-test-role]', 'test');
    await click('[data-test-auth-submit]');
    assert.dom('[data-test-message-error-description]').hasText('Error fetching role: permission denied');
  });

  test('it prompts mfa if configured', async function (assert) {
    assert.expect(1);

    this.setupMocks(assert);
    this.server.get('/auth/foo/oidc/callback', () => setupTotpMfaResponse('foo'));
    await this.selectMethod('oidc');
    later(() => {
      window.postMessage(buildMessage().data, window.origin);
      cancelTimers();
    }, 50);

    await click('[data-test-auth-submit]');
    await waitUntil(() => find('[data-test-mfa-form]'));
    assert.dom('[data-test-mfa-form]').exists('it renders TOTP MFA form');
  });
});
