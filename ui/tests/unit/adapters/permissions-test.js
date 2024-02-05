/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Adapter | permissions', function (hooks) {
  setupTest(hooks);

  test('it correctly calculates namespace access', function (assert) {
    assert.expect(7);
    const adapter = this.owner.lookup('adapter:permissions');
    const combinations = [
      { rootNs: 'ns1', currentNs: 'ns1', expected: true },
      { rootNs: 'ns2', currentNs: 'ns1', expected: false },
      { rootNs: '', currentNs: '', expected: true },
      { rootNs: '', currentNs: 'ns1', expected: true },
      { rootNs: 'ns1', currentNs: '', expected: false },
      { rootNs: 'ns1', currentNs: 'ns1/ns2', expected: true },
      { rootNs: 'ns1/ns2', currentNs: 'ns1/ns3', expected: false },
    ];
    combinations.forEach((c) => {
      assert.strictEqual(
        adapter.allowNsAccess(c.rootNs, c.currentNs),
        c.expected,
        `accessing ${c.currentNs} from ${c.rootNs} should be ${c.expected ? 'allowed' : 'denied'}`
      );
    });
  });
});
