/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { isWildcardString } from 'vault/helpers/is-wildcard-string';
import { module, test } from 'qunit';

module('Unit | Helpers | is-wildcard-string', function () {
  test('it returns true if regular string with wildcard', function (assert) {
    const string = 'foom#*eep';
    const result = isWildcardString([string]);
    assert.true(result);
  });

  test('it returns false if no wildcard', function (assert) {
    const string = 'foo.bar';
    const result = isWildcardString([string]);
    assert.false(result);
  });

  test('it returns true if string with id as in searchSelect selected has wildcard', function (assert) {
    const string = { id: 'foo.bar*baz' };
    const result = isWildcardString([string]);
    assert.true(result);
  });

  test('it returns true if string object has name and no id', function (assert) {
    const string = { name: 'foo.bar*baz' };
    const result = isWildcardString([string]);
    assert.true(result);
  });

  test('it returns true if string object has name and id with at least one wildcard', function (assert) {
    const string = { id: '7*', name: 'seven' };
    const result = isWildcardString([string]);
    assert.true(result);
  });

  test('it returns true if string object has name and id with wildcard in name not id', function (assert) {
    const string = { id: '7', name: 'sev*n' };
    const result = isWildcardString([string]);
    assert.true(result);
  });
});
