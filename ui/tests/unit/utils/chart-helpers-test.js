/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { calculateSum, toFixedDisplay } from 'vault/utils/chart-helpers';
import { module, test } from 'qunit';

module('Unit | Utility | chart-helpers', function () {
  test('calculateSum adds array of numbers', function (assert) {
    assert.strictEqual(calculateSum([2, 3]), 5, 'it sums array');
    assert.strictEqual(calculateSum(['one', 2]), null, 'returns null if array contains non-integers');
    assert.strictEqual(calculateSum('not an array'), null, 'returns null if an array is not passed');
  });

  test('calculateSum with fixedDecimalPlaces parameter', function (assert) {
    assert.strictEqual(calculateSum([2.5, 3.7], 1), 6.2, 'rounds sum to 1 decimal place');
    assert.strictEqual(calculateSum([2.5555, 3.7777], 2), 6.33, 'rounds sum to 2 decimal places');
    assert.strictEqual(
      calculateSum([48.7888, 0.0112], 4),
      48.8,
      'handles floating-point precision issues with 4 decimal places'
    );
    assert.strictEqual(
      calculateSum([73.1832, 0.0168], 4),
      73.2,
      'correctly sums and rounds to 4 decimal places'
    );
    assert.strictEqual(
      calculateSum([10, 20, 30], 4),
      60,
      'works with whole numbers when fixedDecimalPlaces is provided'
    );
    assert.strictEqual(
      calculateSum([1.11111, 2.22222, 3.33333], 4),
      6.6667,
      'rounds sum of multiple numbers to 4 decimal places'
    );
    assert.strictEqual(calculateSum([0.1, 0.2], 4), 0.3, 'handles classic floating-point issue (0.1 + 0.2)');
    assert.strictEqual(calculateSum([2, 3], 0), 5, 'rounds to 0 decimal places (whole number)');
  });

  test('toFixedDisplay formats numbers with fixed decimal places', function (assert) {
    assert.strictEqual(toFixedDisplay(48.8, 4), '48.8000', 'formats number with trailing zeros');
    assert.strictEqual(toFixedDisplay(73.2, 4), '73.2000', 'preserves 4 decimal places');
    assert.strictEqual(toFixedDisplay(100, 2), '100.00', 'formats whole number with decimals');
    assert.strictEqual(toFixedDisplay(0, 4), 0, 'returns 0 as number, not formatted string');
    assert.strictEqual(toFixedDisplay(1.23456, 2), '1.23', 'rounds to specified decimal places');
    assert.strictEqual(toFixedDisplay('not a number', 4), 'not a number', 'returns non-number as-is');
    assert.strictEqual(toFixedDisplay(5.5, -1), 5.5, 'returns number as-is for negative decimal places');
  });

  test('calculateSum and toFixedDisplay work together', function (assert) {
    const sum1 = calculateSum([48.7888, 0.0112], 4);
    assert.strictEqual(sum1, 48.8, 'calculateSum returns number with fixed precision');
    assert.strictEqual(
      toFixedDisplay(sum1, 4),
      '48.8000',
      'toFixedDisplay formats for display with trailing zeros'
    );

    const sum2 = calculateSum([73.1832, 0.0168], 4);
    assert.strictEqual(sum2, 73.2, 'calculateSum handles floating-point precision');
    assert.strictEqual(toFixedDisplay(sum2, 4), '73.2000', 'toFixedDisplay preserves trailing zeros');

    const sum3 = calculateSum([10, 20, 30], 4);
    assert.strictEqual(sum3, 60, 'calculateSum works with whole numbers');
    assert.strictEqual(
      toFixedDisplay(sum3, 4),
      '60.0000',
      'toFixedDisplay adds decimal places to whole numbers'
    );

    const sum4 = calculateSum([0, 0, 0], 4);
    assert.strictEqual(sum4, 0, 'calculateSum returns 0 for zero sum');
    assert.strictEqual(toFixedDisplay(sum4, 4), 0, 'toFixedDisplay returns 0 as-is, not formatted');
  });
});
