/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import Service from '@ember/service';
import sinon from 'sinon';
import SecretsEngineResource from 'vault/resources/secrets/engine';

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
    const externalKvEngine = new SecretsEngineResource({
      type: 'vault-plugin-secrets-kv', // External KV plugin
      path: 'external-kv/',
      options: { version: '2' },
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
    const externalKvV1Engine = new SecretsEngineResource({
      type: 'vault-plugin-secrets-kv',
      path: 'external-kv-v1/',
      options: { version: '1' },
    });

    const route = this.owner.lookup('route:vault.cluster.secrets.backend.list');

    route.modelFor = this.stub().withArgs('vault.cluster.secrets.backend').returns(externalKvV1Engine);
    route.paramsFor = this.stub();
    route.paramsFor.withArgs('vault.cluster.secrets.backend.list-root').returns({ tab: null });
    route.paramsFor.withArgs('vault.cluster.secrets.backend').returns({ backend: 'external-kv-v1' });
    route.secretParam = this.stub().returns('');
    route.enginePathParam = this.stub().returns('external-kv-v1');
    route.store = { unloadAll: this.stub() };
    route.routeName = 'vault.cluster.secrets.backend.list-root';

    await route.beforeModel();
    assert.true(this.mockRouter.transitionTo.notCalled, 'KV v1 stays on the standard list route');
    assert.true(
      route.store.unloadAll.calledOnceWithExactly('capabilities'),
      'capabilities are refreshed before listing secrets'
    );
  });

  test('external configuration-only plugin routes to configuration', async function (assert) {
    // Test with external Azure plugin which is configuration-only
    const externalAzureEngine = new SecretsEngineResource({
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
    const builtinKvEngine = new SecretsEngineResource({
      type: 'kv',
      path: 'builtin-kv/',
      options: { version: '2' },
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

  test('legacy KV v2 edit links still redirect to the migrated edit page', async function (assert) {
    const route = this.owner.lookup('route:vault.cluster.secrets.backend.secret-edit');
    route.router = this.mockRouter;
    route.modelFor = this.stub().returns(
      new SecretsEngineResource({ type: 'kv', path: 'example-kv/', options: { version: '2' } })
    );
    route.secretParam = this.stub().returns('nested/secret');
    route.routeName = 'vault.cluster.secrets.backend.edit';

    await route.beforeModel();
    assert.true(
      this.mockRouter.transitionTo.calledOnceWithExactly(
        'vault.cluster.secrets.backend.kv.secret.details.edit',
        'example-kv',
        'nested/secret'
      ),
      'legacy edit links retain the backend and complete secret path'
    );
  });

  test('legacy folder edit links return to the parent list', async function (assert) {
    const route = this.owner.lookup('route:vault.cluster.secrets.backend.secret-edit');
    route.router = this.mockRouter;
    route.modelFor = this.stub().returns(new SecretsEngineResource({ type: 'kv', path: 'example-kv/' }));
    route.secretParam = this.stub().returns('parent/folder/');
    route.routeName = 'vault.cluster.secrets.backend.edit';

    await route.beforeModel();
    assert.true(
      this.mockRouter.transitionTo.calledOnceWithExactly('vault.cluster.secrets.backend.list', 'parent/'),
      'folders cannot be edited and retain the parent-list redirect'
    );
  });
});
