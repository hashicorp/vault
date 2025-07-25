/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { include } from 'vault/helpers/include';

module('Unit | Helper | include', function () {
  test('it returns true when item is in array', function (assert) {
    const result = include([['a', 'b', 'c'], 'b']);
    assert.true(result, 'returns true when item is found in array');
  });

  test('it returns false when item is not in array', function (assert) {
    const result = include([['a', 'b', 'c'], 'd']);
    assert.false(result, 'returns false when item is not found in array');
  });

  test('it returns false when array is empty', function (assert) {
    const result = include([[], 'a']);
    assert.false(result, 'returns false when array is empty');
  });

  test('it returns false when array is not an array', function (assert) {
    const result = include([null, 'a']);
    assert.false(result, 'returns false when first argument is null');

    const result2 = include([undefined, 'a']);
    assert.false(result2, 'returns false when first argument is undefined');

    const result3 = include(['not-an-array', 'a']);
    assert.false(result3, 'returns false when first argument is not an array');
  });

  test('it works with different data types', function (assert) {
    const result1 = include([[1, 2, 3], 2]);
    assert.true(result1, 'works with numbers');

    const result2 = include([['kv_123', 'pki_456'], 'kv_123']);
    assert.true(result2, 'works with strings');

    const result3 = include([[true, false], true]);
    assert.true(result3, 'works with booleans');
  });
});
