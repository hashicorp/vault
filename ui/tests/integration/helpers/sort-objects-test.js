/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { sortObjects } from 'core/helpers/sort-objects';

module('Integration | Helper | sort-objects', function (hooks) {
  setupRenderingTest(hooks);

  test('it sorts-objects in template view', async function (assert) {
    this.array = [
      { foo: 'something-a', bar: 'third' },
      { foo: 'd-something', bar: 'second' },
      { foo: 'something-b', bar: 'fourth' },
      { foo: 'a-something', bar: 'first' },
    ];
    await render(hbs`
      {{#each (sort-objects this.array "foo") as |item|}}
      <p data-test-item> {{item.bar}} </p>
      {{/each}}
      `);

    assert.dom('[data-test-item]').hasText('first', 'renders items alphabetically');
  });

  test('it sorts array of objects', function (assert) {
    const originalArray = [
      { foo: 'grape', bar: 'third' },
      { foo: 'banana', bar: 'second' },
      { foo: 'lemon', bar: 'fourth' },
      { foo: 'apple', bar: 'first' },
    ];
    const expectedArray = [
      { bar: 'first', foo: 'apple' },
      { bar: 'second', foo: 'banana' },
      { bar: 'third', foo: 'grape' },
      { bar: 'fourth', foo: 'lemon' },
    ];

    assert.propEqual(sortObjects([originalArray, 'foo']), expectedArray, 'it sorts array of objects');

    const originalWithNumbers = [
      { foo: 'Z', bar: 'fourth' },
      { foo: '1', bar: 'first' },
      { foo: '2', bar: 'second' },
      { foo: 'A', bar: 'third' },
    ];
    const expectedWithNumbers = [
      { bar: 'first', foo: '1' },
      { bar: 'second', foo: '2' },
      { bar: 'third', foo: 'A' },
      { bar: 'fourth', foo: 'Z' },
    ];
    assert.propEqual(
      sortObjects([originalWithNumbers, 'foo']),
      expectedWithNumbers,
      'it sorts strings with numbers and letters'
    );
  });

  test('it disregards capitalization', function (assert) {
    // sort() arranges capitalized values before lowercase, the helper removes case by making all strings toUppercase()
    const originalArray = [
      { foo: 'something-a', bar: 'third' },
      { foo: 'D-something', bar: 'second' },
      { foo: 'SOMETHING-b', bar: 'fourth' },
      { foo: 'a-something', bar: 'first' },
    ];
    const expectedArray = [
      { bar: 'first', foo: 'a-something' },
      { bar: 'second', foo: 'D-something' },
      { bar: 'third', foo: 'something-a' },
      { bar: 'fourth', foo: 'SOMETHING-b' },
    ];

    assert.propEqual(
      sortObjects([originalArray, 'foo']),
      expectedArray,
      'it sorts array of objects regardless of capitalization'
    );
  });

  test('it fails gracefully', function (assert) {
    const originalArray = [
      { foo: 'b', bar: 'two' },
      { foo: 'a', bar: 'one' },
    ];
    assert.propEqual(
      sortObjects([originalArray, 'someKey']),
      originalArray,
      'it returns original array if key does not exist'
    );
    assert.deepEqual(
      sortObjects(['not an array']),
      'not an array',
      'it returns original arg if not an array'
    );

    const notStrings = [
      { foo: '1', bar: 'third' },
      { foo: 'Z', bar: 'second' },
      { foo: 1, bar: 'fourth' },
      { foo: 2, bar: 'first' },
    ];
    assert.propEqual(
      sortObjects([notStrings, 'foo']),
      notStrings,
      'it returns original array if values are not all strings'
    );
  });
});
