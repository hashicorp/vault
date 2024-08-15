/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import queryParamString from 'vault/utils/query-param-string';
import { module, test } from 'qunit';

module('Unit | Utility | query-param-string', function () {
  [
    {
      scenario: 'object with nonencoded keys and values',
      obj: { redirect: 'https://hashicorp.com', some$key: 'normal-value', number: 7 },
      expected: '?redirect=https%3A%2F%2Fhashicorp.com&some%24key=normal-value&number=7',
    },
    {
      scenario: 'object with invalid values',
      obj: { redirect: '', null: null, foo: 'bar', number: 0, undefined: undefined, false: false },
      expected: '?foo=bar&number=0&false=false',
    },
    {
      scenario: 'object where every value is invalid',
      obj: { redirect: '', null: null, undefined: undefined },
      expected: '',
    },
    {
      scenario: 'empty object',
      obj: {},
      expected: '',
    },
    {
      scenario: 'array',
      obj: ['some', 'array'],
      expected: '',
    },
    {
      scenario: 'string',
      obj: 'foobar',
      expected: '',
    },
  ].forEach((testCase) => {
    test(`it works when ${testCase.scenario}`, function (assert) {
      const result = queryParamString(testCase.obj);
      assert.strictEqual(result, testCase.expected);
    });
  });
});
