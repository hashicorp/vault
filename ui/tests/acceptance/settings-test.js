/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { currentURL, visit, click, fillIn } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

import mountSecrets from 'vault/tests/pages/settings/mount-secret-backend';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { deleteEngineCmd, mountEngineCmd, runCmd } from 'vault/tests/helpers/commands';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

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
    await visit('/vault/secrets-engines/enable');

    assert.strictEqual(
      currentURL(),
      '/vault/secrets-engines/enable',
      'navigates to the mount secret backend page'
    );
    await click(GENERAL.cardContainer(type));
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

    assert.strictEqual(currentURL(), `/vault/secrets-engines`, 'redirects to secrets page');
    // cleanup
    await runCmd(deleteEngineCmd(path));
  });

  test('it navigates to ember engine configuration page', async function (assert) {
    const type = 'ldap';
    const path = `ldap-${this.uid}`;

    await visit('/vault/secrets-engines/enable');
    await runCmd(mountEngineCmd(type, path), false);
    await visit('/vault/secrets-engines');
    await fillIn(GENERAL.inputSearch('secret-engine-path'), path);
    await click(GENERAL.menuTrigger);
    await click(GENERAL.menuItem('View configuration'));
    // since ldap hasn't been configured yet, it should redirect to configure page
    assert.strictEqual(
      currentURL(),
      `/vault/secrets-engines/${path}/${type}/configure`,
      'navigates to the config page for ember engine'
    );
    // clean up
    await runCmd(deleteEngineCmd(path));
  });

  test('it navigates to non-ember engine general configuration page', async function (assert) {
    const type = 'keymgmt';
    const path = `keymgmt-${this.uid}`;

    await visit('/vault/secrets-engines/enable');
    await runCmd(mountEngineCmd(type, path), false);
    await visit(`/vault/secrets-engines/${path}/configuration/general-settings`);

    // since non-ember engines haven't been configured yet, it should redirect to general settings page
    assert.strictEqual(
      currentURL(),
      `/vault/secrets-engines/${path}/configuration/general-settings`,
      'navigates to the general settings config page for non-ember engine'
    );
    // clean up
    await runCmd(deleteEngineCmd(path));
  });

  test('it navigates to edit configuration page if engine is configurable and not set', async function (assert) {
    const type = 'ssh';
    const path = `ssh-${this.uid}`;

    await visit('/vault/secrets-engines/enable');
    await runCmd(mountEngineCmd(type, path), false);
    await visit('/vault/secrets-engines');
    await fillIn(GENERAL.inputSearch('secret-engine-path'), path);
    await click(GENERAL.menuTrigger);
    await click(GENERAL.menuItem('View configuration'));

    // since the engine hasn't been configured yet & is configurable, it should redirect to configuration edit page
    assert.strictEqual(
      currentURL(),
      `/vault/secrets-engines/${path}/configuration/edit`,
      'navigates to the config page for configurable engine'
    );
    // clean up
    await runCmd(deleteEngineCmd(path));
  });
});
