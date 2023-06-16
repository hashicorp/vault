/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { kvId } from 'vault/utils/kv-id';

module('Unit | Adapter | kv/data', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.secretMountPath = this.owner.lookup('service:secret-mount-path');
    this.backend = 'kv-backend';
    this.secretMountPath.currentPath = this.backend;
    this.path = 'beep/bop/my-secret';
    this.version = '2';
    this.id = kvId(this.backend, this.path, 'data', this.version);
    this.data = {
      options: {
        cas: 2,
      },
      data: {
        foo: 'bar',
      },
    };
    this.payload = {
      backend: this.backend,
      path: this.path,
      version: 2,
    };
  });

  hooks.afterEach(function () {
    this.store.unloadAll('kv/data');
    this.server.shutdown();
  });

  test('it should make request to correct endpoint on createRecord', async function (assert) {
    assert.expect(1);
    this.server.post(`${this.backend}/data/${this.path}`, () => {
      assert.ok('POST request made to correct endpoint when creating new record');
    });
    const record = this.store.createRecord('kv/data', { backend: this.backend, path: this.path });
    await record.save();
  });

  test('it should make request to correct endpoint on queryRecord', async function (assert) {
    assert.expect(2);
    this.server.get(`${this.backend}/data/${this.path}`, (schema, req) => {
      assert.strictEqual(
        req.queryParams.version,
        this.version,
        'request includes the version flag on queryRecord.'
      );
      assert.ok(true, 'request is made to correct url on queryRecord.');
    });

    await this.store.queryRecord('kv/data', this.payload);
  });

  test('it should make request to correct endpoint on delete latest version', async function (assert) {
    assert.expect(1);
    this.server.get(`${this.backend}/data/${this.path}`, () => {
      return { id: this.id };
    });
    this.server.delete(`${this.backend}/data/${this.path}`, () => {
      assert.ok(true, 'request made to correct endpoint on delete latest version.');
    });

    this.store.pushPayload('kv/data', {
      modelName: 'kv/data',
      id: this.id,
      ...this.payload,
    });
    const record = await this.store.findRecord('kv/data', this.id);
    await record.destroyRecord({ adapterOptions: { deleteType: 'delete-latest-version' } });
  });

  test('it should make request to correct endpoint on delete specific versions', async function (assert) {
    assert.expect(2);
    this.server.get(`${this.backend}/data/${this.path}`, () => {
      return { id: this.id };
    });
    this.server.post(`${this.backend}/data/${this.path}`, (schema, req) => {
      const { versions } = JSON.parse(req.requestBody);
      assert.strictEqual(versions, 2, 'version array is sent in the payload.');
      assert.ok(true, 'request made to correct endpoint on delete specific version.');
    });

    this.store.pushPayload('kv/data', {
      modelName: 'kv/data',
      id: this.id,
      ...this.payload,
    });
    const record = await this.store.findRecord('kv/data', this.id);
    await record.destroyRecord({
      adapterOptions: { deleteType: 'delete-specific-version', deleteVersions: 2 },
    });
  });

  test('it should make request to correct endpoint on undelete', async function (assert) {
    assert.expect(2);
    this.server.get(`${this.backend}/data/${this.path}`, () => {
      return { id: this.id };
    });
    this.server.post(`${this.backend}/undelete/${this.path}`, (schema, req) => {
      const { versions } = JSON.parse(req.requestBody);
      assert.strictEqual(versions, 2, 'version array is sent in the payload.');
      assert.ok(true, 'request made to correct endpoint on undelete specific version.');
    });

    this.store.pushPayload('kv/data', {
      modelName: 'kv/data',
      id: this.id,
      ...this.payload,
    });
    const record = await this.store.findRecord('kv/data', this.id);

    await record.destroyRecord({
      adapterOptions: { deleteType: 'undelete-specific-version', deleteVersions: 2 },
    });
  });

  test('it should make request to correct endpoint on destroy specific versions', async function (assert) {
    assert.expect(2);
    this.server.get(`${this.backend}/data/${this.path}`, () => {
      return { id: this.id };
    });
    this.server.put(`${this.backend}/destroy/${this.path}`, (schema, req) => {
      const { versions } = JSON.parse(req.requestBody);
      assert.strictEqual(versions, 2, 'version array is sent in the payload.');
      assert.ok(true, 'request made to correct endpoint on destroy specific version.');
    });

    this.store.pushPayload('kv/data', {
      modelName: 'kv/data',
      id: this.id,
      ...this.payload,
    });
    const record = await this.store.findRecord('kv/data', this.id);
    await record.destroyRecord({
      adapterOptions: { deleteType: 'destroy-specific-version', deleteVersions: 2 },
    });
  });

  test('it should make request to correct endpoint on destroy everything', async function (assert) {
    assert.expect(1);
    this.server.get(`${this.backend}/data/${this.path}`, () => {
      return { id: this.id };
    });
    this.server.delete(`${this.backend}/metadata/${this.path}`, () => {
      assert.ok(true, 'request made to correct endpoint on destroy everything.');
    });

    this.store.pushPayload('kv/data', {
      modelName: 'kv/data',
      id: this.id,
      ...this.payload,
    });
    const record = await this.store.findRecord('kv/data', this.id);
    await record.destroyRecord({
      adapterOptions: { deleteType: 'destroy-everything' },
    });
  });
});
