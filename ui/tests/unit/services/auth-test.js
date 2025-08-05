/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Service | auth', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    this.service = this.owner.lookup('service:auth');
  });

  module('#calculateExpiration', function () {
    [
      ['#calculateExpiration w/ttl', { ttl: 30 }, 30],
      ['#calculateExpiration w/lease_duration', { lease_duration: 15 }, 15],
    ].forEach(([testName, response, ttlValue]) => {
      test(testName, function (assert) {
        const now = Date.now();

        const resp = this.service.calculateExpiration(response, now);

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

      const resp = this.service.calculateExpiration(
        { ttl: 30, expire_time: '2024-06-13T09:10:21-07:00' },
        now
      );

      assert.strictEqual(resp.ttl, 30, 'returns ttl');
      assert.strictEqual(
        resp.tokenExpirationEpoch,
        expectedExpirationEpoch,
        'calculates expiration from expire_time'
      );
    });
  });

  module('#setExpirationSettings', function () {
    test('#setExpirationSettings for a renewable token', function (assert) {
      const now = Date.now();
      const ttl = 30;
      const response = { ttl, renewable: true };

      this.service.setExpirationSettings(response, now);

      assert.false(this.service.allowExpiration, 'sets allowExpiration to false');
      assert.strictEqual(this.service.expirationCalcTS, now, 'sets expirationCalcTS to now');
    });

    test('#setExpirationSettings for a non-renewable token', function (assert) {
      const now = Date.now();
      const ttl = 30;
      const response = { ttl, renewable: false };

      this.service.setExpirationSettings(response, now);

      assert.true(this.service.allowExpiration, 'sets allowExpiration to true');
      assert.strictEqual(this.service.expirationCalcTS, null, 'keeps expirationCalcTS as null');
    });
  });
});
