/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { stringArrayToCamelCase } from 'vault/helpers/string-array-to-camel';

module('Integration | Helper | string-array-to-camel', function (hooks) {
  setupRenderingTest(hooks);

  test('it returns camelCase string with all caps and two words separated by space', async function (assert) {
    const string = 'FOO Bar';
    const expected = 'fooBar';
    const result = stringArrayToCamelCase(string);

    assert.strictEqual(
      result,
      expected,
      'camelCase string returned for call caps and two words separated by space'
    );
  });

  test('it returns an array of camelCased strings if an array of strings passed in', function (assert) {
    const string = ['FOO Bar', 'Baz Qux', 'wibble wobble', 'wobble WIBBLes'];
    const expected = ['fooBar', 'bazQux', 'wibbleWobble', 'wobbleWibbles'];
    const result = stringArrayToCamelCase(string);
    assert.deepEqual(result, expected, 'camelCase array of strings returned for all sorts of strings');
  });

  test('it returns string if string is numbers', function (assert) {
    const string = '123';
    const expected = '123';
    const result = stringArrayToCamelCase(string);
    assert.strictEqual(result, expected, 'camelCase kind of handles strings with numbers');
  });

  test('it returns error if str is not a string', function (assert) {
    const string = { name: 'foo.bar*baz' };
    let result;
    try {
      result = stringArrayToCamelCase(string);
    } catch (e) {
      result = e.message;
    }
    assert.deepEqual(result, 'Assertion Failed: must pass in a string or array of strings');
  });

  test('it returns error if str is not an array of string', function (assert) {
    const string = [{ name: 'foo.bar*baz' }];
    let result;
    try {
      result = stringArrayToCamelCase(string);
    } catch (e) {
      result = e.message;
    }
    assert.deepEqual(result, 'Assertion Failed: must pass in a string or array of strings');
  });
});
