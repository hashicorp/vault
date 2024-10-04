/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { isAdvancedSecret } from 'core/utils/advanced-secret';
import { module, test } from 'qunit';

module('Unit | Utility | advanced-secret', function () {
  module('isAdvancedSecret', function () {
    test('it returns false for non-valid JSON', function (assert) {
      assert.expect(5);
      let result;
      ['some-string', 'character{string', '{value}', '[blah]', 'multi\nline\nstring'].forEach((value) => {
        result = isAdvancedSecret(value);
        assert.false(result, `returns false for ${value}`);
      });
    });

    test('it returns false for single-level objects', function (assert) {
      assert.expect(6);
      let result;
      [{ single: 'one' }, { first: '1', two: 'three' }, ['my', 'array']].forEach((value) => {
        const stringValue = JSON.stringify(value);
        result = isAdvancedSecret(value);
        assert.false(result, `returns false for object ${stringValue}`);
        result = isAdvancedSecret(stringValue);
        assert.false(result, `returns false for json ${stringValue}`);
      });
    });

    test('it returns true for any nested object or number value', function (assert) {
      assert.expect(8);
      let result;
      [
        { single: { one: 'uno' } },
        { first: ['this', 'counts\ntoo'] },
        { deeply: { nested: { item: 1 } } },
        { number: 5 },
      ].forEach((value) => {
        const stringValue = JSON.stringify(value);
        result = isAdvancedSecret(value);
        assert.true(result, `returns true for object ${stringValue}`);
        result = isAdvancedSecret(stringValue);
        assert.true(result, `returns true for json ${stringValue}`);
      });
    });
  });
});
