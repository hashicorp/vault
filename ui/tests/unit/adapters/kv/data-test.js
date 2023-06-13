/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { encodePath } from 'vault/utils/path-encoding-helpers';

module('Unit | Adapter | kv/data', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.secretMountPath = this.owner.lookup('service:secret-mount-path');
    this.backend = 'kv-backend';
    this.secretMountPath.currentPath = this.backend;
    this.path = 'my-secret-path';
    this.data = {
      options: {
        cas: 2,
      },
      data: {
        foo: 'bar',
      },
    };
  });

  hooks.afterEach(function () {
    this.store.unloadAll('kv/data');
    this.server.shutdown();
  });

  test('it should make request to correct endpoint on queryRecord', async function (assert) {
    assert.expect(2);
    this.server.get(`${this.backend}/data/${this.path}`, (schema, req) => {
      assert.strictEqual(req.queryParams.version, '2', 'request includes the version flag on queryRecord.');
      assert.ok(true, 'request is made to correct url on queryRecord.');
    });

    this.store.queryRecord('kv/data', { backend: this.backend, path: this.path, version: 2 });
  });

  test('it should make request to correct endpoint on createRecord', async function (assert) {
    assert.expect(1);
    this.server.post(`${this.backend}/data/${this.path}`, () => {
      assert.ok('POST request made to correct endpoint when creating new record');
    });
    const record = this.store.createRecord('kv/data', { backend: this.backend, path: this.path });
    await record.save();
  });

  test('it should make request to correct endpoint on delete', async function (assert) {
    assert.expect(1);
    const id = `${encodePath(this.backend)}/2/${encodePath(this.path)}`;

    this.server.get(`${this.backend}/data/${this.path}`, () => {});

    this.server.delete(`${this.backend}/key/${this.data.key_id}`, (schema, req) => {
      assert.strictEqual(req.queryParams.version, '2', 'request includes the version flag on queryRecord.');
      assert.ok(true, 'request made to correct endpoint on delete');
    });

    this.store.pushPayload('kv/data', {
      modelName: 'kv/data',
      backend: this.backend,
      path: this.path,
      version: 2,
      deleteVersions: [2],
      id,
    });

    const model = this.store.peekRecord('kv/data', id);

    await model.destroyRecord();
  });
});
