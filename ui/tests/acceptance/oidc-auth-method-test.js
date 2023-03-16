/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { click, fillIn, find, waitUntil } from '@ember/test-helpers';
import authPage from 'vault/tests/pages/auth';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { fakeWindow, buildMessage } from '../helpers/oidc-window-stub';
import sinon from 'sinon';
import { later, _cancelTimers as cancelTimers } from '@ember/runloop';
import timestamp from 'core/utils/timestamp';

module('Acceptance | oidc auth method', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.openStub = sinon.stub(window, 'open').callsFake(() => fakeWindow.create());
    // OIDC test fails when using fake timestamps
    this.dateStub = sinon.stub(timestamp, 'now').callsFake(() => new Date());
    this.server.post('/auth/oidc/oidc/auth_url', () => ({
      data: { auth_url: 'http://example.com' },
    }));
    this.server.get('/auth/foo/oidc/callback', () => ({
      auth: { client_token: 'root' },
    }));
    // ensure clean state
    localStorage.removeItem('selectedAuth');
  });

  hooks.afterEach(function () {
    this.openStub.restore();
    this.dateStub.restore();
  });

  test('it should login with oidc when selected from auth methods dropdown', async function (assert) {
    assert.expect(1);

    this.server.get('/auth/token/lookup-self', (schema, req) => {
      assert.ok(true, 'request made to auth/token/lookup-self after oidc callback');
      req.passthrough();
    });
    authPage.logout();
    // select from dropdown or click auth path tab
    await waitUntil(() => find('[data-test-select="auth-method"]'));
    await fillIn('[data-test-select="auth-method"]', 'oidc');
    later(() => {
      window.postMessage(buildMessage().data, window.origin);
      cancelTimers();
    }, 50);
    await click('[data-test-auth-submit]');
  });

  test('it should login with oidc from listed auth mount tab', async function (assert) {
    assert.expect(3);

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
    // there was a bug that would result in the /auth/:path/login endpoint hit with an empty payload rather than lookup-self
    // ensure that the correct endpoint is hit after the oidc callback
    this.server.get('/auth/token/lookup-self', (schema, req) => {
      assert.ok(true, 'request made to auth/token/lookup-self after oidc callback');
      req.passthrough();
    });

    authPage.logout();
    // select from dropdown or click auth path tab
    await waitUntil(() => find('[data-test-auth-method-link="oidc"]'));
    await click('[data-test-auth-method-link="oidc"]');
    later(() => {
      window.postMessage(buildMessage().data, window.origin);
      cancelTimers();
    }, 50);
    await click('[data-test-auth-submit]');
  });

  // coverage for bug where token was selected as auth method for oidc and jwt
  test('it should populate oidc auth method on logout', async function (assert) {
    authPage.logout();
    // select from dropdown or click auth path tab
    await waitUntil(() => find('[data-test-select="auth-method"]'));
    await fillIn('[data-test-select="auth-method"]', 'oidc');
    later(() => {
      window.postMessage(buildMessage().data, window.origin);
      cancelTimers();
    }, 50);
    await click('[data-test-auth-submit]');
    await waitUntil(() => find('.nav-user-button button'), { timeout: 2000 });
    await click('.nav-user-button button');
    await click('#logout');
    assert
      .dom('[data-test-select="auth-method"]')
      .hasValue('oidc', 'Previous auth method selected on logout');
  });
});
