/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { addToArray } from '../../../helpers/add-to-array';

module('Integration | Helper | add-to-array', function (hooks) {
  setupRenderingTest(hooks);

  test('it correctly adds a value to an array without mutating the original', function (assert) {
    const ARRAY = ['horse', 'cow', 'chicken'];
    const result = addToArray(ARRAY, 'pig');
    assert.deepEqual(result, [...ARRAY, 'pig'], 'Result has additional item');
    assert.deepEqual(ARRAY, ['horse', 'cow', 'chicken'], 'original array is not mutated');
  });

  test('it fails if the first value is not an array', function (assert) {
    let result;
    try {
      result = addToArray('not-array', 'string');
    } catch (e) {
      result = e.message;
    }
    assert.strictEqual(result, 'Assertion Failed: Value provided is not an array');
  });

  test('it works with non-string arrays', function (assert) {
    const ARRAY = ['five', 6, '7'];
    const result = addToArray(ARRAY, 10);
    assert.deepEqual(result, ['five', 6, '7', 10], 'added number value');
  });

  test('it de-dupes the result', function (assert) {
    const ARRAY = ['horse', 'cow', 'chicken'];
    const result = addToArray(ARRAY, 'horse');
    assert.deepEqual(result, ['horse', 'cow', 'chicken']);
  });
});
