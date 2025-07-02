/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { currentURL, visit, click, fillIn } from '@ember/test-helpers';
import { selectChoose } from 'ember-power-select/test-support';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

import mountSecrets from 'vault/tests/pages/settings/mount-secret-backend';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { deleteEngineCmd, mountEngineCmd, runCmd } from 'vault/tests/helpers/commands';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { MOUNT_BACKEND_FORM } from 'vault/tests/helpers/components/mount-backend-form-selectors';

module('Acceptance | secret engine mount settings', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    this.uid = uuidv4();
    return login();
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
    await click(GENERAL.button('Method Options'));
    await click(GENERAL.toggleInput('Default Lease TTL'));
    await mountSecrets.defaultTTLUnit('s').defaultTTLVal(100);
    await click(GENERAL.submitButton);

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
    await selectChoose(GENERAL.searchSelect.trigger('filter-by-engine-name'), path);
    await click(GENERAL.menuTrigger);
    await click(GENERAL.menuItem('view-configuration'));
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
    await selectChoose(GENERAL.searchSelect.trigger('filter-by-engine-name'), path);
    await click(GENERAL.menuTrigger);
    await click(GENERAL.menuItem('view-configuration'));
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${path}/configuration`,
      'navigates to the config page for non-ember engine'
    );
    // clean up
    await runCmd(deleteEngineCmd(path));
  });
});
