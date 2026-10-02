/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import Service from '@ember/service';
import KvEngine from 'kv/engine';
import KvSecretsListRoute from 'kv/routes/list-directory';

module('Unit | Engine | kv', function (hooks) {
  setupTest(hooks);

  test('it exposes the KV mount retry service to engine routes', function (assert) {
    const engine = KvEngine.create();

    assert.true(
      engine.dependencies.services.includes('kv-mount-retry'),
      'the KV engine can inject the host retry service'
    );
  });

  test('it injects the named KV mount retry service into the list route', function (assert) {
    class KvMountRetryStub extends Service {
      listKvSecrets = {};
    }

    this.owner.register('service:app-router', Service);
    this.owner.register('service:secret-mount-path', Service);
    this.owner.register('service:api', Service);
    this.owner.register('service:capabilities', Service);
    this.owner.register('service:kv-mount-retry', KvMountRetryStub);
    this.owner.register('route:kv-list-directory', KvSecretsListRoute);

    const kvMountRetry = this.owner.lookup('service:kv-mount-retry');
    const route = this.owner.lookup('route:kv-list-directory');

    assert.strictEqual(route.kvMountRetry, kvMountRetry, 'the list route resolves the host service by name');
  });
});
