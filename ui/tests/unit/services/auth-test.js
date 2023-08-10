/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Service | auth', function (hooks) {
  setupTest(hooks);

  [
    ['#calculateExpiration w/ttl', { ttl: 30 }, 30],
    ['#calculateExpiration w/lease_duration', { ttl: 15 }, 15],
  ].forEach(([testName, response, ttlValue]) => {
    test(testName, function (assert) {
      const now = Date.now();
      const service = this.owner.factoryFor('service:auth').create({
        now() {
          return now;
        },
      });

      const resp = service.calculateExpiration(response);

      assert.strictEqual(resp.ttl, ttlValue, 'returns the ttl');
      assert.strictEqual(
        resp.tokenExpirationEpoch,
        now + ttlValue * 1e3,
        'calculates expiration from ttl as epoch timestamp'
      );
    });
  });
});
