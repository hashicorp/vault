/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { currentURL, find, visit, settled, click } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

import backendListPage from 'vault/tests/pages/secrets/backends';
import mountSecrets from 'vault/tests/pages/settings/mount-secret-backend';
import authPage from 'vault/tests/pages/auth';
import { deleteEngineCmd, mountEngineCmd, runCmd } from 'vault/tests/helpers/commands';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

const { searchSelect } = GENERAL;

module('Acceptance | secret engine mount settings', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    this.uid = uuidv4();
    return authPage.login();
  });

  test('it allows you to mount a secret engine', async function (assert) {
    const type = 'consul';
    const path = `settings-path-${this.uid}`;

    // mount unsupported backend
    await visit('/vault/settings/mount-secret-backend');

    assert.strictEqual(currentURL(), '/vault/settings/mount-secret-backend');

    await mountSecrets.selectType(type);
    await mountSecrets
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
    // cleanup
    await runCmd(deleteEngineCmd(path));
  });

  test('it navigates to ember engine configuration page', async function (assert) {
    const type = 'ldap';
    const path = `ldap-${this.uid}`;

    await visit('/vault/settings/mount-secret-backend');
    await runCmd(mountEngineCmd(type, path), false);
    await visit('/vault/secrets');
    await click(searchSelect.trigger('filter-by-engine-name'));
    await click(searchSelect.option(searchSelect.optionIndex(path)));
    await click(GENERAL.menuTrigger);
    await backendListPage.configLink();
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${path}/${type}/configuration`,
      'navigates to the config page for ember engine'
    );
    // clean up
    await runCmd(deleteEngineCmd(path));
  });

  test('it navigates to non-ember engine configuration page', async function (assert) {
    const type = 'ssh';
    const path = `ssh-${this.uid}`;

    await visit('/vault/settings/mount-secret-backend');
    await runCmd(mountEngineCmd(type, path), false);
    await visit('/vault/secrets');
    await click(searchSelect.trigger('filter-by-engine-name'));
    await click(searchSelect.option(searchSelect.optionIndex(path)));
    await click(GENERAL.menuTrigger);
    await backendListPage.configLink();
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${path}/configuration`,
      'navigates to the config page for non-ember engine'
    );
    // clean up
    await runCmd(deleteEngineCmd(path));
  });
});
