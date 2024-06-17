/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Service | auth', function (hooks) {
  setupTest(hooks);

  [
    ['#calculateExpiration w/ttl', { ttl: 30 }, 30],
    ['#calculateExpiration w/lease_duration', { lease_duration: 15 }, 15],
  ].forEach(([testName, response, ttlValue]) => {
    test(testName, function (assert) {
      const now = Date.now();
      const service = this.owner.lookup('service:auth');

      const resp = service.calculateExpiration(response, now);

      assert.strictEqual(resp.ttl, ttlValue, 'returns the ttl');
      assert.strictEqual(
        resp.tokenExpirationEpoch,
        now + ttlValue * 1e3,
        'calculates expiration from ttl as epoch timestamp'
      );
    });
  });

  test('#calculateExpiration w/ expire_time', function (assert) {
    const now = Date.now();
    const expirationString = '2024-06-13T09:10:21-07:00';
    const expectedExpirationEpoch = new Date(expirationString).getTime();

    const service = this.owner.lookup('service:auth');

    const resp = service.calculateExpiration({ ttl: 30, expire_time: '2024-06-13T09:10:21-07:00' }, now);

    assert.strictEqual(resp.ttl, 30, 'returns ttl');
    assert.strictEqual(
      resp.tokenExpirationEpoch,
      expectedExpirationEpoch,
      'calculates expiration from expire_time'
    );
  });
});
