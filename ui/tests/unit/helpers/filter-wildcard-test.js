/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { filterWildcard } from 'vault/helpers/filter-wildcard';
import { module, test } from 'qunit';

module('Unit | Helpers | filter-wildcard', function () {
  test('it returns a count if array contains a wildcard', function (assert) {
    const string = { id: 'foo*' };
    const array = ['foobar', 'foozar', 'boo', 'oof'];
    const result = filterWildcard([string, array]);
    assert.strictEqual(result, 2);
  });

  test('it returns zero if no wildcard is string', function (assert) {
    const string = { id: 'foo#' };
    const array = ['foobar', 'foozar', 'boo', 'oof'];
    const result = filterWildcard([string, array]);
    assert.strictEqual(result, 0);
  });

  test('it escapes function and does not error if no id is in string', function (assert) {
    const string = '*bar*';
    const array = ['foobar', 'foozar', 'boobarboo', 'oof'];
    const result = filterWildcard([string, array]);
    assert.strictEqual(result, 2);
  });
});
