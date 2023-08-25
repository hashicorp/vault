/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { kvDataPath } from 'vault/utils/kv-path';
import { encodePath } from 'vault/utils/path-encoding-helpers';
import { Response } from 'miragejs';

const EXAMPLE_KV_DATA_CREATE_RESPONSE = {
  request_id: 'foobar',
  data: {
    created_time: '2023-06-21T16:18:31.479993Z',
    custom_metadata: null,
    deletion_time: '',
    destroyed: false,
    version: 1,
  },
};

const EXAMPLE_KV_DATA_GET_RESPONSE = {
  request_id: 'foobar',
  data: {
    data: { foo: 'bar' },
    metadata: {
      created_time: '2023-06-20T21:26:47.592306Z',
      custom_metadata: null,
      deletion_time: '',
      destroyed: false,
      version: 2,
    },
  },
};

const EXAMPLE_CONTROL_GROUP_RESPONSE = {
  data: null,
  wrap_info: {
    token: 'some-token',
    accessor: 'some-accessor',
    ttl: 86400,
    creation_time: '2023-08-09T16:08:06-05:00',
    creation_path: 'some/path/here',
  },
};

const EXAMPLE_KV_DATA_DESTROYED = {
  data: {
    data: null,
    metadata: {
      created_time: '2023-08-09T20:10:24.4825Z',
      custom_metadata: null,
      deletion_time: '',
      destroyed: true,
      version: 2,
    },
  },
};

const EXAMPLE_KV_DATA_DELETED = {
  data: {
    data: null,
    metadata: {
      created_time: '2023-08-09T20:10:24.571332Z',
      custom_metadata: null,
      deletion_time: '2023-08-09T20:10:24.70176Z',
      destroyed: false,
      version: 2,
    },
  },
};

