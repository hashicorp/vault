/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { FEATURE_VERSIONS } from 'vault/utils/feature-versions';

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
