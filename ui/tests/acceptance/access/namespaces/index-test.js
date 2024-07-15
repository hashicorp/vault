/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { currentRouteName, visit } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { login } from 'vault/tests/helpers/auth/auth-helpers';

module('Acceptance | Enterprise | /access/namespaces', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    return login();
  });

  test('it navigates to namespaces page', async function (assert) {
    assert.expect(1);
    await visit('/vault/access/namespaces');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.namespaces.index',
      'navigates to the correct route'
    );
  });

  test('it should render correct number of namespaces', async function (assert) {
    assert.expect(3);
    await visit('/vault/access/namespaces');
    const store = this.owner.lookup('service:store');
    // Default page size is 15
    assert.strictEqual(store.peekAll('namespace').length, 15, 'Store has 15 namespaces records');
    assert.dom('.list-item-row').exists({ count: 15 });
    assert.dom('.hds-pagination').exists();
  });
});
