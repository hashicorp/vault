/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { currentRouteName, settled } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

import mountSecrets from 'vault/tests/pages/settings/mount-secret-backend';
import backendsPage from 'vault/tests/pages/secrets/backends';
import authPage from 'vault/tests/pages/auth';
import { mountBackend } from 'vault/tests/helpers/components/mount-backend-form-helpers';

module('Acceptance | gcpkms/enable', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    this.uid = uuidv4();
    return authPage.login();
  });

  test('enable gcpkms', async function (assert) {
    // Error: Cannot call `visit` without having first called `setupApplicationContext`.
    const enginePath = `gcpkms-${this.uid}`;
    await mountSecrets.visit();
    await settled();
    await mountBackend('gcpkms', enginePath);

    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backends',
      'redirects to the backends page'
    );
    assert.ok(backendsPage.rows.filterBy('path', `${enginePath}/`)[0], 'shows the gcpkms engine');
  });
});
