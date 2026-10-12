/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { MAX_DURATION_SECONDS, isDurationOverMax } from 'core/utils/duration-utils';

module('Unit | Util | duration-utils', function () {
  test('MAX_DURATION_SECONDS matches the largest duration Vault accepts', function (assert) {
    // Go's time.Duration is int64 nanoseconds: 9223372036854775807ns truncates to 9223372036 whole seconds.
    assert.strictEqual(MAX_DURATION_SECONDS, 9223372036, 'max is 9223372036 seconds');
  });

  test('isDurationOverMax accepts values Vault accepts', function (assert) {
    [
      ['', 'empty string'],
      [undefined, 'undefined'],
      [null, 'null'],
      [0, 'zero as a number'],
      ['0', 'zero as a string'],
      ['0s', 'zero seconds'],
      ['30s', 'seconds'],
      ['24h', 'hours'],
      ['9223372036s', 'the maximum in seconds'],
      ['153722867m', 'the largest whole number of minutes'],
      ['2562047h', 'the largest whole number of hours'],
      ['106751d', 'the largest whole number of days'],
      [9223372036, 'the maximum as a number'],
      ['9223372036', 'the maximum as plain digits'],
    ].forEach(([value, label]) => {
      assert.false(isDurationOverMax(value), `${label} is not over the max`);
    });
  });

  test('isDurationOverMax rejects values Vault rejects as too large', function (assert) {
    [
      ['9223372037s', 'one second over the maximum'],
      ['153722868m', 'one minute over in minutes'],
      ['2562048h', 'one hour over in hours'],
      ['106752d', 'one day over in days'],
      ['863999999913600s', '9999999999 days as sent by the TTL picker'],
      ['1e+21s', 'exponent notation produced by very large numbers'],
      [9223372037, 'one second over as a number'],
      ['9223372037', 'one second over as plain digits'],
      ['Infinitys', 'the overflow value the TTL picker produces for very long digit input'],
      [Infinity, 'infinity as a number'],
    ].forEach(([value, label]) => {
      assert.true(isDurationOverMax(value), `${label} is over the max`);
    });
  });

  test('isDurationOverMax leaves values it cannot parse to the API', function (assert) {
    ['abc', 'off', '5x'].forEach((value) => {
      assert.false(isDurationOverMax(value), `"${value}" is not reported as over the max`);
    });
  });
});
