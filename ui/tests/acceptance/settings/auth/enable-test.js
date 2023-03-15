/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { currentRouteName, settled } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import page from 'vault/tests/pages/settings/auth/enable';
import listPage from 'vault/tests/pages/access/methods';
import authPage from 'vault/tests/pages/auth';

module('Acceptance | settings/auth/enable', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    this.timestamp = new Date().getTime();
    return authPage.login();
  });

  test('it mounts and redirects', async function (assert) {
    // always force the new mount to the top of the list
    const path = `approle-${this.timestamp}`;
    const type = 'approle';
    await page.visit();
    assert.strictEqual(currentRouteName(), 'vault.cluster.settings.auth.enable');
    await page.enable(type, path);
    await settled();
    await assert.strictEqual(
      page.flash.latestMessage,
      `Successfully mounted the ${type} auth method at ${path}.`,
      'success flash shows'
    );
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.settings.auth.configure.section',
      'redirects to the auth config page'
    );

    await listPage.visit();
    assert.ok(listPage.findLinkById(path), 'mount is present in the list');
  });
});
