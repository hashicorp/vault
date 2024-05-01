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
});
