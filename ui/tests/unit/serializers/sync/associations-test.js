/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'vault/tests/helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';

module('Unit | Serializer | sync | association', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.serializer = this.owner.lookup('serializer:sync/association');
  });

  test('it normalizes query payload from the server', async function (assert) {
    const updated_at = '2023-09-20T10:51:53.961861096-04:00';
    const destinationName = 'us-west-1';
    const destinationType = 'aws-sm';
    const associations = [
      { mount: 'foo', secret_name: 'bar', sync_status: 'SYNCED', updated_at },
      { mount: 'test', secret_name: 'my-secret', sync_status: 'UNSYNCED', updated_at },
    ];
    const payload = {
      data: {
        associated_secrets: {
          'foo/bar': associations[0],
          'test/my-secret': associations[1],
        },
        store_name: destinationName,
        store_type: destinationType,
      },
    };
    const expected = [
      { id: 'foo/bar', destinationName, destinationType, ...associations[0] },
      { id: 'test/my-secret', destinationName, destinationType, ...associations[1] },
    ];
    const normalized = this.serializer.extractLazyPaginatedData(payload);

    assert.deepEqual(normalized, expected, 'lazy paginated data is extracted from payload');
  });
});
