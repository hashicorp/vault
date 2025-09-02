/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { removeManyFromArray, removeFromArray } from 'vault/helpers/remove-from-array';

module('Integration | Helper | remove-from-array', function (hooks) {
  setupRenderingTest(hooks);

  test('it correctly removes a value from an array without mutating the original', function (assert) {
    const ARRAY = ['horse', 'cow', 'chicken'];
    const result = removeFromArray(ARRAY, 'horse');
    assert.deepEqual(result, ['cow', 'chicken'], 'Result does not have removed item');
    assert.deepEqual(ARRAY, ['horse', 'cow', 'chicken'], 'original array is not mutated');
  });

  test('it returns the same value if the item is not found', function (assert) {
    const ARRAY = ['horse', 'cow', 'chicken'];
    const result = removeFromArray(ARRAY, 'pig');
    assert.deepEqual(result, ARRAY, 'Results are the same as original array');
  });

  test('it fails if the first value is not an array', function (assert) {
    let result;
    try {
      result = removeFromArray('not-array', 'string');
    } catch (e) {
      result = e.message;
    }
    assert.strictEqual(result, 'Assertion Failed: Value provided is not an array');
  });

  test('it works with non-string arrays', function (assert) {
    const ARRAY = ['five', 6, '7'];
    const result1 = removeFromArray(ARRAY, 6);
    const result2 = removeFromArray(ARRAY, 7);
    assert.deepEqual(result1, ['five', '7'], 'removed number value');
    assert.deepEqual(result2, ARRAY, 'did not match on different types');
  });

  test('it de-dupes the result', function (assert) {
    const ARRAY = ['horse', 'cow', 'chicken', 'cow'];
    const result = removeFromArray(ARRAY, 'horse');
    assert.deepEqual(result, ['cow', 'chicken']);
  });

  test('it works with two arrays', function (assert) {
    const ARRAY = ['five', 6, '7'];
    const result1 = removeManyFromArray(ARRAY, [6, '7']);
    const result2 = removeManyFromArray(ARRAY, ['foo', 'five']);
    assert.deepEqual(result1, ['five'], 'removed multiple values');
    assert.deepEqual(result2, [6, '7'], 'did nothing with values that were not in the array');
  });
});
