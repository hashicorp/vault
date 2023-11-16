/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Serializer | kv/data', function (hooks) {
  setupTest(hooks);

  test('it should always pass the cas option when creating/updating a secret', function (assert) {
    const store = this.owner.lookup('service:store');
    const record = store.createRecord('kv/data', {
      path: 'my-secret-path',
      backend: 'kv-test',
      version: 2,
      casVersion: 3,
      secretData: { foo: 'bar' },
    });
    const expectedResult = {
      data: { foo: 'bar' },
      options: {
        cas: 3,
      },
    };

    const serializedRecord = record.serialize();
    assert.deepEqual(serializedRecord, expectedResult, 'cas option was correctly added to the payload.');
  });
});
