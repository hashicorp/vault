/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { setupTest } from 'ember-qunit';
import { module, test } from 'qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { destinationTypes } from 'vault/helpers/sync-destinations';

module('Unit | Adapter | sync | destination', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
  });

  test('it calls the correct endpoint for findRecord', async function (assert) {
    const types = destinationTypes();
    assert.expect(types.length);

    for (const type of types) {
      this.server.get(`sys/sync/destinations/${type}/my-dest`, () => {
        assert.ok(true, `request is made to GET sys/sync/destinations/${type}/my-dest endpoint on find`);
        return {
          data: {
            connection_details: {},
            name: 'my-dest',
            type,
          },
        };
      });
      this.store.findRecord(`sync/destinations/${type}`, 'my-dest');
    }
  });

  test('it calls the correct endpoint for query', async function (assert) {
    assert.expect(2);

    this.server.get('sys/sync/destinations', (schema, req) => {
      assert.propEqual(req.queryParams, { list: 'true' }, 'it passes { list: true } as query params');
      assert.ok(true, `request is made to LIST sys/sync/destinations endpoint on query`);
      return {
        data: {
          key_info: {
            'aws-sm': ['my-dest-1'],
            gh: ['my-dest-1'],
          },
          keys: ['aws-sm', 'gh'],
        },
      };
    });
    this.store.query('sync/destination', {});
  });
});
