/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import Route from '@ember/routing/route';
import { getBackendEffectiveType, getEnginePathParam } from 'vault/utils/backend-route-helpers';
import sinon from 'sinon';

/**
 * Test the backend route helper utilities to ensure external plugin mapping
 * works correctly.
 */
module('Unit | Utility | backend-route-helpers', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    // Create a test route
    this.owner.register('route:test', Route);
    this.route = this.owner.lookup('route:test');
    this.stub = sinon.stub;
  });

  hooks.afterEach(function () {
    sinon.restore();
  });

  module('getBackendEffectiveType', function () {
    test('returns effective type for external keymgmt plugins', function (assert) {
      // Mock modelFor to return an external keymgmt engine
      this.route.modelFor = this.stub().withArgs('vault.cluster.secrets.backend').returns({
        engineType: 'vault-plugin-secrets-keymgmt',
      });

      const effectiveType = getBackendEffectiveType(this.route);

      assert.strictEqual(effectiveType, 'keymgmt', 'External keymgmt plugin returns effective type keymgmt');
    });

    test('returns effective type for external KV plugins', function (assert) {
      this.route.modelFor = this.stub().withArgs('vault.cluster.secrets.backend').returns({
        engineType: 'vault-plugin-secrets-kv',
      });

      const effectiveType = getBackendEffectiveType(this.route);

      assert.strictEqual(effectiveType, 'kv', 'External KV plugin returns effective type kv');
    });

    test('returns original type for builtin engines', function (assert) {
      this.route.modelFor = this.stub().withArgs('vault.cluster.secrets.backend').returns({
        engineType: 'keymgmt',
      });

      const effectiveType = getBackendEffectiveType(this.route);

      assert.strictEqual(effectiveType, 'keymgmt', 'Builtin keymgmt returns original type');
    });

    test('returns original type for unknown external plugins', function (assert) {
      this.route.modelFor = this.stub().withArgs('vault.cluster.secrets.backend').returns({
        engineType: 'vault-plugin-secrets-unknown',
      });

      const effectiveType = getBackendEffectiveType(this.route);

      assert.strictEqual(
        effectiveType,
        'vault-plugin-secrets-unknown',
        'Unknown external plugin returns original type'
      );
    });

    test('handles external Azure plugins', function (assert) {
      this.route.modelFor = this.stub().withArgs('vault.cluster.secrets.backend').returns({
        engineType: 'vault-plugin-secrets-azure',
      });

      const effectiveType = getBackendEffectiveType(this.route);

      assert.strictEqual(effectiveType, 'azure', 'External Azure plugin returns effective type azure');
    });
  });

  module('getEnginePathParam', function () {
    test('returns backend parameter from route params', function (assert) {
      this.route.paramsFor = this.stub().withArgs('vault.cluster.secrets.backend').returns({
        backend: 'external-keymgmt',
      });

      const enginePath = getEnginePathParam(this.route);

      assert.strictEqual(enginePath, 'external-keymgmt', 'Returns backend parameter from route');
    });

    test('handles different backend paths', function (assert) {
      this.route.paramsFor = this.stub().withArgs('vault.cluster.secrets.backend').returns({
        backend: 'my-custom-engine-path',
      });

      const enginePath = getEnginePathParam(this.route);

      assert.strictEqual(enginePath, 'my-custom-engine-path', 'Returns custom backend path');
    });
  });

  module('integration with route operations', function () {
    test('utility functions can be used together in route logic', function (assert) {
      // Mock both functions used in typical route scenarios
      this.route.modelFor = this.stub().withArgs('vault.cluster.secrets.backend').returns({
        engineType: 'vault-plugin-secrets-keymgmt',
      });

      this.route.paramsFor = this.stub().withArgs('vault.cluster.secrets.backend').returns({
        backend: 'external-keymgmt',
      });

      const effectiveType = getBackendEffectiveType(this.route);
      const enginePath = getEnginePathParam(this.route);

      assert.strictEqual(effectiveType, 'keymgmt', 'Gets effective type correctly');
      assert.strictEqual(enginePath, 'external-keymgmt', 'Gets engine path correctly');
    });
  });
});
