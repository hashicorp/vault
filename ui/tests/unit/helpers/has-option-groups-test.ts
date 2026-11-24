/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import type Helper from '@ember/component/helper';
import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Helper | has-option-groups', function (hooks) {
  setupTest(hooks);

  let helper: Helper;

  hooks.beforeEach(function () {
    helper = this.owner.lookup('helper:has-option-groups') as Helper;
  });

  test('returns false for null or undefined', function (assert) {
    assert.false(helper.compute([null], {}), 'returns false for null');
    assert.false(helper.compute([undefined], {}), 'returns false for undefined');
  });

  test('returns false for empty array', function (assert) {
    assert.false(helper.compute([[]], {}), 'returns false for empty array');
  });

  test('returns false for flat array of strings', function (assert) {
    const flatArray = ['option1', 'option2', 'option3'];
    assert.false(helper.compute([flatArray], {}), 'returns false for flat string array');
  });

  test('returns false for flat array of objects without group property', function (assert) {
    const flatObjects = [
      { value: 'opt1', displayName: 'Option 1' },
      { value: 'opt2', displayName: 'Option 2' },
    ];
    assert.false(helper.compute([flatObjects], {}), 'returns false for objects without group property');
  });

  test('returns true for array with grouped options', function (assert) {
    const groupedArray = [
      {
        group: 'Group A',
        options: ['option1', 'option2'],
      },
      {
        group: 'Group B',
        options: [
          { value: 'opt3', displayName: 'Option 3' },
          { value: 'opt4', displayName: 'Option 4' },
        ],
      },
    ];
    assert.true(helper.compute([groupedArray], {}), 'returns true for grouped options');
  });

  test('returns true if any item has group property', function (assert) {
    const mixedArray = [
      'simpleString',
      { value: 'opt1', displayName: 'Option 1' },
      {
        group: 'Group C',
        options: ['option3'],
      },
    ];
    assert.true(helper.compute([mixedArray], {}), 'returns true when at least one item has group property');
  });
});
