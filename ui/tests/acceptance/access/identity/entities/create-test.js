/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { currentRouteName } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import page from 'vault/tests/pages/access/identity/create';
import { testCRUD, testDeleteFromForm } from '../_shared-tests';
import authPage from 'vault/tests/pages/auth';
import { setRunOptions } from 'ember-a11y-testing/test-support';

module('Acceptance | /access/identity/entities/create', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    // Popup menu causes flakiness
    setRunOptions({
      rules: {
        'color-contrast': { enabled: false },
      },
    });
    return authPage.login();
  });

  test('it visits the correct page', async function (assert) {
    await page.visit({ item_type: 'entities' });
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.identity.create',
      'navigates to the correct route'
    );
  });

  test('it allows create, list, delete of an entity', async function (assert) {
    assert.expect(6);
    const name = `entity-${Date.now()}`;
    await testCRUD(name, 'entities', assert);
  });

  test('it can be deleted from the edit form', async function (assert) {
    assert.expect(6);
    const name = `entity-${Date.now()}`;
    await testDeleteFromForm(name, 'entities', assert);
  });
});
