/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, currentRouteName, visit } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';
import { mountBackend } from 'vault/tests/helpers/components/mount-backend-form-helpers';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { deleteAuthCmd, runCmd } from 'vault/tests/helpers/commands';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Acceptance | settings/auth/enable', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    this.uid = uuidv4();
    return login();
  });

  test('it mounts and redirects', async function (assert) {
    // always force the new mount to the top of the list
    const path = `aaa-approle-${this.uid}`;
    const type = 'approle';
    await visit('/vault/settings/auth/enable');
    assert.strictEqual(currentRouteName(), 'vault.cluster.settings.auth.enable');
    await mountBackend(type, path);
    assert
      .dom(GENERAL.latestFlashContent)
      .hasText(`Successfully mounted the ${type} auth method at ${path}.`);
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.settings.auth.configure.section',
      'redirects to the auth config page'
    );

    await visit('/vault/access/');
    assert.dom(`[data-test-auth-backend-link=${path}]`).exists('mount is present in the list');

    // cleanup
    await runCmd(deleteAuthCmd(path));
  });

  test('it renders default config details', async function (assert) {
    const path = `approle-config-${this.uid}`;
    const type = 'approle';
    await visit('/vault/settings/auth/enable');
    await mountBackend(type, path);
    // the config details is updated to query mount details from sys/internal/ui/mounts
    // but we still want these forms to continue using sys/auth which returns 0 for default ttl values
    // check tune form (right after enabling)
    assert.dom(GENERAL.toggleInput('Default Lease TTL')).isNotChecked('default lease ttl is unset');
    assert.dom(GENERAL.toggleInput('Max Lease TTL')).isNotChecked('max lease ttl is unset');
    await click(GENERAL.breadcrumbAtIdx(1));
    assert
      .dom(GENERAL.infoRowValue('Default Lease TTL'))
      .hasText('1 month 1 day', 'shows system default TTL');
    assert.dom(GENERAL.infoRowValue('Max Lease TTL')).hasText('1 month 1 day', 'shows the proper max TTL');

    // check edit form TTL values
    await click('[data-test-configure-link]');
    assert.dom(GENERAL.toggleInput('Default Lease TTL')).isNotChecked('default lease ttl is still unset');
    assert.dom(GENERAL.toggleInput('Max Lease TTL')).isNotChecked('max lease ttl is still unset');

    // cleanup
    await runCmd(deleteAuthCmd(path));
  });
});
