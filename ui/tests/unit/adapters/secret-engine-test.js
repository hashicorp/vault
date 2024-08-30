/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';

module('Unit | Adapter | secret engine', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  const storeStub = {
    serializerFor() {
      return {
        serializeIntoHash() {},
      };
    },
  };
  const type = {
    modelName: 'secret-engine',
  };

  test('Empty query', function (assert) {
    assert.expect(1);
    this.server.get('/sys/internal/ui/mounts', () => {
      assert.ok('query calls the correct url');
      return {};
    });
    const adapter = this.owner.lookup('adapter:secret-engine');
    adapter['query'](storeStub, type, {});
  });
  test('Query with a path', function (assert) {
    assert.expect(1);
    this.server.get('/sys/internal/ui/mounts/foo', () => {
      assert.ok('query calls the correct url');
      return {};
    });
    const adapter = this.owner.lookup('adapter:secret-engine');
    adapter['query'](storeStub, type, { path: 'foo' });
  });

  test('Query with nested path', function (assert) {
    assert.expect(1);
    this.server.get('/sys/internal/ui/mounts/foo/bar/baz', () => {
      assert.ok('query calls the correct url');
      return {};
    });
    const adapter = this.owner.lookup('adapter:secret-engine');
    adapter['query'](storeStub, type, { path: 'foo/bar/baz' });
  });

  module('WIF secret engines', function (hooks) {
    hooks.beforeEach(function () {
      this.store = this.owner.lookup('service:store');
    });

    test('it should make request to correct endpoint when creating new record', async function (assert) {
      assert.expect(1);
      this.server.post('/sys/mounts/aws-wif', (schema, req) => {
        assert.deepEqual(
          JSON.parse(req.requestBody),
          {
            path: 'aws-wif',
            type: 'aws',
            config: { id: 'aws-wif', identity_token_key: 'test-key', listing_visibility: 'hidden' },
          },
          'Correct payload is sent when adding aws secret engine with identity_token_key set'
        );
        return {};
      });
      const mountData = {
        id: 'aws-wif',
        path: 'aws-wif',
        type: 'aws',
        config: this.store.createRecord('mount-config', {
          identityTokenKey: 'test-key',
        }),
        uuid: 'f1739f9d-dfc0-83c8-011f-ec17103a06c2',
      };
      const record = this.store.createRecord('secret-engine', mountData);
      await record.save();
    });

    test('it should not send identity_token_key if not set', async function (assert) {
      assert.expect(1);
      this.server.post('/sys/mounts/aws-wif', (schema, req) => {
        assert.deepEqual(
          JSON.parse(req.requestBody),
          {
            path: 'aws-wif',
            type: 'aws',
            config: { id: 'aws-wif', max_lease_ttl: '125h', listing_visibility: 'hidden' },
          },
          'Correct payload is sent when adding aws secret engine with no identity_token_key set'
        );
        return {};
      });
      const mountData = {
        id: 'aws-wif',
        path: 'aws-wif',
        type: 'aws',
        config: this.store.createRecord('mount-config', {
          maxLeaseTtl: '125h',
        }),
        uuid: 'f1739f9d-dfc0-83c8-011f-ec17103a06c2',
      };
      const record = this.store.createRecord('secret-engine', mountData);
      await record.save();
    });

    test('it should not send identity_token_key if set on a non-WIF secret engine', async function (assert) {
      assert.expect(1);
      this.server.post('/sys/mounts/cubbyhole-test', (schema, req) => {
        assert.deepEqual(
          JSON.parse(req.requestBody),
          {
            path: 'cubbyhole-test',
            type: 'cubbyhole',
            config: { id: 'cubbyhole-test', max_lease_ttl: '125h', listing_visibility: 'hidden' },
          },
          'Correct payload is sent when sending a non-wif secret engine with identity_token_key accidentally set'
        );
        return {};
      });
      const mountData = {
        id: 'cubbyhole-test',
        path: 'cubbyhole-test',
        type: 'cubbyhole',
        config: this.store.createRecord('mount-config', {
          maxLeaseTtl: '125h',
          identity_token_key: 'test-key',
        }),
        uuid: 'f1739f9d-dfc0-83c8-011f-ec17103a06c4',
      };
      const record = this.store.createRecord('secret-engine', mountData);
      await record.save();
    });
  });
});
