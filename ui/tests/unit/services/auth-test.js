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
    test('with a non-zero ttl value', function (assert) {
      const now = Date.now();
      const ttl = 30;
      const expireTime = null;
      const calculatedExpiry = this.service.calculateExpiration({ now, ttl, expireTime });

      assert.strictEqual(calculatedExpiry.ttl, 30, 'returns the ttl');
      assert.strictEqual(
        calculatedExpiry.tokenExpirationEpoch,
        now + ttl * 1e3,
        'calculates expiration from ttl as epoch timestamp'
      );
    });

    test('with a zero ttl value', function (assert) {
      const now = Date.now();
      const ttl = 0;
      const expireTime = null;
      const calculatedExpiry = this.service.calculateExpiration({ now, ttl, expireTime });

      assert.strictEqual(calculatedExpiry.ttl, null, 'returns `null` for the ttl');
      assert.strictEqual(calculatedExpiry.tokenExpirationEpoch, null, 'tokenExpirationEpoch is null');
    });

    test('#calculateExpiration w/ expireTime', function (assert) {
      const now = Date.now();
      const expirationString = '2024-06-13T09:10:21-07:00';
      const expectedExpirationEpoch = new Date(expirationString).getTime();

      const calculatedExpiry = this.service.calculateExpiration({
        now,
        ttl: 30,
        expireTime: '2024-06-13T09:10:21-07:00',
      });

      assert.strictEqual(calculatedExpiry.ttl, 30, 'returns ttl');
      assert.strictEqual(
        calculatedExpiry.tokenExpirationEpoch,
        expectedExpirationEpoch,
        'calculates expiration from expireTime'
      );
    });
  });

  module('#setExpirationSettings', function () {
    test('#setExpirationSettings for a renewable token', function (assert) {
      const now = Date.now();
      const renewable = true;

      this.service.setExpirationSettings(renewable, now);

      assert.false(this.service.allowExpiration, 'sets allowExpiration to false');
      assert.strictEqual(this.service.expirationCalcTS, now, 'sets expirationCalcTS to now');
    });

    test('#setExpirationSettings for a non-renewable token', function (assert) {
      const now = Date.now();
      const renewable = false;

      this.service.setExpirationSettings(renewable, now);

      assert.true(this.service.allowExpiration, 'sets allowExpiration to true');
      assert.strictEqual(this.service.expirationCalcTS, null, 'keeps expirationCalcTS as null');
    });
  });
});
