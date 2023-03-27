/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { currentRouteName, settled, find, waitUntil } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

import page from 'vault/tests/pages/settings/configure-secret-backends/pki/section';
import authPage from 'vault/tests/pages/auth';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';

module('Acceptance | settings/configure/secrets/pki/urls', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    this.uid = uuidv4();
    return authPage.login();
  });

  test('it saves urls config', async function (assert) {
    const path = `pki-config-urls-${this.uid}`;
    await enablePage.enable('pki', path);
    await settled();
    await page.visit({ backend: path, section: 'urls' });
    await settled();
    assert.strictEqual(currentRouteName(), 'vault.cluster.settings.configure-secret-backend.section');

    await page.form.fields.objectAt(0).textarea('foo').change();
    await page.form.submit();
    await waitUntil(() => find('[data-test-error]'));
    assert.ok(page.form.hasError, 'shows error on invalid input');

    await page.form.fields.objectAt(0).textarea('foo.example.com').change();
    await page.form.submit();
    await settled();
    assert.strictEqual(page.lastMessage, 'The urls config for this backend has been updated.');
  });
});
