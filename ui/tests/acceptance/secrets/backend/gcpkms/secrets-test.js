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

module('Acceptance | gcpkms/enable', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    this.timestamp = new Date().getTime();
    return authPage.login();
  });

  test('enable gcpkms', async function (assert) {
    // Error: Cannot call `visit` without having first called `setupApplicationContext`.
    const enginePath = `gcpkms-${this.timestamp}`;
    await mountSecrets.visit();
    await settled();
    await mountSecrets.selectType('gcpkms');
    await mountSecrets.next().path(enginePath).submit();
    await settled();
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backends',
      'redirects to the backends page'
    );
    assert.ok(backendsPage.rows.filterBy('path', `${enginePath}/`)[0], 'shows the gcpkms engine');
  });
});
