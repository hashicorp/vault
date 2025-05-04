/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { currentURL, visit, currentRouteName, click, waitFor } from '@ember/test-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { login, logout } from 'vault/tests/helpers/auth/auth-helpers';
import { runCmd } from 'vault/tests/helpers/commands';
import { AUTH_FORM } from 'vault/tests/helpers/auth/auth-form-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';

module(`Acceptance | wrapped_token query param functionality`, function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    await login();
    // create wrapped token
    const token = await runCmd(`write -field=token auth/token/create policies=default -wrap-ttl=3m`);
    await logout();
    this.token = token;
  });

  test('it authenticates you if the query param is present', async function (assert) {
    await visit(`/vault/auth?wrapped_token=${this.token}`);
    assert.strictEqual(
      currentURL(),
      '/vault/dashboard',
      'authenticates and redirects to home (dashboard page)'
    );
  });

  test('it authenticates when used with the with=token query param', async function (assert) {
    await visit(`/vault/auth?wrapped_token=${this.token}&with=token`);
    assert.strictEqual(
      currentURL(),
      '/vault/dashboard',
      'authenticates and redirects to home (dashboard page)'
    );
  });

  test('it should authenticate when hitting logout url with wrapped_token when logged out', async function (assert) {
    await login();
    assert.strictEqual(currentRouteName(), 'vault.cluster.dashboard');
    await visit(`/vault/logout?wrapped_token=${this.token}`);
    assert.strictEqual(
      currentURL(),
      '/vault/dashboard',
      'authenticates and redirects to home (dashboard page)'
    );
  });

  test('it shows error if unwrap fails and goes back to login form', async function (assert) {
    await visit(`/vault/auth?wrapped_token=54321`);
    assert
      .dom(GENERAL.pageError.error)
      .hasText(
        'Authentication error Token unwrap failed Error: wrapping token is not valid or does not exist Go back'
      );
    assert.dom(AUTH_FORM.form).doesNotExist();
    await click(`${GENERAL.pageError.error} button`);
    await waitFor(AUTH_FORM.form);
    const url = currentURL();
    assert.false(url.includes('wrapped_token='), `url does not include wrapped_token param: ${url}`);
    assert.strictEqual(currentRouteName(), 'vault.cluster.auth', 'it navigates back to auth route');
    assert.dom(AUTH_FORM.form).exists('it navigates back to login form');
  });

  test('it makes request to authentication service with expected args', async function (assert) {
    // stub response so we know what the client_token value will be
    this.server.post('/sys/wrapping/unwrap', () => {
      return { auth: { client_token: '12345' } };
    });
    const authSpy = sinon.spy(this.owner.lookup('service:auth'), 'authenticate');
    await visit(`/vault/logout?wrapped_token=${this.token}`);
    const [actual] = authSpy.lastCall.args;
    assert.propEqual(
      actual,
      {
        backend: 'token',
        clusterId: '1',
        data: { token: '12345' },
        selectedAuth: 'token',
      },
      `it calls auth service authenticate method with correct args: ${JSON.stringify(actual)} `
    );
  });
});