module('Unit | Adapter | kv/data', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.version = this.owner.lookup('service:version');
    this.version.version = 'example+ent'; // Required for testing control-group flow
    this.secretMountPath = this.owner.lookup('service:secret-mount-path');
    this.backend = 'my/kv-back&end';
    this.secretMountPath.currentPath = this.backend;
    this.path = 'beep/bop/my secret';
    this.version = '2';
    this.id = kvDataPath(this.backend, this.path, this.version);
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
    this.endpoint = (noun) => `${encodePath(this.backend)}/${noun}/${encodePath(this.path)}`;
  });

  test('it should make request to correct endpoint on createRecord', async function (assert) {
    assert.expect(8);
    this.server.post(this.endpoint('data'), (schema, req) => {
      assert.ok('POST request made to correct endpoint when creating new record');
      const body = JSON.parse(req.requestBody);
      assert.deepEqual(body, {
        data: {
          foo: 'bar',
        },
        options: {
          cas: 0,
        },
      });
      return EXAMPLE_KV_DATA_CREATE_RESPONSE;
    });
    const record = this.store.createRecord('kv/data', {
      backend: this.backend,
      path: this.path,
      secretData: { foo: 'bar' },
      casVersion: 0,
    });
    await record.save();
    assert.strictEqual(record.path, this.path, 'record has correct path');
    assert.strictEqual(record.backend, this.backend, 'record has correct backend');
    assert.strictEqual(record.version, 1, 'record has correct version');
    assert.deepEqual(record.secretData, { foo: 'bar' }, 'record has correct data');
    assert.strictEqual(record.createdTime, '2023-06-21T16:18:31.479993Z', 'record has correct createdTime');
    assert.strictEqual(
      record.id,
      `${encodePath(this.backend)}/data/${encodePath(this.path)}?version=1`,
      'record has correct id'
    );
  });

  test('it should not send cas if casVersion is not a number', async function (assert) {
    assert.expect(8);
    this.server.post(this.endpoint('data'), (schema, req) => {
      assert.ok('POST request made to correct endpoint when creating new record');
      const body = JSON.parse(req.requestBody);
      assert.deepEqual(body, {
        data: {
          foo: 'bar',
        },
      });
      return EXAMPLE_KV_DATA_CREATE_RESPONSE;
    });
    const record = this.store.createRecord('kv/data', {
      backend: this.backend,
      path: this.path,
      secretData: { foo: 'bar' },
    });
    await record.save();
    assert.strictEqual(record.path, this.path, 'record has correct path');
    assert.strictEqual(record.backend, this.backend, 'record has correct backend');
    assert.strictEqual(record.version, 1, 'record has correct version');
    assert.deepEqual(record.secretData, { foo: 'bar' }, 'record has correct data');
    assert.strictEqual(record.createdTime, '2023-06-21T16:18:31.479993Z', 'record has correct createdTime');
    assert.strictEqual(
      record.id,
      `${encodePath(this.backend)}/data/${encodePath(this.path)}?version=1`,
      'record has correct id'
    );
  });

  test('it should make request to correct endpoint on queryRecord', async function (assert) {
    assert.expect(8);
    this.server.get(this.endpoint('data'), (schema, req) => {
      assert.ok(true, 'request is made to correct url on queryRecord.');
      assert.strictEqual(
        req.queryParams.version,
        this.version,
        'request includes the version flag on queryRecord.'
      );
      return EXAMPLE_KV_DATA_GET_RESPONSE;
    });

    const record = await this.store.queryRecord('kv/data', this.payload);
    assert.strictEqual(record.path, this.path, 'record has correct path');
    assert.strictEqual(record.backend, this.backend, 'record has correct backend');
    assert.strictEqual(record.version, 2, 'record has correct version');
    assert.deepEqual(record.secretData, { foo: 'bar' }, 'record has correct data');
    assert.strictEqual(record.createdTime, '2023-06-20T21:26:47.592306Z', 'record has correct createdTime');
    assert.strictEqual(
      record.id,
      `${encodePath(this.backend)}/data/${encodePath(this.path)}?version=${this.version}`,
      'record has correct id'
    );
  });

  test('it should handle a 404 not found response properly', async function (assert) {
    assert.expect(1);
    this.server.get(this.endpoint('data'), () => {
      // This is what the API currently returns for not found
      return new Response(404, {}, { errors: [] });
    });

    try {
      await this.store.queryRecord('kv/data', this.payload);
    } catch (e) {
      assert.ok('throws the error');
    }
  });

  test('it should handle a 403 permission denied properly', async function (assert) {
    assert.expect(8);
    this.server.get(this.endpoint('data'), (schema, req) => {
      assert.ok(true, 'request is made to correct url on queryRecord.');
      assert.strictEqual(
        req.queryParams.version,
        this.version,
        'request includes the version flag on queryRecord.'
      );
      return new Response(403, {}, { errors: ['1 error occurred:\n\t* permission denied\n\n'] });
    });

    const record = await this.store.queryRecord('kv/data', this.payload);
    assert.strictEqual(record.path, this.path, 'record has correct path');
    assert.strictEqual(record.backend, this.backend, 'record has correct backend');
    assert.strictEqual(record.version, 2, 'record has version based on request');
    assert.strictEqual(record.secretData, undefined, 'record does not include data');
    assert.strictEqual(record.failReadErrorCode, 403, 'record has error response recorded');
    assert.strictEqual(
      record.id,
      `${encodePath(this.backend)}/data/${encodePath(this.path)}?version=${this.version}`,
      'record has correct id'
    );
  });

  test('it should handle a soft-deleted version properly', async function (assert) {
    this.server.get(this.endpoint('data'), () => {
      return new Response(404, {}, EXAMPLE_KV_DATA_DELETED);
    });

    const record = await this.store.queryRecord('kv/data', this.payload);
    assert.strictEqual(record.path, this.path, 'record has correct path');
    assert.strictEqual(record.backend, this.backend, 'record has correct backend');
    assert.strictEqual(record.version, 2, 'record has version based on request');
    assert.strictEqual(record.deletionTime, '2023-08-09T20:10:24.70176Z', 'record includes deletion time');
    assert.strictEqual(record.failReadErrorCode, undefined, 'record does not have failed error code');
    assert.strictEqual(
      record.id,
      `${encodePath(this.backend)}/data/${encodePath(this.path)}?version=${this.version}`,
      'record has correct id'
    );
  });

  test('it should handle a destroyed version properly', async function (assert) {
    this.server.get(this.endpoint('data'), () => {
      return new Response(404, {}, EXAMPLE_KV_DATA_DESTROYED);
    });

    const record = await this.store.queryRecord('kv/data', this.payload);
    assert.strictEqual(record.path, this.path, 'record has correct path');
    assert.strictEqual(record.backend, this.backend, 'record has correct backend');
    assert.strictEqual(record.version, 2, 'record has version based on request');
    assert.true(record.destroyed, 'record has destroyed value');
    assert.strictEqual(record.failReadErrorCode, undefined, 'record does not have error code');
    assert.strictEqual(
      record.id,
      `${encodePath(this.backend)}/data/${encodePath(this.path)}?version=${this.version}`,
      'record has correct id'
    );
  });

  test('it should handle a control group response properly', async function (assert) {
    assert.expect(1);
    this.server.get(this.endpoint('data'), () => {
      return EXAMPLE_CONTROL_GROUP_RESPONSE;
    });

    try {
      await this.store.queryRecord('kv/data', this.payload);
    } catch (e) {
      assert.ok('throws the error');
    }
  });

  test('it should make request to correct endpoint on delete latest version', async function (assert) {
    assert.expect(3);
    this.server.delete(this.endpoint('data'), () => {
      assert.ok(true, 'request made to correct endpoint on delete latest version.');
      return new Response(204);
    });

    this.store.pushPayload('kv/data', {
      modelName: 'kv/data',
      id: this.id,
      ...this.payload,
    });
    let record = await this.store.peekRecord('kv/data', this.id);
    await record.destroyRecord({ adapterOptions: { deleteType: 'delete-latest-version' } });
    assert.true(record.isDeleted, 'record is deleted');
    record = await this.store.peekRecord('kv/data', this.id);
    assert.strictEqual(record, null, 'record is no longer in store');
  });

  test('it should make request to correct endpoint on delete specific versions', async function (assert) {
    assert.expect(4);
    this.server.post(this.endpoint('delete'), (schema, req) => {
      const { versions } = JSON.parse(req.requestBody);
      assert.strictEqual(versions, 2, 'version array is sent in the payload.');
      assert.ok(true, 'request made to correct endpoint on delete specific version.');
    });

    this.store.pushPayload('kv/data', {
      modelName: 'kv/data',
      id: this.id,
      ...this.payload,
    });
    let record = await this.store.peekRecord('kv/data', this.id);
    await record.destroyRecord({
      adapterOptions: { deleteType: 'delete-version', deleteVersions: 2 },
    });
    assert.true(record.isDeleted, 'record is deleted');
    record = await this.store.peekRecord('kv/data', this.id);
    assert.strictEqual(record, null, 'record is no longer in store');
  });

  test('it should make request to correct endpoint on undelete', async function (assert) {
    assert.expect(4);
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
    let record = await this.store.peekRecord('kv/data', this.id);

    await record.destroyRecord({
      adapterOptions: { deleteType: 'undelete', deleteVersions: 2 },
    });
    assert.true(record.isDeleted, 'record is deleted');
    record = await this.store.peekRecord('kv/data', this.id);
    assert.strictEqual(record, null, 'record is no longer in store');
  });

  test('it should make request to correct endpoint on destroy specific versions', async function (assert) {
    assert.expect(4);
    this.server.put(`${encodePath(this.backend)}/destroy/${encodePath(this.path)}`, (schema, req) => {
      const { versions } = JSON.parse(req.requestBody);
      assert.strictEqual(versions, 2, 'version array is sent in the payload.');
      assert.ok(true, 'request made to correct endpoint on destroy specific version.');
    });

    this.store.pushPayload('kv/data', {
      modelName: 'kv/data',
      id: this.id,
      ...this.payload,
    });
    let record = await this.store.peekRecord('kv/data', this.id);
    await record.destroyRecord({
      adapterOptions: { deleteType: 'destroy', deleteVersions: 2 },
    });
    assert.true(record.isDeleted, 'record is deleted');
    record = await this.store.peekRecord('kv/data', this.id);
    assert.strictEqual(record, null, 'record is no longer in store');
  });
});
