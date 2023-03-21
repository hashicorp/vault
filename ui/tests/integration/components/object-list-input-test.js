/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import sinon from 'sinon';

const SELECTORS = {
  listLabel: (key) => `[data-test-object-list-label="${key}"]`,
  listInput: (key, row) => `[data-test-object-list-input="${key}-${row}"]`,
  addButton: '[data-test-object-list-add-button]',
  deleteButton: (row) => `[data-test-object-list-delete-button="${row}"]`,
};
module('Integration | Component | object-list-input', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.keys = [
      { label: 'Input A', key: 'inputA', placeholder: 'Input something here' },
      { label: 'Input B', key: 'inputB', placeholder: 'Input another thing' },
      { label: 'Input C', key: 'inputC' },
      { label: 'Input D', key: 'inputD' },
      { label: 'Input E', key: 'inputE' },
    ];
    this.inputValue = [
      { inputA: 'foo', inputB: 'bar' },
      { inputA: 'another', inputB: 'value' },
    ];
    this.onChange = sinon.spy();
  });

  test('it renders with correct number of inputs and labels', async function (assert) {
    assert.expect(12);
    await render(hbs`<ObjectListInput @objectKeys={{this.keys}} @onChange={{this.onChange}} />`);
    for (let i = 0; i < this.keys.length; i++) {
      const element = this.keys[i];
      assert.dom(SELECTORS.listLabel(element.key)).hasText(element.label, 'it renders labels');
      assert.dom(SELECTORS.listInput(element.key, 0)).exists('it renders input');
    }
    const firstKey = this.keys[0];
    assert
      .dom(SELECTORS.listInput(firstKey.key, 0))
      .hasAttribute('placeholder', firstKey.placeholder, 'it renders placeholder text');
    assert.dom(SELECTORS.addButton).isDisabled('add button is disabled with empty inputs');
  });

  test('it adds a new row', async function (assert) {
    assert.expect(12);
    const expectedArray = [
      {
        inputA: 'foo-0',
        inputB: 'foo-1',
        inputC: 'foo-2',
        inputD: 'foo-3',
        inputE: 'foo-4',
      },
      {
        inputA: 'bar-0',
        inputB: 'bar-1',
        inputC: 'bar-2',
        inputD: 'bar-3',
        inputE: 'bar-4',
      },
    ];
    await render(hbs`<ObjectListInput @objectKeys={{this.keys}} @onChange={{this.onChange}} />`);

    assert.dom(SELECTORS.addButton).isDisabled('add button is disabled with empty inputs');
    await fillIn(SELECTORS.listInput(this.keys[0].key, 0), 'foo');
    assert.propEqual(
      this.onChange.lastCall.args[0],
      [
        {
          inputA: 'foo',
          inputB: '',
          inputC: '',
          inputD: '',
          inputE: '',
        },
      ],
      'onChange called with single input'
    );
    assert.dom(SELECTORS.addButton).isDisabled('add button is still disabled with single input');

    // fill in each input
    for (let i = 0; i < this.keys.length; i++) {
      await fillIn(SELECTORS.listInput(this.keys[i].key, 0), `foo-${i}`);
    }
    assert.dom(SELECTORS.addButton).isEnabled('add button enabled when all inputs are filled');
    assert.propEqual(
      this.onChange.lastCall.args[0],
      expectedArray.slice(0, 1),
      'onChange called with first row of inputs'
    );

    // add a row
    await click(SELECTORS.addButton);
    assert.propEqual(
      this.onChange.lastCall.args[0],
      expectedArray.slice(0, 1),
      'onChange is called with only filled in row (does not include empty row)'
    );

    for (let i = 0; i < this.keys.length; i++) {
      const element = this.keys[i];
      assert.dom(SELECTORS.listLabel(element.key)).exists({ count: 1 }, 'label only renders for first row');
    }

    // fill in another row of inputs
    for (let i = 0; i < this.keys.length; i++) {
      await fillIn(SELECTORS.listInput(this.keys[i].key, 1), `bar-${i}`);
    }

    assert.propEqual(
      this.onChange.lastCall.args[0],
      expectedArray,
      'onChange includes both first and second row of objects'
    );
  });

  test('it renders with inputValues and deletes a row', async function (assert) {
    assert.expect(12);
    this.keys = this.keys.slice(0, 2);
    const [firstColumn, secondColumn] = this.keys;
    await render(hbs`
      <ObjectListInput
        @objectKeys={{this.keys}}
        @onChange={{this.onChange}} 
        @inputValue={{this.inputValue}} 
      />`);

    assert.dom(SELECTORS.listInput(firstColumn.key, 0)).hasValue('foo', 'input exists in first row');
    assert.dom(SELECTORS.listInput(secondColumn.key, 0)).hasValue('bar', 'input exists in first row');
    assert.dom(SELECTORS.listInput(firstColumn.key, 1)).hasValue('another', 'input exists in second row');
    assert.dom(SELECTORS.listInput(secondColumn.key, 1)).hasValue('value', 'input exists in second row');
    assert.dom(SELECTORS.listInput(firstColumn.key, 2)).hasNoValue('empty input renders for first key');
    assert.dom(SELECTORS.listInput(secondColumn.key, 2)).hasNoValue('empty input renders for second key');

    assert.dom(SELECTORS.deleteButton(0)).exists({ count: 1 }, 'renders delete button first row');
    assert.dom(SELECTORS.deleteButton(1)).exists({ count: 1 }, 'renders delete button second row');
    assert.dom(SELECTORS.addButton).exists({ count: 1 }, 'renders one add button');
    assert.dom(SELECTORS.addButton).isDisabled('add button is disabled when inputValue exists');

    assert.ok(this.onChange.notCalled, 'on change does not fire when rendering input values');

    await click(SELECTORS.deleteButton(1));

    assert.propEqual(
      this.onChange.lastCall.args[0],
      this.inputValue.slice(0, 1),
      'onChange fires with deleted row'
    );
  });
});
