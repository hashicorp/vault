/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Serializer | kv/metadata', function (hooks) {
  setupTest(hooks);

  test('it should properly normalize a list response', function (assert) {
    const serializer = this.owner.lookup('serializer:kv/metadata');
    const serverData = {
      request_id: 'foo',
      backend: 'my-backend',
      data: {
        keys: ['first', 'second'],
      },
    };
    const expectedData = [
      {
        id: 'my-backend/metadata/first',
        path: 'first',
        backend: 'my-backend',
      },
      {
        id: 'my-backend/metadata/second',
        path: 'second',
        backend: 'my-backend',
      },
    ];

    const serializedRecord = serializer.normalizeItems(serverData);
    assert.deepEqual(serializedRecord.data.keys, expectedData, 'transformed keys into proper IDs');
  });

  test('it throws an assertionn if backend not on payload', function (assert) {
    const serializer = this.owner.lookup('serializer:kv/metadata');
    const serverData = {
      request_id: 'foo',
      data: {
        keys: ['first', 'second'],
      },
    };
    let result;
    try {
      result = serializer.normalizeItems(serverData);
    } catch (e) {
      result = e.message;
    }
    assert.strictEqual(
      result,
      'Assertion Failed: payload.backend must be provided on kv/metadata list response'
    );
  });
});
