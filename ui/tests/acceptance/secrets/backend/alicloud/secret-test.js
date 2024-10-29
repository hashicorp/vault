/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { currentRouteName, settled, click, fillIn } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

import mountSecrets from 'vault/tests/pages/settings/mount-secret-backend';
import backendsPage from 'vault/tests/pages/secrets/backends';
import authPage from 'vault/tests/pages/auth';
import { MOUNT_BACKEND_FORM } from 'vault/tests/helpers/components/mount-backend-form-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Acceptance | alicloud/enable', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    this.uid = uuidv4();
    return authPage.login();
  });

  test('enable alicloud', async function (assert) {
    const enginePath = `alicloud-${this.uid}`;
    await mountSecrets.visit();
    await settled();
    await click(MOUNT_BACKEND_FORM.mountType('alicloud'));
    await fillIn(GENERAL.inputByAttr('path'), enginePath);
    await click(GENERAL.saveButton);

    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backends',
      'redirects to the backends page'
    );
    await settled();
    assert.ok(backendsPage.rows.filterBy('path', `${enginePath}/`)[0], 'shows the alicloud engine');
  });
});
