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

  test('it calls the correct endpoint for createRecord', async function (assert) {
    const types = destinationTypes();
    assert.expect(types.length);

    for (const type of types) {
      const name = 'my-dest';
      this.server.post(`sys/sync/destinations/${type}/${name}`, () => {
        assert.ok(true, `request is made to GET sys/sync/destinations/${type}/my-dest endpoint on create`);
        return {
          data: {
            connection_details: {},
            name,
            type,
          },
        };
      });
      this.model = this.store.createRecord(`sync/destinations/${type}`, { type, name });
      this.model.save();
    }
  });

  test('it calls the correct endpoint for findRecord', async function (assert) {
    const types = destinationTypes();
    assert.expect(types.length);

    for (const type of types) {
      const name = 'my-dest';
      this.server.get(`sys/sync/destinations/${type}/${name}`, () => {
        assert.ok(true, `request is made to GET sys/sync/destinations/${type}/${name} endpoint on find`);
        return {
          data: {
            connection_details: {},
            name,
            type,
          },
        };
      });
      this.store.findRecord(`sync/destinations/${type}`, name);
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

  test('it should make request to correct endpoint for deleteRecord', async function (assert) {
    assert.expect(1);

    const destination = this.server.create('sync-destination', 'aws-sm');

    this.server.delete(`/sys/sync/destinations/${destination.type}/${destination.name}`, () => {
      assert.ok(true, 'DELETE request made to correct endpoint');
      return {};
    });

    const modelName = 'sync/destinations/aws-sm';
    this.store.pushPayload(modelName, {
      modelName,
      ...destination,
      id: destination.name,
    });
    const model = this.store.peekRecord(modelName, destination.name);
    await model.destroyRecord();
  });
});
