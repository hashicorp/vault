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

  test('Fails gracefully finding records for non ssh engines', function (assert) {
    assert.expect(1);
    const snapshot = {
      attr() {
        return { type: 'aws', path: 'aws/' };
      },
    };
    const adapter = this.owner.lookup('adapter:secret-engine');
    const response = adapter.findRecord(storeStub, 'aws', { path: 'aws' }, snapshot);
    assert.propEqual(response, { data: {} }, 'returns empty data object');
  });
});
