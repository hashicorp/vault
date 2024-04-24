/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { numericalAxisLabel, calculateAverage, calculateSum } from 'vault/utils/chart-helpers';
import { module, test } from 'qunit';

const SMALL_NUMBERS = [0, 7, 27, 103, 999];
const LARGE_NUMBERS = {
  1001: '1k',
  1245: '1.2k',
  33777: '34k',
  532543: '530k',
  2100100: '2.1M',
  54500200100: '55B',
};

module('Unit | Utility | chart-helpers', function () {
  test('numericalAxisLabel renders number correctly', function (assert) {
    assert.expect(12);
    const method = numericalAxisLabel();
    assert.ok(method);
    SMALL_NUMBERS.forEach(function (num) {
      assert.strictEqual(numericalAxisLabel(num), num, `Does not format small number ${num}`);
    });
    Object.keys(LARGE_NUMBERS).forEach(function (num) {
      const expected = LARGE_NUMBERS[num];
      assert.strictEqual(numericalAxisLabel(num), expected, `Formats ${num} as ${expected}`);
    });
  });

  test('calculateAverage is accurate', function (assert) {
    const testArray1 = [
      { label: 'foo', value: 10 },
      { label: 'bar', value: 22 },
    ];
    const testArray2 = [
      { label: 'foo', value: undefined },
      { label: 'bar', value: 22 },
    ];
    const testArray3 = [{ label: 'foo' }, { label: 'bar' }];
    const getAverage = (array) => array.reduce((a, b) => a + b, 0) / array.length;
    assert.strictEqual(calculateAverage(null), null, 'returns null if dataset it null');
    assert.strictEqual(calculateAverage([]), null, 'returns null if dataset it empty array');
    assert.strictEqual(
      calculateAverage(testArray1, 'value'),
      getAverage([10, 22]),
      `returns correct average for array of objects`
    );
    assert.strictEqual(
      calculateAverage(testArray2, 'value'),
      getAverage([0, 22]),
      `returns correct average for array of objects containing undefined values`
    );
    assert.strictEqual(
      calculateAverage(testArray3, 'value'),
      null,
      'returns null when object key does not exist at all'
    );
  });

  test('calculateSum adds array of numbers', function (assert) {
    assert.strictEqual(calculateSum([2, 3]), 5, 'it sums array');
    assert.strictEqual(calculateSum(['one', 2]), null, 'returns null if array contains non-integers');
    assert.strictEqual(calculateSum('not an array'), null, 'returns null if an array is not passed');
  });
});
