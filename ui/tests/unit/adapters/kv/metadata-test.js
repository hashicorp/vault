/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { kvMetadataPath } from 'vault/utils/kv-path';
import { Response } from 'miragejs';

const EXAMPLE_KV_METADATA_GET_RESPONSE = {
  request_id: 'foobar',
  data: {
    cas_required: true,
    created_time: 'created-time',
    current_version: 2,
    custom_metadata: { application: 'staging' },
    delete_version_after: '0s',
    max_versions: 10,
    oldest_version: 0, // TODO: is this a bug? payload from real API
    updated_time: 'updated-time',
    versions: {
      1: {
        created_time: 'created-time',
        deletion_time: 'deletion-time',
        destroyed: false,
      },
      2: { created_time: 'created-time', deletion_time: '', destroyed: false },
    },
  },
};

module('Unit | Adapter | kv/metadata', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.secretMountPath = this.owner.lookup('service:secret-mount-path');
    this.backend = 'some/kv-back&end';
    this.secretMountPath.currentPath = this.backend;
    this.path = 'beep/bop my/secret';
    this.id = kvMetadataPath(this.backend, this.path);
    this.data = {
      options: {
        cas: 2,
      },
      data: {
        foo: 'bar',
      },
    };
    this.payload = {
      max_versions: 2,
      cas_required: false,
      delete_version_after: '0s',
      custom_metadata: {
        admin: 'bob',
      },
    };
    this.endpoint = kvMetadataPath(this.backend, this.path);
  });

  test('it should make request to correct endpoint on createRecord', async function (assert) {
    assert.expect(10);
    const recordData = {
      backend: this.backend,
      path: this.path,
      deleteVersionAfter: '45h',
      customMetadata: { application: 'staging' },
      oldestVersion: 4,
      currentVersion: 6,
      createdTime: 'created',
      updatedTime: 'updated',
      versions: {
        1: {
          created_time: 'created-time',
          deletion_time: 'deletion-time',
          destroyed: false,
        },
        2: {
          created_time: 'created-time',
          deletion_time: '',
          destroyed: false,
        },
      },
    };
    const expectedBody = {
      max_versions: 0,
      delete_version_after: '45h',
      cas_required: false,
      custom_metadata: { application: 'staging' },
    };
    this.server.post(this.endpoint, (schema, req) => {
      const body = JSON.parse(req.requestBody);
      assert.ok('POST request made to correct endpoint when creating new record');
      assert.propEqual(body, expectedBody, 'POST request has correct body');
      return new Response(204);
    });
    const record = this.store.createRecord('kv/metadata', recordData);
    await record.save();
    assert.strictEqual(record.id, this.id, 'record has correct id');
    assert.strictEqual(record.backend, this.backend, 'record has correct backend');
    assert.strictEqual(record.path, this.path, 'record has correct path');
    assert.strictEqual(record.maxVersions, 0, 'record has default maxVersions');
    assert.false(record.casRequired, 'record has correct casRequired');
    assert.strictEqual(record.deleteVersionAfter, '45h', 'record has correct deleteVersionAfter');
    assert.deepEqual(record.customMetadata, { application: 'staging' }, 'record has correct customMetadata');
    assert.deepEqual(
      record.versions,
      EXAMPLE_KV_METADATA_GET_RESPONSE.data.versions,
      'record has correct versions data'
    );
  });

  test('it should make request to correct endpoint on update record', async function (assert) {
    assert.expect(1);
    const data = this.server.create('kv-metadatum');
    data.id = kvMetadataPath('kv-engine', 'my-secret');
    this.store.pushPayload('kv/metadata', {
      modelName: 'kv/metadata',
      ...data,
    });
    this.server.post(kvMetadataPath('kv-engine', 'my-secret'), () => {
      assert.ok(true, 'request made to correct endpoint on delete metadata.');
    });

    const record = await this.store.peekRecord('kv/metadata', data.id);
    await record.save();
  });

  test('it should make request to correct endpoint on queryRecord', async function (assert) {
    assert.expect(13);
    this.server.get(this.endpoint, () => {
      assert.ok(true, 'request is made to correct url on queryRecord.');
      return EXAMPLE_KV_METADATA_GET_RESPONSE;
    });

    const record = await this.store.queryRecord('kv/metadata', { backend: this.backend, path: this.path });
    assert.strictEqual(record.id, this.id, 'record has correct id');
    assert.strictEqual(record.backend, this.backend, 'record has correct backend');
    assert.strictEqual(record.path, this.path, 'record has correct path');
    assert.strictEqual(record.maxVersions, 10, 'record has correct maxVersions');
    assert.true(record.casRequired, 'record has correct casRequired');
    assert.strictEqual(record.deleteVersionAfter, '0s', 'record has correct deleteVersionAfter');
    assert.deepEqual(record.customMetadata, { application: 'staging' }, 'record has correct customMetadata');
    assert.strictEqual(record.createdTime, 'created-time', 'record has correct createdTime');
    assert.strictEqual(record.currentVersion, 2, 'record has correct currentVersion');
    assert.strictEqual(record.oldestVersion, 0, 'record has correct oldestVersion');
    assert.strictEqual(record.updatedTime, 'updated-time', 'record has correct updatedTime');
    assert.deepEqual(
      record.versions,
      EXAMPLE_KV_METADATA_GET_RESPONSE.data.versions,
      'record has correct versions data'
    );
  });

  test('it should make request to correct endpoint on query', async function (assert) {
    assert.expect(1);
    this.server.get(kvMetadataPath(this.backend, 'directory/'), (schema, req) => {
      assert.ok(req.queryParams.list, 'list query param sent when listing secrets');
      return { data: { keys: [] } };
    });

    this.store.query('kv/metadata', { backend: this.backend, pathToSecret: 'directory/' });
  });

  test('it should make request to correct endpoint on delete metadata', async function (assert) {
    assert.expect(3);
    const data = this.server.create('kv-metadatum');
    data.id = kvMetadataPath('kv-engine', 'my-secret');
    this.store.pushPayload('kv/metadata', {
      modelName: 'kv/metadata',
      ...data,
    });
    this.server.delete(kvMetadataPath('kv-engine', 'my-secret'), () => {
      assert.ok(true, 'request made to correct endpoint on delete metadata.');
    });

    let record = await this.store.peekRecord('kv/metadata', data.id);

    await record.destroyRecord();
    assert.true(record.isDestroyed, 'record is destroyed');
    record = await this.store.peekRecord('kv/metadata', this.id);
    assert.strictEqual(record, null, 'record is no longer in store');
  });
});
