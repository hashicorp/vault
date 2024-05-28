/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { settled } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import { testAliasCRUD, testAliasDeleteFromForm } from '../../_shared-alias-tests';
import authPage from 'vault/tests/pages/auth';
import { v4 as uuidv4 } from 'uuid';
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

  test('it allows create, list, delete of an entity alias', async function (assert) {
    assert.expect(6);
    const name = `alias-${uuidv4()}`;
    await testAliasCRUD(name, 'groups', assert);
    await settled();
  });

  test('it allows delete from the edit form', async function (assert) {
    assert.expect(4);
    const name = `alias-${uuidv4()}`;
    await testAliasDeleteFromForm(name, 'groups', assert);
    await settled();
  });
});
