/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { currentURL, visit, click, fillIn } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

import backendListPage from 'vault/tests/pages/secrets/backends';
import mountSecrets from 'vault/tests/pages/settings/mount-secret-backend';
import authPage from 'vault/tests/pages/auth';
import { deleteEngineCmd, mountEngineCmd, runCmd } from 'vault/tests/helpers/commands';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { MOUNT_BACKEND_FORM } from 'vault/tests/helpers/components/mount-backend-form-selectors';

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

    assert.strictEqual(
      currentURL(),
      '/vault/settings/mount-secret-backend',
      'navigates to the mount secret backend page'
    );
    await click(MOUNT_BACKEND_FORM.mountType(type));
    await fillIn(GENERAL.inputByAttr('path'), path);
    await click(GENERAL.toggleGroup('Method Options'));
    await mountSecrets.enableDefaultTtl().defaultTTLUnit('s').defaultTTLVal(100);
    await click(GENERAL.saveButton);

    assert
      .dom(`${GENERAL.flashMessage}.is-success`)
      .includesText(
        `Success Successfully mounted the ${type} secrets engine at ${path}`,
        'flash message is shown after mounting'
      );

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
