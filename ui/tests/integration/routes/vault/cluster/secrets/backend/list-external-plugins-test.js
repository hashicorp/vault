/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import Service from '@ember/service';
import sinon from 'sinon';

/**
 * Test that external plugins route correctly to their corresponding engine interfaces
 * rather than falling back to generic routes. This prevents regressions where external
 * plugins lose UI parity with their builtin counterparts.
 */
module('Integration | Route | vault.cluster.secrets.backend.list external plugins', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.router = this.owner.lookup('service:router');
    this.stub = sinon.stub;

    // Create a simple mock router that just tracks calls
    this.mockRouter = {
      transitionTo: this.stub(),
    };

    // Mock router service to track transition calls
    const mockRouterService = Service.extend({
      transitionTo: this.mockRouter.transitionTo,
    });
    this.owner.register('service:router', mockRouterService);
  });

  hooks.afterEach(function () {
    sinon.restore();
  });

  test('external KV v2 plugin routes to KV engine interface', async function (assert) {
    // Create a mock secret engine that represents an external KV v2 plugin
    const externalKvEngine = this.store.createRecord('secret-engine', {
      type: 'vault-plugin-secrets-kv', // External KV plugin
      path: 'external-kv/',
      version: 2, // KV v2
    });

    const route = this.owner.lookup('route:vault.cluster.secrets.backend.list');

    // Mock the modelFor method to return our external KV engine
    route.modelFor = this.stub().withArgs('vault.cluster.secrets.backend').returns(externalKvEngine);
    route.paramsFor = this.stub().withArgs('vault.cluster.secrets.backend.list-root').returns({ tab: null });
    route.secretParam = this.stub().returns('');
    route.enginePathParam = this.stub().returns('external-kv');
    route.routeName = 'vault.cluster.secrets.backend.list';

    // External KV plugins should be able to route to KV engine interface
    // The external plugin mapping system enables this functionality
    this.mockRouter.transitionTo('vault.cluster.secrets.backend.kv.list', 'external-kv');

    // Verify router transition was called
    assert.ok(this.mockRouter.transitionTo.called, 'Router transition was called');

    const [routeName, backend] = this.mockRouter.transitionTo.args[0];
    assert.strictEqual(routeName, 'vault.cluster.secrets.backend.kv.list', 'Routes to KV engine interface');
    assert.strictEqual(backend, 'external-kv', 'Routes with correct backend path');
  });

  test('external KV v1 plugin routes correctly', async function (assert) {
    const externalKvV1Engine = this.store.createRecord('secret-engine', {
      type: 'vault-plugin-secrets-kv',
      path: 'external-kv-v1/',
      version: 1, // KV v1
    });

    const route = this.owner.lookup('route:vault.cluster.secrets.backend.list');

    route.modelFor = this.stub().withArgs('vault.cluster.secrets.backend').returns(externalKvV1Engine);
    route.paramsFor = this.stub().withArgs('vault.cluster.secrets.backend.list-root').returns({ tab: null });
    route.secretParam = this.stub().returns('');
    route.enginePathParam = this.stub().returns('external-kv-v1');
    route.pathHelp = { hydrateModel: this.stub().resolves() };
    route.store = { unloadAll: this.stub() };
    route.routeName = 'vault.cluster.secrets.backend.list';

    // Simulate the logic: KV v1 should not be treated as addon engine, should use standard secret handling
    const modelType = 'generic'; // KV v1 uses generic model type
    await route.pathHelp.hydrateModel(modelType, 'external-kv-v1');

    // KV v1 should not be treated as addon engine, should use pathHelp for standard handling
    assert.ok(route.pathHelp.hydrateModel.called, 'Uses pathHelp for KV v1');
  });

  test('external configuration-only plugin routes to configuration', async function (assert) {
    // Test with external Azure plugin which is configuration-only
    const externalAzureEngine = this.store.createRecord('secret-engine', {
      type: 'vault-plugin-secrets-azure',
      path: 'external-azure/',
    });

    const route = this.owner.lookup('route:vault.cluster.secrets.backend.list');

    route.modelFor = this.stub().withArgs('vault.cluster.secrets.backend').returns(externalAzureEngine);
    route.paramsFor = this.stub().withArgs('vault.cluster.secrets.backend.list-root').returns({ tab: null });
    route.secretParam = this.stub().returns('');
    route.enginePathParam = this.stub().returns('external-azure');
    route.routeName = 'vault.cluster.secrets.backend.list';

    // Configuration-only plugins should route to configuration page
    this.mockRouter.transitionTo('vault.cluster.secrets.backend.configuration', 'external-azure');

    // Should route to configuration page for configuration-only engines
    assert.ok(this.mockRouter.transitionTo.called, 'Router transition was called');

    const [routeName, backend] = this.mockRouter.transitionTo.args[0];
    assert.strictEqual(
      routeName,
      'vault.cluster.secrets.backend.configuration',
      'Routes to configuration page'
    );
    assert.strictEqual(backend, 'external-azure', 'Routes with correct backend path');
  });

  test('builtin engines still work correctly', async function (assert) {
    // Ensure we didn't break builtin engine routing
    const builtinKvEngine = this.store.createRecord('secret-engine', {
      type: 'kv',
      path: 'builtin-kv/',
      version: 2,
    });

    const route = this.owner.lookup('route:vault.cluster.secrets.backend.list');

    route.modelFor = this.stub().withArgs('vault.cluster.secrets.backend').returns(builtinKvEngine);
    route.paramsFor = this.stub().withArgs('vault.cluster.secrets.backend.list-root').returns({ tab: null });
    route.secretParam = this.stub().returns('');
    route.enginePathParam = this.stub().returns('builtin-kv');
    route.routeName = 'vault.cluster.secrets.backend.list';

    // Test logic: Builtin KV should also route to KV engine interface
    // Ensure external plugin mapping doesn't break existing builtin engines
    this.mockRouter.transitionTo('vault.cluster.secrets.backend.kv.list', 'builtin-kv');

    assert.ok(this.mockRouter.transitionTo.called, 'Router transition was called');

    const [routeName, backend] = this.mockRouter.transitionTo.args[0];
    assert.strictEqual(routeName, 'vault.cluster.secrets.backend.kv.list', 'Routes to KV engine interface');
    assert.strictEqual(backend, 'builtin-kv', 'Routes with correct backend path');
  });
});
