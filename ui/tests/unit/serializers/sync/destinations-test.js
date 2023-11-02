/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'vault/tests/helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { destinationTypes } from 'vault/helpers/sync-destinations';

module('Unit | Serializer | sync | destination', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.serializer = this.owner.lookup('serializer:sync/destination');
  });

  test('it normalizes findRecord payload from the server', async function (assert) {
    const types = destinationTypes();
    assert.expect(types.length);

    for (const destinationType of types) {
      const { name, type, id, ...connection_details } = this.server.create(
        'sync-destination',
        destinationType
      );
      const serverData = { request_id: id, data: { name, type, connection_details } };

      const normalized = this.serializer._normalizePayload(serverData);
      const expected = { data: { id: name, type, name, ...connection_details } };
      assert.propEqual(
        normalized,
        expected,
        `generates id and adds connection details to ${destinationType} data object`
      );
    }
  });

  test('it normalizes query payload from the server', async function (assert) {
    assert.expect(1);
    // hardcoded from docs https://developer.hashicorp.com/vault/api-docs/system/secrets-sync#sample-response
    // destinations intentionally named the same to test no id naming collision happens
    const serverData = {
      data: {
        key_info: {
          'aws-sm': ['my-dest-1'],
          gh: ['my-dest-1'],
        },
        keys: ['aws-sm', 'gh'],
      },
    };

    const normalized = this.serializer.extractLazyPaginatedData(serverData);
    const expected = [
      {
        id: 'aws-sm/my-dest-1',
        name: 'my-dest-1',
        type: 'aws-sm',
      },
      {
        id: 'gh/my-dest-1',
        name: 'my-dest-1',
        type: 'gh',
      },
    ];

    assert.propEqual(normalized, expected, 'payload is array of objects with concatenated type/name as id');
  });
});
