/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { currentRouteName, settled } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

import page from 'vault/tests/pages/settings/configure-secret-backends/pki/section';
import authPage from 'vault/tests/pages/auth';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';

module('Acceptance | settings/configure/secrets/pki/tidy', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    return authPage.login();
  });

  test('it saves tidy config', async function (assert) {
    const path = `pki-tidy-${uuidv4()}`;
    await enablePage.enable('pki', path);
    await settled();
    await page.visit({ backend: path, section: 'tidy' });
    await settled();
    assert.strictEqual(currentRouteName(), 'vault.cluster.settings.configure-secret-backend.section');
    await page.form.fields.objectAt(0).clickLabel();

    await page.form.submit();
    await settled();
    assert.strictEqual(page.lastMessage, 'The tidy config for this backend has been updated.');
  });
});
