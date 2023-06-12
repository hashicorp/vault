/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';

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
      assert.strictEqual(req.queryParams.version, '2', 'request includes the version flag on query.');
      assert.ok(true, 'request is made to correct url on query.');
    });

    this.store.queryRecord('kv/data', { backend: this.backend, path: this.path, version: 2 });
    // unsure why I"m getting the failure.
  });
});
