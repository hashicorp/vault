/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { currentURL, find, visit, settled } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

import backendListPage from 'vault/tests/pages/secrets/backends';
import mountSecrets from 'vault/tests/pages/settings/mount-secret-backend';
import authPage from 'vault/tests/pages/auth';
import logout from 'vault/tests/pages/logout';

module('Acceptance | settings', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    this.uid = uuidv4();
    return authPage.login();
  });

  hooks.afterEach(function () {
    return logout.visit();
  });

  test('settings', async function (assert) {
    const type = 'consul';
    const path = `settings-path-${this.uid}`;

    // mount unsupported backend
    await visit('/vault/settings/mount-secret-backend');

    assert.strictEqual(currentURL(), '/vault/settings/mount-secret-backend');

    await mountSecrets.selectType(type);
    await mountSecrets
      .next()
      .path(path)
      .toggleOptions()
      .enableDefaultTtl()
      .defaultTTLUnit('s')
      .defaultTTLVal(100)
      .submit();
    await settled();
    assert.ok(
      find('[data-test-flash-message]').textContent.trim(),
      `Successfully mounted '${type}' at '${path}'!`
    );
    await settled();
    assert.strictEqual(currentURL(), `/vault/secrets`, 'redirects to secrets page');
    const row = backendListPage.rows.filterBy('path', path + '/')[0];
    await row.menu();
    await backendListPage.configLink();
    assert.strictEqual(currentURL(), `/vault/secrets/${path}/configuration`, 'navigates to the config page');
  });
});
