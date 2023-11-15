/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { setupTest } from 'ember-qunit';
import { module, test } from 'qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { associationsResponse } from 'vault/mirage/handlers/sync';

module('Unit | Adapter | sync | association', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.params = { type: 'aws-sm', name: 'us-west-1' };
    this.destination = this.server.create('sync-destination', this.params.type, { name: this.params.name });
    this.association = this.server.create('sync-association', {
      ...this.params,
      mount: 'foo',
      secret_name: 'bar',
    });
    this.newModel = this.store.createRecord('sync/association', {
      destinationType: 'aws-sm',
      destinationName: 'us-west-1',
      mount: 'foo',
      secretName: 'bar',
    });
  });

  test('it should make request to correct endpoint when querying', async function (assert) {
    assert.expect(2);

    this.server.get('/sys/sync/destinations/:type/:name/associations', (schema, req) => {
      // list query param not required for this endpoint
      assert.deepEqual(req.queryParams, {}, 'query params stripped from request');
      assert.deepEqual(req.params, this.params, 'request is made to correct endpoint when querying');
      return associationsResponse(schema, req);
    });

    await this.store.lazyPaginatedQuery('sync/association', {
      responsePath: 'data.keys',
      page: 1,
      destinationType: 'aws-sm',
      destinationName: 'us-west-1',
    });
  });

  test('it should make request to correct endpoint when creating record', async function (assert) {
    assert.expect(2);

    this.server.post('/sys/sync/destinations/:type/:name/associations/set', (schema, req) => {
      assert.deepEqual(req.params, this.params, 'request is made to correct endpoint when querying');
      assert.deepEqual(
        JSON.parse(req.requestBody),
        { mount: 'foo', secret_name: 'bar' },
        'Correct payload is sent when creating association'
      );
      return associationsResponse(schema, req);
    });

    await this.newModel.save({ adapterOptions: { action: 'set' } });
  });

  test('it should make request to correct endpoint when updating record', async function (assert) {
    assert.expect(2);

    this.server.post('/sys/sync/destinations/:type/:name/associations/remove', (schema, req) => {
      assert.deepEqual(req.params, this.params, 'request is made to correct endpoint when querying');
      assert.deepEqual(
        JSON.parse(req.requestBody),
        { mount: 'foo', secret_name: 'bar' },
        'Correct payload is sent when removing association'
      );
      return associationsResponse(schema, req);
    });

    this.store.pushPayload('sync/association', {
      modelName: 'sync/association',
      destinationType: 'aws-sm',
      destinationName: 'us-west-1',
      mount: 'foo',
      secret_name: 'bar',
      sync_status: 'SYNCED',
      id: 'foo/bar',
    });
    const model = this.store.peekRecord('sync/association', 'foo/bar');

    await model.save({ adapterOptions: { action: 'remove' } });
  });

  test('it should parse response from set/remove request', async function (assert) {
    this.server.post('/sys/sync/destinations/:type/:name/associations/set', associationsResponse);

    const adapter = this.store.adapterFor('sync/association');
    // mock snapshot
    const snapshot = {
      attributes() {
        return { destinationName: 'us-west-1', destinationType: 'aws-sm' };
      },
      serialize() {
        return { mount: 'foo', secret_name: 'bar' };
      },
      adapterOptions: { action: 'set' },
    };
    const response = await adapter._setOrRemove(this.store, { modelName: 'sync/association' }, snapshot);
    const { accessor, mount, secret_name, sync_status, name, type, updated_at } = this.association;
    const expected = {
      id: 'foo/bar',
      accessor,
      mount,
      secret_name,
      sync_status,
      updated_at,
      destinationType: type,
      destinationName: name,
    };

    assert.deepEqual(
      response,
      expected,
      'Custom create/update record method makes request and parses response'
    );
  });

  test('it should throw error if save action is not passed in adapterOptions', async function (assert) {
    assert.expect(1);

    try {
      await this.newModel.save();
    } catch (e) {
      assert.strictEqual(
        e.message,
        "Assertion Failed: action type of set or remove required when saving association => association.save({ adapterOptions: { action: 'set' }})"
      );
    }
  });
});
