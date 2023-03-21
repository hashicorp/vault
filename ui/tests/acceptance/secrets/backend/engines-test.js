/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { currentRouteName, settled } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import mountSecrets from 'vault/tests/pages/settings/mount-secret-backend';
import backendsPage from 'vault/tests/pages/secrets/backends';
import authPage from 'vault/tests/pages/auth';

module('Acceptance | engine/disable', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    return authPage.login();
  });

  test('disable engine', async function (assert) {
    // first mount an engine so we can disable it.
    const enginePath = `alicloud-${new Date().getTime()}`;
    await mountSecrets.enable('alicloud', enginePath);
    await settled();
    assert.ok(backendsPage.rows.filterBy('path', `${enginePath}/`)[0], 'shows the mounted engine');

    await backendsPage.visit();
    await settled();
    const row = backendsPage.rows.filterBy('path', `${enginePath}/`)[0];
    await row.menu();
    await settled();
    await backendsPage.disableButton();
    await settled();
    await backendsPage.confirmDisable();
    await settled();
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backends',
      'redirects to the backends page'
    );

    assert.strictEqual(
      backendsPage.rows.filterBy('path', `${enginePath}/`).length,
      0,
      'does not show the disabled engine'
    );
  });
});
