/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import sinon from 'sinon';

module('Unit | Route | vault/cluster/secrets/backend/list', function (hooks) {
  setupTest(hooks);

  test('it filters generic secrets before fetching strict capabilities for the page', async function (assert) {
    const route = this.owner.lookup('route:vault.cluster.secrets.backend.list');
    const list = sinon.stub(route.api.secrets, 'kvV1List').resolves({
      keys: ['first', 'matching-secret', 'folder/'],
    });
    const fetch = sinon.stub(route.capabilitiesService, 'fetch').resolves({
      'generic/matching-secret': { canRead: true, canUpdate: true, canDelete: false },
    });

    const items = await route.fetchGenericSecrets('generic', '', 1, 'matching');

    assert.deepEqual(
      list.firstCall.args,
      ['', 'generic', 'true'],
      'the route lists keys from the requested KV v1 path'
    );
    assert.deepEqual(
      fetch.firstCall.args,
      [['generic/matching-secret'], { strict: true }],
      'capabilities are batched only for filtered page items'
    );
    assert.strictEqual(items.length, 1, 'the page contains only the matching secret');
    assert.strictEqual(items[0].id, 'matching-secret', 'the listed key becomes the item id');
    assert.true(items[0].canRead, 'read permission is mapped to the list item');
    assert.true(items[0].canEdit, 'update permission is mapped to the list item');
    assert.false(items[0].canDelete, 'delete permission is mapped to the list item');
    assert.false(items[0].secretPath.isLoading, 'the item does not show a capability loading state');
    assert.strictEqual(items.meta.total, 3, 'pagination metadata retains the total key count');
    assert.strictEqual(items.meta.filteredTotal, 1, 'pagination metadata retains the filtered key count');
  });
});
