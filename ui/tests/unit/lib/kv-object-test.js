/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import KVObject from 'vault/lib/kv-object';

module('Unit | Lib | kv object', function () {
  const fromJSONTests = [
    [
      'types',
      { string: 'string', false: false, zero: 0, number: 1, null: null, true: true, object: { one: 'two' } },
      [
        { name: 'false', value: false },
        { name: 'null', value: null },
        { name: 'number', value: 1 },
        { name: 'object', value: { one: 'two' } },
        { name: 'string', value: 'string' },
        { name: 'true', value: true },
        { name: 'zero', value: 0 },
      ],
    ],
    [
      'ordering',
      { b: 'b', 1: '1', z: 'z', A: 'A', a: 'a' },
      [
        { name: '1', value: '1' },
        { name: 'a', value: 'a' },
        { name: 'A', value: 'A' },
        { name: 'b', value: 'b' },
        { name: 'z', value: 'z' },
      ],
    ],
  ];

  fromJSONTests.forEach(function ([name, input, content]) {
    test(`fromJSON: ${name}`, function (assert) {
      const data = KVObject.create({ content: [] }).fromJSON(input);
      assert.deepEqual(data.get('content'), content, 'has expected content');
    });
  });

  test(`fromJSON: non-object input`, function (assert) {
    const input = [{ foo: 'bar' }];
    assert.throws(
      () => {
        KVObject.create({ content: [] }).fromJSON(input);
      },
      /Vault expects data to be formatted as an JSON object/,
      'throws when non-object input is used to construct the KVObject'
    );
  });

  fromJSONTests.forEach(function ([name, input, content]) {
    test(`fromJSONString: ${name}`, function (assert) {
      const inputString = JSON.stringify(input, null, 2);
      const data = KVObject.create({ content: [] }).fromJSONString(inputString);
      assert.deepEqual(data.get('content'), content, 'has expected content');
    });
  });

  const toJSONTests = [
    [
      'types',
      false,
      { string: 'string', false: false, zero: 0, number: 1, null: null, true: true, object: { one: 'two' } },
      { false: false, null: null, number: 1, object: { one: 'two' }, string: 'string', true: true, zero: 0 },
    ],
    ['include blanks = true', true, { string: 'string', '': '' }, { string: 'string', '': '' }],
    ['include blanks = false', false, { string: 'string', '': '' }, { string: 'string' }],
  ];

  toJSONTests.forEach(function ([name, includeBlanks, input, output]) {
    test(`toJSON: ${name}`, function (assert) {
      const data = KVObject.create({ content: [] }).fromJSON(input);
      const result = data.toJSON(includeBlanks);
      assert.deepEqual(result, output, 'has expected output');
    });
  });

  toJSONTests.forEach(function ([name, includeBlanks, input, output]) {
    test(`toJSONString: ${name}`, function (assert) {
      const expected = JSON.stringify(output, null, 2);
      const data = KVObject.create({ content: [] }).fromJSON(input);
      const result = data.toJSONString(includeBlanks);
      assert.strictEqual(result, expected, 'has expected output string');
    });
  });

  const isAdvancedTests = [
    [
      'advanced',
      { string: 'string', false: false, zero: 0, number: 1, null: null, true: true, object: { one: 'two' } },
      true,
    ],
    ['string-only', { string: 'string', one: 'two' }, false],
  ];

  isAdvancedTests.forEach(function ([name, input, expected]) {
    test(`isAdvanced: ${name}`, function (assert) {
      const data = KVObject.create({ content: [] }).fromJSON(input);

      assert.strictEqual(data.isAdvanced(), expected, 'calculates isAdvanced correctly');
    });
  });
});
