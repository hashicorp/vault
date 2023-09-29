/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Service | version', function (hooks) {
  setupTest(hooks);

  test('setting version computes isOSS properly', function (assert) {
    const service = this.owner.lookup('service:version');
    service.version = '0.9.5';
    assert.true(service.isOSS);
    assert.false(service.isEnterprise);
  });

  test('setting version computes isEnterprise properly', function (assert) {
    const service = this.owner.lookup('service:version');
    service.version = '0.9.5+ent';
    assert.false(service.isOSS);
    assert.true(service.isEnterprise);
  });

  test('setting version with hsm ending computes isEnterprise properly', function (assert) {
    const service = this.owner.lookup('service:version');
    service.version = '0.9.5+ent.hsm';
    assert.false(service.isOSS);
    assert.true(service.isEnterprise);
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
