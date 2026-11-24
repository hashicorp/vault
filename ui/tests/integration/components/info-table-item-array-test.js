/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { findAll, render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

const DISPLAY_ARRAY = ['role-1', 'role-2', 'role-3', 'role-4', 'role-5'];

module('Integration | Component | InfoTableItemArray', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.set('displayArray', DISPLAY_ARRAY);
    this.set('isLink', true);
    this.set('queryParam', 'role');
    this.set('label', 'Roles');
    this.set('arrayOptions', [...DISPLAY_ARRAY, 'role-6']);
    this.set('wildcardLabel', 'role');
  });

  test('it renders', async function (assert) {
    await render(hbs`<InfoTableItemArray
        @displayArray={{this.displayArray}}
        @label="my label"
      />`);

    assert.dom('[data-test-info-table-item-array]').exists();
    const noLinkString = document.querySelector('code').textContent.trim();
    assert.strictEqual(
      noLinkString.length,
      DISPLAY_ARRAY.toString().length,
      'renders a string of the array if isLink is not provided'
    );
  });

  test('it renders links if isLink is true', async function (assert) {
    await render(hbs`
      <InfoTableItemArray
        @displayArray={{this.displayArray}}
        @isLink={{this.isLink}}
        @label="my label"
        @queryParam={{this.queryParam}}
      />
    `);
    assert.strictEqual(
      document.querySelectorAll('a > span').length,
      DISPLAY_ARRAY.length,
      'renders each item in array with link'
    );
  });

  test('it renders wildcard items as plain text and still shows view all', async function (assert) {
    const displayArrayWithWildcard = ['role-1', 'role-2', 'role-3', 'r*'];
    this.set('displayArrayWithWildcard', displayArrayWithWildcard);
    await render(hbs`
    <InfoTableItemArray
      @label={{this.label}}
      @displayArray={{this.displayArrayWithWildcard}}
      @isLink={{this.isLink}}
      @queryParam={{this.queryParam}}
      @arrayOptions={{this.arrayOptions}}
      @wildcardLabel={{this.wildcardLabel}}
    />`);
    assert.strictEqual(
      document.querySelectorAll('a > span').length,
      displayArrayWithWildcard.length - 1,
      'renders only the non-wildcard items in array with link'
    );
    assert.dom('[data-test-item="r*"]').hasText('r*', 'renders wildcard as plain text');
    assert.dom('[data-test-count="6"]').exists('renders a wildcard count badge from preloaded options');
    assert.dom('[data-test-view-all="roles"]').doesNotExist('does not render view all text for short arrays');
  });

  test('it renders a wildcard item in the truncated array', async function (assert) {
    const displayArrayWithWildcard = [
      'role-1',
      'role-2',
      'role-3',
      'r*',
      'role-4',
      'role-5',
      'role-6',
      'role-7',
      'role-8',
      'role-9',
      'role-10',
    ];
    this.set('displayArrayWithWildcard', displayArrayWithWildcard);
    await render(hbs`
    <InfoTableItemArray
      @label={{this.label}}
      @displayArray={{this.displayArrayWithWildcard}}
      @isLink={{this.isLink}}
      @queryParam={{this.queryParam}}
      @arrayOptions={{this.arrayOptions}}
      @wildcardLabel={{this.wildcardLabel}}
    />`);
    const numberCutOffTruncatedArray = displayArrayWithWildcard.length - 5;
    assert.strictEqual(
      document.querySelectorAll('a > span').length,
      5,
      'renders only the non-wildcard items in the truncated array with link'
    );
    assert.dom('[data-test-item="r*"]').hasText('r*', 'renders wildcard item in the truncated array');
    assert.dom('[data-test-count="6"]').exists('renders a wildcard count badge in the truncated array');
    assert
      .dom(`[data-test-and="${numberCutOffTruncatedArray}"]`)
      .exists("renders correct 'and N others' text");
    assert.dom('[data-test-view-all="roles"]').hasText('View all roles.', 'renders correct view all text');
  });

  test('it omits the wildcard badge when preloaded options are not provided', async function (assert) {
    this.set('displayArrayWithWildcard', ['role-1', 'r*']);
    await render(hbs`
      <InfoTableItemArray
        @label={{this.label}}
        @displayArray={{this.displayArrayWithWildcard}}
        @isLink={{this.isLink}}
        @queryParam={{this.queryParam}}
      />
    `);

    assert.dom('[data-test-item="r*"]').hasText('r*', 'renders wildcard as plain text');
    assert.dom('[data-test-count]').doesNotExist('does not render a badge without preloaded options');
  });

  test('it truncates arrays with linked items', async function (assert) {
    const value = ['1', '2', '3-id', '4', '5', '6', '7', '8', '9', '10'];
    this.set('value', value);
    await render(hbs`
    <InfoTableItemArray
      @label="Entities"
      @displayArray={{this.value}}
      @isLink={{this.isLink}}
    />`);
    assert.dom('[data-test-item="1"]').hasText('1', 'renders the raw item text');
    assert.dom('[data-test-item="3-id"]').hasText('3-id', 'renders item text without name lookup');
    assert.strictEqual(findAll('[data-test-item]').length, 5, 'only lists 5 entities');
  });

  test('it truncates using read more component when overflows div', async function (assert) {
    const value = [
      'entity-name-1-with-extra-text',
      'entity-name-2-with-extra-text',
      'entity-name-3-with-extra-text',
      'entity-name-4-with-extra-text',
      'entity-name-5-with-extra-text',
      'entity-name-6-with-extra-text',
      'entity-name-7-with-extra-text',
      'entity-name-8-with-extra-text',
      'entity-name-9-with-extra-text',
      'entity-name-10-with-extra-text',
    ];
    this.set('value', value);
    await render(hbs`
      <div style="width: 200px">
        <InfoTableItemArray
          @label="Entities"
          @displayArray={{this.value}}
          @isLink={{this.isLink}}
          @doNotTruncate={{true}}
        />
      </div>
    `);
    assert.dom('[data-test-readmore-toggle]').exists('renders see more toggle');
    assert.dom('[data-test-view-all]').doesNotExist('Does not render view all text');
  });
});
