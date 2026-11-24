/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Service | replication-mode', function (hooks) {
  setupTest(hooks);

  test('it has no mode until one is set', function (assert) {
    const service = this.owner.lookup('service:replication-mode');
    assert.strictEqual(service.mode, null);
    assert.strictEqual(service.getMode(), null);
  });

  test('setMode updates the mode returned by mode and getMode', function (assert) {
    const service = this.owner.lookup('service:replication-mode');

    service.setMode('dr');
    assert.strictEqual(service.mode, 'dr');
    assert.strictEqual(service.getMode(), 'dr');

    service.setMode('performance');
    assert.strictEqual(service.mode, 'performance');
    assert.strictEqual(service.getMode(), 'performance');
  });
});
