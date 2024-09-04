/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { currentURL, visit, currentRouteName } from '@ember/test-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { login, logout } from 'vault/tests/helpers/auth/auth-helpers';
import { runCmd } from 'vault/tests/helpers/commands';

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
});
