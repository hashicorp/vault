/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { isAdvancedSecret, obfuscateData } from 'core/utils/advanced-secret';
import { module, test } from 'qunit';

module('Unit | Utility | advanced-secret', function () {
  module('isAdvancedSecret', function () {
    test('it returns false for non-valid JSON', function (assert) {
      assert.expect(5);
      let result;
      ['some-string', 'character{string', '{value}', '[blah]', 'multi\nline\nstring'].forEach((value) => {
        result = isAdvancedSecret('some-string');
        assert.false(result, `returns false for ${value}`);
      });
    });

    test('it returns false for single-level objects', function (assert) {
      assert.expect(3);
      let result;
      [{ single: 'one' }, { first: '1', two: 'three' }, ['my', 'array']].forEach((value) => {
        result = isAdvancedSecret(JSON.stringify(value));
        assert.false(result, `returns false for object ${JSON.stringify(value)}`);
      });
    });

    test('it returns true for any nested object', function (assert) {
      assert.expect(3);
      let result;
      [
        { single: { one: 'uno' } },
        { first: ['this', 'counts\ntoo'] },
        { deeply: { nested: { item: 1 } } },
      ].forEach((value) => {
        result = isAdvancedSecret(JSON.stringify(value));
        assert.true(result, `returns true for object ${JSON.stringify(value)}`);
      });
    });
  });
  module('obfuscateData', function () {
    test('it obfuscates values of an object', function (assert) {
      assert.expect(4);
      [
        {
          name: 'flat map',
          data: {
            first: 'one',
            second: 'two',
            third: 'three',
          },
          obscured: {
            first: '********',
            second: '********',
            third: '********',
          },
        },
        {
          name: 'nested map',
          data: {
            first: 'one',
            second: {
              third: 'two',
            },
          },
          obscured: {
            first: '********',
            second: {
              third: '********',
            },
          },
        },
        {
          name: 'numbers and arrays',
          data: {
            first: 1,
            list: ['one', 'two'],
            second: {
              third: ['one', 'two'],
              number: 5,
            },
          },
          obscured: {
            first: '********',
            list: ['********', '********'],
            second: {
              third: ['********', '********'],
              number: '********',
            },
          },
        },
        {
          name: 'object arrays',
          data: {
            list: [{ one: 'one' }, { two: 'two' }],
          },
          obscured: {
            list: ['********', '********'],
          },
        },
      ].forEach((test) => {
        const result = obfuscateData(test.data);
        assert.deepEqual(result, test.obscured, `obfuscates values of ${test.name}`);
      });
    });

    test('it does not obfuscate non-object values', function (assert) {
      assert.expect(3);
      ['some-string', 5, ['my', 'array']].forEach((test) => {
        const result = obfuscateData(test);
        assert.deepEqual(result, test, `does not obfuscate value ${JSON.stringify(test)}`);
      });
    });
  });
});
