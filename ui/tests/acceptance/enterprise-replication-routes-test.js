/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { currentRouteName, currentURL, visit } from '@ember/test-helpers';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

// These run against the real Enterprise backend without mocks. The mocked variants live in
// replication-routes-test.js and run in both editions.
module('Acceptance | Enterprise | replication routes against the backend', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    return login();
  });

  test('the replication index renders a page header', async function (assert) {
    await visit('/vault/replication');

    assert.strictEqual(currentURL(), '/vault/replication');
    assert.strictEqual(currentRouteName(), 'vault.cluster.replication.index');
    assert.dom(GENERAL.hdsPageHeaderTitle).exists();
  });

  test('an unknown replication mode redirects to the replication index', async function (assert) {
    await visit('/vault/replication/not-a-mode');

    assert.strictEqual(currentURL(), '/vault/replication');
    assert.strictEqual(currentRouteName(), 'vault.cluster.replication.index');
  });

  test('the replication sidebar links navigate to the mode pages', async function (assert) {
    await visit('/vault/replication');

    assert.dom(GENERAL.navLink('Overview')).hasAttribute('href', /\/vault\/replication$/);
    assert.dom(GENERAL.navLink('Disaster recovery')).hasAttribute('href', /\/vault\/replication\/dr$/);
    assert
      .dom(GENERAL.navLink('Performance replication'))
      .hasAttribute('href', /\/vault\/replication\/performance$/);
  });
});
