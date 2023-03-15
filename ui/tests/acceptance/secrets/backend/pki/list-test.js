/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { click, currentRouteName, settled } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import page from 'vault/tests/pages/secrets/backend/list';
import authPage from 'vault/tests/pages/auth';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';

module('Acceptance | secrets/pki/list', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    this.timestamp = new Date().getTime();
    return authPage.login();
  });

  const mountAndNav = async (uid) => {
    const path = `pki-${uid}`;
    await enablePage.enable('pki', path);
    await page.visitRoot({ backend: path });
  };

  test('it renders an empty list', async function (assert) {
    assert.expect(5);
    await mountAndNav(this.timestamp);
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.list-root',
      'redirects from the index'
    );
    assert.ok(page.createIsPresent, 'create button is present');
    await click('[data-test-configuration-tab]');
    assert.ok(page.configureIsPresent, 'configure button is present');
    assert.strictEqual(page.tabs.length, 2, 'shows 2 tabs');
    assert.ok(page.backendIsEmpty);
  });

  test('it navigates to the create page', async function (assert) {
    assert.expect(1);
    await mountAndNav(this.timestamp);
    await page.create();
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.create-root',
      'links to the create page'
    );
  });

  test('it navigates to the configure page', async function (assert) {
    assert.expect(1);
    await mountAndNav(this.timestamp);
    await click('[data-test-configuration-tab]');
    await page.configure();
    await settled();
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.settings.configure-secret-backend.section',
      'links to the configure page'
    );
  });
});
