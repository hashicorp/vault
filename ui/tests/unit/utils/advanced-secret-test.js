/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { isAdvancedSecret } from 'core/utils/advanced-secret';
import { module, test } from 'qunit';

module('Unit | Utility | advanced-secret', function () {
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
