/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { kvDataPath } from 'vault/utils/kv-path';

module('Unit | Adapter | kv/metadata', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.secretMountPath = this.owner.lookup('service:secret-mount-path');
    this.backend = 'kv-backend';
    this.secretMountPath.currentPath = this.backend;
    this.path = 'beep/bop/my-secret';
    this.id = kvDataPath(this.backend, this.path, 'metadata');
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
  });

  hooks.afterEach(function () {
    this.store.unloadAll('kv/metadata');
    this.server.shutdown();
  });

  test('it should make request to correct endpoint on createRecord', async function (assert) {
    assert.expect(1);
    this.server.post(`${this.backend}/metadata/${this.path}`, () => {
      assert.ok('POST request made to correct endpoint when creating new record');
    });
    const record = this.store.createRecord('kv/metadata', { backend: this.backend, path: this.path });
    await record.save();
  });

  test('it should make request to correct endpoint on queryRecord', async function (assert) {
    assert.expect(1);
    this.server.get(`${this.backend}/metadata/${this.path}`, () => {
      assert.ok(true, 'request is made to correct url on queryRecord.');
    });

    await this.store.queryRecord('kv/metadata', { backend: this.backend, path: this.path });
  });
});
