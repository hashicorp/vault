/**
 * Copyright IBM Corp. 2016, 2026
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

  test('normalization retains plain snake-case config without creating embedded records', function (assert) {
    const config = {
      ...this.backend.config,
      allowed_managed_keys: ['managed-key'],
      identity_token_key: 'signing-key',
      plugin_version: 'v1.0.0',
    };
    const store = this.owner.lookup('service:store');
    const normalized = this.serializer.normalize(store.modelFor('secret-engine'), {
      ...this.backend,
      id: 'plain-config',
      config,
    });

    assert.deepEqual(normalized.data.attributes.config, config, 'config fields retain their API casing');
    assert.notOk(normalized.data.relationships?.config, 'config is not an Ember Data relationship');
    assert.notOk('id' in config, 'normalization does not add an embedded-record ID');
  });

  test('serialization does not mutate plain config or share default config between records', function (assert) {
    const store = this.owner.lookup('service:store');
    const config = {
      default_lease_ttl: '2h',
      max_lease_ttl: '4h',
      identity_token_key: 'wif-only-key',
      audit_non_hmac_request_keys: ['request-key'],
      allowed_response_headers: ['X-Custom'],
    };
    const record = store.createRecord('secret-engine', { type: 'cubbyhole', config });
    const serialized = record.serialize();
    const { identity_token_key, ...expected } = config;

    assert.deepEqual(serialized.config, expected, 'only the non-WIF key is omitted');
    assert.strictEqual(config.identity_token_key, identity_token_key, 'record config is not mutated');
    const first = store.createRecord('secret-engine');
    const second = store.createRecord('secret-engine');
    first.config.default_lease_ttl = '1h';
    assert.deepEqual(second.config, {}, 'new records have independent config objects');
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
