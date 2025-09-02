/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Serializer | secret-engine', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    this.serializer = this.owner.lookup('serializer:secret-engine');
    this.path = 'kv-engine/';
    this.backend = {
      accessor: 'kv_77813cc8',
      config: {
        default_lease_ttl: 0,
        force_no_cache: false,
        max_lease_ttl: 0,
      },
      deprecation_status: 'supported',
      description: '',
      external_entropy_access: false,
      local: true,
      options: null,
      plugin_version: '',
      running_plugin_version: 'v0.16.1+builtin',
      running_sha256: '',
      seal_wrap: false,
      type: 'kv',
      uuid: '400a4673-6bd9-1336-b84c-caf43ee28340',
    };
  });

  test('it should not overwrite options for version 2', async function (assert) {
    assert.expect(1);
    this.backend.options = { version: '2' };
    const expectedData = {
      ...this.backend,
      id: 'kv-engine',
      path: 'kv-engine/',
      options: {
        version: '2',
      },
    };
    assert.propEqual(
      this.serializer.normalizeBackend(this.path, this.backend),
      expectedData,
      'options contain version 2'
    );
  });

  test('it should add version 1 for kv mounts when options is null', async function (assert) {
    assert.expect(1);

    const expectedData = {
      ...this.backend,
      id: 'kv-engine',
      path: 'kv-engine/',
      options: {
        version: '1',
      },
    };
    assert.propEqual(
      this.serializer.normalizeBackend(this.path, this.backend),
      expectedData,
      'options contains version 1'
    );
  });

  test('it should add version 1 for kv mounts if options has data but no version key', async function (assert) {
    assert.expect(1);

    this.backend.options = { foo: 'bar' };
    const expectedData = {
      ...this.backend,
      id: 'kv-engine',
      path: 'kv-engine/',
      options: {
        foo: 'bar',
        version: '1',
      },
    };

    assert.propEqual(
      this.serializer.normalizeBackend(this.path, this.backend),
      expectedData,
      'it adds version 1 to existing options'
    );
  });

  test('it should not update options for non-kv engines', async function (assert) {
    assert.expect(1);

    const cubbyholeData = {
      accessor: 'cubbyhole_8a89fbc7',
      config: {
        default_lease_ttl: 0,
        force_no_cache: false,
        max_lease_ttl: 0,
      },
      description: 'per-token private secret storage',
      external_entropy_access: false,
      local: true,
      options: null,
      plugin_version: '',
      running_plugin_version: 'v1.15.0+builtin.vault',
      running_sha256: '',
      seal_wrap: false,
      type: 'cubbyhole',
      uuid: 'a7638176-6c6e-2c65-0e50-05d689ef7fc8',
    };

    const expectedData = {
      ...cubbyholeData,
      id: 'cubbyhole',
      path: 'cubbyhole/',
    };
    assert.propEqual(
      this.serializer.normalizeBackend('cubbyhole/', cubbyholeData),
      expectedData,
      'options are still null'
    );
  });
});
