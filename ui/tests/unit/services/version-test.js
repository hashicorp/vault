/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { FEATURE_VERSIONS } from 'vault/utils/feature-versions';
import Sinon from 'sinon';

module('Unit | Service | version', function (hooks) {
  setupTest(hooks);

  test('setting type computes isCommunity properly', function (assert) {
    const service = this.owner.lookup('service:version');
    service.type = 'community';
    assert.true(service.isCommunity);
    assert.false(service.isEnterprise);
  });

  test('setting type computes isEnterprise properly', function (assert) {
    const service = this.owner.lookup('service:version');
    service.type = 'enterprise';
    assert.false(service.isCommunity);
    assert.true(service.isEnterprise);
  });

  test('calculates versionDisplay correctly', function (assert) {
    const service = this.owner.lookup('service:version');
    service.type = 'community';
    service.version = '1.2.3';
    assert.strictEqual(service.versionDisplay, 'v1.2.3');
    service.type = 'enterprise';
    service.version = '1.4.7+ent';
    assert.strictEqual(service.versionDisplay, 'v1.4.7');
  });

  test('fetchVersion stores seal status data and only fetches it once', async function (assert) {
    const service = this.owner.lookup('service:version');
    const api = service.api;
    const response = { version: '1.2.3+ent', cluster_name: 'test-cluster' };
    const sealStatus = Sinon.stub(api.sys, 'sealStatusRaw').resolves({
      raw: new Response(JSON.stringify(response)),
    });

    await service.fetchVersion();

    assert.strictEqual(service.version, response.version, 'stores the Vault version');
    assert.strictEqual(service.clusterName, response.cluster_name, 'stores the cluster name');
    assert.deepEqual(
      sealStatus.firstCall.args[0],
      { headers: api.buildHeaders({ namespace: '' }).headers },
      'requests seal status in the root namespace with the current token'
    );

    await service.fetchVersion();

    assert.true(sealStatus.calledOnce, 'does not fetch seal status again after the version is loaded');
  });

  test('fetchType derives and caches the enterprise type from health status', async function (assert) {
    const service = this.owner.lookup('service:version');
    const api = service.api;
    const health = Sinon.stub(api.sys, 'readHealthStatus').resolves({ enterprise: true });

    await service.fetchType();

    assert.strictEqual(service.type, 'enterprise', 'uses the enterprise value from health status');
    assert.true(service.isEnterprise, 'exposes the enterprise classification');
    const initOverrides = health.firstCall.args[0];
    const requestContext = { context: { query: {} }, init: {} };
    const init = await initOverrides(requestContext);
    assert.deepEqual(
      requestContext.context.query,
      {
        standbycode: 200,
        sealedcode: 200,
        uninitcode: 200,
        drsecondarycode: 200,
        performancestandbycode: 200,
      },
      'preserves health status query parameters'
    );
    assert.deepEqual(
      init.headers,
      api.buildHeaders({ token: '', namespace: '' }).headers,
      'reads health status without a token in the root namespace'
    );

    await service.fetchType();

    assert.true(health.calledOnce, 'does not fetch health status again after the type is loaded');
  });

  test('fetchType derives the community type from health status', async function (assert) {
    const service = this.owner.lookup('service:version');
    Sinon.stub(service.api.sys, 'readHealthStatus').resolves({ enterprise: false });

    await service.fetchType();

    assert.strictEqual(service.type, 'community', 'uses the community value from health status');
    assert.true(service.isCommunity, 'exposes the community classification');
  });

  test('fetchType treats a chroot namespace health response as enterprise', async function (assert) {
    const service = this.owner.lookup('service:version');
    Sinon.stub(service.api.sys, 'readHealthStatus').resolves({ has_chroot_namespace: true });

    await service.fetchType();

    assert.strictEqual(
      service.type,
      'enterprise',
      'preserves enterprise classification for chroot namespaces'
    );
  });

  test('fetchType treats a failed health request as enterprise for chroot namespaces', async function (assert) {
    const service = this.owner.lookup('service:version');
    Sinon.stub(service.api.sys, 'readHealthStatus').rejects(new Error('health unavailable'));

    await service.fetchType();

    assert.strictEqual(service.type, 'enterprise', 'preserves the chroot namespace fallback');
  });

  test('fetchFeatures stores and caches license features', async function (assert) {
    const service = this.owner.lookup('service:version');
    const api = service.api;
    const features = ['Namespaces', 'Control Groups'];
    const fetchFeatures = Sinon.stub(api.request, 'get').resolves(new Response(JSON.stringify({ features })));
    service.type = 'enterprise';

    await service.fetchFeatures();

    assert.deepEqual(service.features, features, 'stores the features returned by the backend');
    assert.true(service.hasNamespaces, 'derives namespace availability from the returned features');
    assert.true(service.hasControlGroups, 'derives control group availability from the returned features');
    assert.strictEqual(
      fetchFeatures.firstCall.args[0],
      'sys/license/features',
      'reads license features from the root endpoint'
    );
    assert.deepEqual(
      fetchFeatures.firstCall.args[1],
      api.buildHeaders({ token: '', namespace: '' }).headers,
      'reads license features without a token in the root namespace'
    );

    await service.fetchFeatures();

    assert.true(fetchFeatures.calledOnce, 'does not fetch features again after they are loaded');
  });

  test('fetchFeatures does not request license features for community Vault', async function (assert) {
    const service = this.owner.lookup('service:version');
    const fetchFeatures = Sinon.stub(service.api.request, 'get');
    service.type = 'community';

    await service.fetchFeatures();

    assert.true(fetchFeatures.notCalled, 'skips the license feature request for community Vault');
  });

  test('fetchFeatures keeps the route usable when the license request fails', async function (assert) {
    const service = this.owner.lookup('service:version');
    Sinon.stub(service.api.request, 'get').rejects(new Error('license features unavailable'));
    service.type = 'enterprise';

    await service.fetchFeatures();

    assert.deepEqual(service.features, [], 'keeps an empty feature list when the request fails');
  });

  test('hasPerfReplication', function (assert) {
    const service = this.owner.lookup('service:version');
    assert.false(service.hasPerfReplication);
    service.features = ['Performance Replication'];
    assert.true(service.hasPerfReplication);
  });

  test('hasDRReplication', function (assert) {
    const service = this.owner.lookup('service:version');
    assert.false(service.hasDRReplication);
    service.features = ['DR Replication'];
    assert.true(service.hasDRReplication);
  });

  test('hasPKIOnly', function (assert) {
    const service = this.owner.lookup('service:version');
    assert.false(service.hasPKIOnly);
    service.features = ['PKI-only Secrets'];
    assert.true(service.hasPKIOnly);
  });

  // SHOW SECRETS SYNC TESTS
  test('hasSecretsSync: it returns false when version is community', function (assert) {
    const service = this.owner.lookup('service:version');
    service.type = 'community';
    assert.false(service.hasSecretsSync);
  });

  test('hasSecretsSync: it returns true when HVD managed', function (assert) {
    this.owner.lookup('service:flags').featureFlags = ['VAULT_CLOUD_ADMIN_NAMESPACE'];
    const service = this.owner.lookup('service:version');
    service.type = 'enterprise';
    assert.true(service.hasSecretsSync);
  });

  test('hasSecretsSync: it returns false when not on enterprise license', function (assert) {
    const service = this.owner.lookup('service:version');
    service.type = 'enterprise';
    service.features = ['replication'];
    assert.false(service.hasSecretsSync);
  });
  test('hasSecretsSync: it returns true when  on enterprise license', function (assert) {
    const service = this.owner.lookup('service:version');
    service.type = 'enterprise';
    service.features = ['secrets-sync'];
    assert.false(service.hasSecretsSync);
  });

  module('hasFeature', function (hooks) {
    let service;
    const TEST_KEY = 'test-version-service-feature';

    hooks.beforeEach(function () {
      service = this.owner.lookup('service:version');
      // Register a known test entry so tests don't depend on real FEATURE_VERSIONS content
      FEATURE_VERSIONS.push({ key: TEST_KEY, name: 'Test Feature', version: '2.0.0' });
    });

    hooks.afterEach(function () {
      const idx = FEATURE_VERSIONS.findIndex((f) => f.key === TEST_KEY);
      if (idx !== -1) FEATURE_VERSIONS.splice(idx, 1);
    });

    test('returns false when version is not yet loaded', function (assert) {
      service.version = null;
      assert.false(service.hasFeature(TEST_KEY), 'returns false when version is null');
    });

    test('returns false for an unknown key', function (assert) {
      service.version = '2.0.0';
      assert.false(service.hasFeature('this-key-does-not-exist'), 'returns false for unknown key');
    });

    test('returns true when current version equals the minimum version', function (assert) {
      service.version = '2.0.0';
      assert.true(service.hasFeature(TEST_KEY), 'returns true when version equals minimum');
    });

    test('returns true when current version exceeds the minimum version', function (assert) {
      service.version = '2.1.0';
      assert.true(service.hasFeature(TEST_KEY), 'returns true when version exceeds minimum');
    });

    test('returns false when current version is below the minimum version', function (assert) {
      service.version = '1.19.0';
      assert.false(service.hasFeature(TEST_KEY), 'returns false when version is below minimum');
    });

    test('handles enterprise version strings with +ent suffix', function (assert) {
      service.version = '2.0.0+ent';
      assert.true(service.hasFeature(TEST_KEY), 'handles +ent suffix correctly');
      service.version = '1.19.0+ent';
      assert.false(service.hasFeature(TEST_KEY), 'handles +ent suffix below minimum');
    });

    test('handles version strings with v prefix', function (assert) {
      service.version = 'v2.0.0';
      assert.true(service.hasFeature(TEST_KEY), 'handles v prefix correctly');
    });
  });
});
