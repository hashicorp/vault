/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';

const TRANSFORM_TYPES = ['fpe', 'masking', 'tokenization'];
module('Unit | Adapter | transform', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.backend = 'my-transform-engine';
    this.name = 'my-transform';
  });

  hooks.afterEach(function () {
    this.store.unloadAll('transform');
  });

  test('it should make request to correct endpoint when querying all records', async function (assert) {
    assert.expect(2);
    this.server.get(`${this.backend}/transformation`, (schema, req) => {
      assert.ok(true, 'GET request made to correct endpoint when querying record');
      assert.propEqual(req.queryParams, { list: 'true' }, 'query params include list: true');
      return { data: { key_info: {}, keys: [] } };
    });
    await this.store.query('transform', { backend: this.backend });
  });

  test('it should make request to correct endpoint when querying a record', async function (assert) {
    assert.expect(1);
    this.server.get(`${this.backend}/transformation/${this.name}`, () => {
      assert.ok(true, 'GET request made to correct endpoint when querying record');
      return { data: { backend: this.backend, name: this.name } };
    });
    await this.store.queryRecord('transform', { backend: this.backend, id: this.name });
  });

  test('it should make request to correct endpoint when creating new record', async function (assert) {
    assert.expect(3);

    for (const type of TRANSFORM_TYPES) {
      const name = `transform-${type}-test`;
      this.server.post(`${this.backend}/transformations/${type}/${name}`, () => {
        assert.ok(true, `POST request made to transformations/${type}/:name creating a record`);
        return { data: { backend: this.backend, name, type } };
      });
      const record = this.store.createRecord('transform', { backend: this.backend, name, type });
      await record.save();
    }
  });

  test('it should make request to correct endpoint when updating record', async function (assert) {
    assert.expect(3);
    for (const type of TRANSFORM_TYPES) {
      const name = `transform-${type}-test`;
      this.server.post(`${this.backend}/transformations/${type}/${name}`, () => {
        assert.ok(true, `POST request made to transformations/${type}/:name endpoint`);
      });
      this.store.pushPayload('transform', {
        modelName: 'transform',
        backend: this.backend,
        id: name,
        type,
        name,
      });
      const record = this.store.peekRecord('transform', name);
      await record.save();
    }
  });

  test('it should make request to correct endpoint when deleting record', async function (assert) {
    assert.expect(3);
    for (const type of TRANSFORM_TYPES) {
      const name = `transform-${type}-test`;
      this.server.delete(`${this.backend}/transformation/${name}`, () => {
        assert.ok(true, `type: ${type} - DELETE request to transformation/:name endpoint`);
      });
      this.store.pushPayload('transform', {
        modelName: 'transform',
        backend: this.backend,
        id: name,
        type,
        name,
      });
      const record = this.store.peekRecord('transform', name);
      await record.destroyRecord();
    }
  });
});
