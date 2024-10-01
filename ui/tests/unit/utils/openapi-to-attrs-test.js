/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { combineFieldGroups } from 'vault/utils/openapi-to-attrs';
import { module, test } from 'qunit';

module('Unit | Util | combineFieldGroups', function () {
  const NEW_FIELDS = ['one', 'two', 'three'];

  test('it adds new fields from OpenAPI to fieldGroups except for exclusions', function (assert) {
    assert.expect(3);
    const modelFieldGroups = [
      { default: ['name', 'awesomePeople'] },
      {
        Options: ['ttl'],
      },
    ];
    const excludedFields = ['two'];
    const expectedGroups = [
      { default: ['name', 'awesomePeople', 'one', 'three'] },
      {
        Options: ['ttl'],
      },
    ];
    const newFieldGroups = combineFieldGroups(modelFieldGroups, NEW_FIELDS, excludedFields);
    for (const groupName in modelFieldGroups) {
      assert.deepEqual(
        newFieldGroups[groupName],
        expectedGroups[groupName],
        'it incorporates all new fields except for those excluded'
      );
    }
  });
  test('it adds all new fields from OpenAPI to fieldGroups when excludedFields is empty', function (assert) {
    assert.expect(3);
    const modelFieldGroups = [
      { default: ['name', 'awesomePeople'] },
      {
        Options: ['ttl'],
      },
    ];
    const excludedFields = [];
    const expectedGroups = [
      { default: ['name', 'awesomePeople', 'one', 'two', 'three'] },
      {
        Options: ['ttl'],
      },
    ];
    const nonExcludedFieldGroups = combineFieldGroups(modelFieldGroups, NEW_FIELDS, excludedFields);
    for (const groupName in modelFieldGroups) {
      assert.deepEqual(
        nonExcludedFieldGroups[groupName],
        expectedGroups[groupName],
        'it incorporates all new fields'
      );
    }
  });
  test('it keeps fields the same when there are no brand new fields from OpenAPI', function (assert) {
    assert.expect(3);
    const modelFieldGroups = [
      { default: ['name', 'awesomePeople', 'two', 'one', 'three'] },
      {
        Options: ['ttl'],
      },
    ];
    const excludedFields = [];
    const expectedGroups = [
      { default: ['name', 'awesomePeople', 'two', 'one', 'three'] },
      {
        Options: ['ttl'],
      },
    ];
    const fieldGroups = combineFieldGroups(modelFieldGroups, NEW_FIELDS, excludedFields);
    for (const groupName in modelFieldGroups) {
      assert.deepEqual(fieldGroups[groupName], expectedGroups[groupName], 'it incorporates all new fields');
    }
  });
});
