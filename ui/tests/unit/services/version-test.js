/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

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
});
