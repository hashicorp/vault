/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import Service from '@ember/service';
import { setupRenderingTest } from 'ember-qunit';
import { findAll, render } from '@ember/test-helpers';
import { run } from '@ember/runloop';
import hbs from 'htmlbars-inline-precompile';

const DISPLAY_ARRAY = ['role-1', 'role-2', 'role-3', 'role-4', 'role-5'];

const storeService = Service.extend({
  query(modelType) {
    return new Promise((resolve, reject) => {
      switch (modelType) {
        case 'transform/role':
          resolve([
            { id: 'role-1' },
            { id: 'role-2' },
            { id: 'role-3' },
            { id: 'role-4' },
            { id: 'role-5' },
            { id: 'role-6' },
          ]);
          break;
        case 'model/no-permission':
          reject({ httpStatus: 403, message: 'permission denied' });
          break;
        case 'identity/entity':
          resolve([
            { id: '1', name: 'one' },
            { id: '6', name: 'six' },
            { id: '7', name: 'seven' },
            { id: '8', name: 'eight' },
            { id: '9', name: 'nine' },
          ]);
          break;
        default:
          reject({ httpStatus: 404 });
          break;
      }
    });
  },
});

module('Integration | Component | InfoTableItemArray', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.set('displayArray', DISPLAY_ARRAY);
    this.set('isLink', true);
    this.set('modelType', 'transform/role');
    this.set('queryParam', 'role');
    this.set('backend', 'transform');
    this.set('wildcardLabel', 'role');
    this.set('label', 'Roles');
    run(() => {
      this.owner.unregister('service:store');
      this.owner.register('service:store', storeService);
    });
  });

  hooks.afterEach(function () {
    this.owner.unregister('service:store');
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
        @modelType={{this.modelType}}
        @queryParam={{this.queryParam}}
        @backend={{this.backend}}
      />
    `);
    assert.strictEqual(
      document.querySelectorAll('a > span').length,
      DISPLAY_ARRAY.length,
      'renders each item in array with link'
    );
  });

  test('it renders a badge and view all if wildcard in display array && < 10', async function (assert) {
    const displayArrayWithWildcard = ['role-1', 'role-2', 'role-3', 'r*'];
    this.set('displayArrayWithWildcard', displayArrayWithWildcard);
    await render(hbs`
    <InfoTableItemArray
      @label={{this.label}}
      @displayArray={{this.displayArrayWithWildcard}}
      @isLink={{this.isLink}}
      @modelType={{this.modelType}}
      @queryParam={{this.queryParam}}
      @backend={{this.backend}}
    />`);
    assert.strictEqual(
      document.querySelectorAll('a > span').length,
      DISPLAY_ARRAY.length - 1,
      'renders each item in array with link'
    );
    // 6 here comes from the six roles setup in the store service.
    assert.dom('[data-test-count="6"]').exists('correctly counts with wildcard filter and shows the count');
    assert.dom('[data-test-view-all="roles"]').exists({ count: 1 }, 'renders 1 view all roles');
    assert.dom('[data-test-view-all="roles"]').hasText('View all roles.', 'renders correct view all text');
  });

  test('it renders a badge and view all if wildcard in display array && >= 10', async function (assert) {
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
      @modelType={{this.modelType}}
      @queryParam={{this.queryParam}}
      @backend={{this.backend}}
    />`);
    const numberCutOffTruncatedArray = displayArrayWithWildcard.length - 5;
    assert.strictEqual(document.querySelectorAll('a > span').length, 5, 'renders truncated array of five');
    assert
      .dom(`[data-test-and="${numberCutOffTruncatedArray}"]`)
      .exists('correctly counts with wildcard filter and shows the count');
    assert.dom('[data-test-view-all="roles"]').hasText('View all roles.', 'renders correct view all text');
  });

  test('it fails gracefully if query returns 403 and display array contains wildcard', async function (assert) {
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
    this.set('modelType', 'model/no-permission');
    await render(hbs`
    <InfoTableItemArray
      @label={{this.label}}
      @displayArray={{this.displayArrayWithWildcard}}
      @isLink={{this.isLink}}
      @modelType={{this.modelType}}
      @queryParam={{this.queryParam}}
      @backend={{this.backend}}
    />`);
    assert.strictEqual(findAll('[data-test-item]').length, 4, 'lists 4 roles');
    assert.dom('[data-test-readmore-content]').hasTextContaining('r*', 'renders wildcard');
    assert.dom('[data-test-count="0"]').doesNotExist('does not render badge');
    assert.dom('[data-test-view-all="roles"]').hasText('View all roles.', 'renders correct view all text');
    assert.dom('[data-test-and="6"]').exists(`renders correct 'and 6 others' text`);
  });

  test('it fails gracefully if query returns 404 and display array contains wildcard', async function (assert) {
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
    this.set('modelType', 'model-not-found');
    await render(hbs`
    <InfoTableItemArray
      @label={{this.label}}
      @displayArray={{this.displayArrayWithWildcard}}
      @isLink={{this.isLink}}
      @modelType={{this.modelType}}
      @queryParam={{this.queryParam}}
      @backend={{this.backend}}
    />`);
    assert.dom('[data-test-count="0"]').hasText('includes 0', 'renders badge');
    assert.strictEqual(findAll('[data-test-item]').length, 4, 'renders list of 4 roles');
    assert.dom('[data-test-view-all="roles"]').hasText('View all roles.', 'renders view all text');
  });

  test('it renders name if renderItemName=true or id if name not found', async function (assert) {
    const value = ['6', '8', '123-id'];
    this.set('value', value);
    this.set('modelType', 'identity/entity');
    await render(hbs`
    <InfoTableItemArray
      @label={{this.label}}
      @displayArray={{this.value}}
      @isLink={{this.isLink}}
      @modelType={{this.modelType}}
      @renderItemName={{true}}
    />`);
    assert.dom('[data-test-item="6"]').hasText('six', `renders name of 'six' instead of id`);
    assert.dom('[data-test-item="8"]').hasText('eight', `renders 'eight' instead of id`);
    assert.strictEqual(findAll('[data-test-item]').length, 3, 'renders all entities');
    assert
      .dom('[data-test-item="123-id"]')
      .hasText('123-id', 'renders id instead of name if no record for name');
  });

  test('it truncates and renders name when renderItemName=true', async function (assert) {
    const value = ['1', '2', '3-id', '4', '5', '6', '7', '8', '9', '10'];
    this.set('value', value);
    this.set('modelType', 'identity/entity');
    await render(hbs`
    <InfoTableItemArray
      @label="Entities"
      @displayArray={{this.value}}
      @isLink={{this.isLink}}
      @modelType={{this.modelType}}
      @renderItemName={{true}}
    />`);
    assert.dom('[data-test-item="1"]').hasText('one', `renders name of 'one' instead of id`);
    assert.dom('[data-test-item="3-id"]').hasText('3-id', 'renders id instead of name if no record for name');
    assert.strictEqual(findAll('[data-test-item]').length, 5, 'only lists 5 entities');
  });

  test('it truncates using read more component when overflows div', async function (assert) {
    const value = ['1', '2', '3-id', '4', '5', '6', '7', '8', '9', '10'];
    this.set('value', value);
    this.set('modelType', 'identity/entity');
    await render(hbs`
      <div style="width: 200px">
        <InfoTableItemArray
          @label="Entities"
          @displayArray={{this.value}}
          @isLink={{this.isLink}}
          @modelType={{this.modelType}}
          @renderItemName={{true}}
          @doNotTruncate={{true}}
        />
      </div>
    `);
    assert.dom('[data-test-readmore-toggle]').exists('renders see more toggle');
    assert.dom('[data-test-view-all]').doesNotExist('Does not render view all text');
  });
});
