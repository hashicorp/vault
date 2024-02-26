/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, skip, test } from 'qunit';
import { settled } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import { testAliasCRUD, testAliasDeleteFromForm } from '../../_shared-alias-tests';
import authPage from 'vault/tests/pages/auth';
import { setRunOptions } from 'ember-a11y-testing/test-support';

module('Acceptance | /access/identity/groups/aliases/add', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    // Popup menu causes flakiness
    setRunOptions({
      rules: {
        'color-contrast': { enabled: false },
      },
    });
    await authPage.login();
    return;
  });

  skip('it allows create, list, delete of an entity alias', async function (assert) {
    // TODO figure out what is wrong with this test
    assert.expect(6);
    const name = `alias-${Date.now()}`;
    await testAliasCRUD(name, 'groups', assert);
    await settled();
  });

  test('it allows delete from the edit form', async function (assert) {
    // TODO figure out what is wrong with this test
    assert.expect(4);
    const name = `alias-${Date.now()}`;
    await testAliasDeleteFromForm(name, 'groups', assert);
    await settled();
  });
});
