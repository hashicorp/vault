/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, fillIn } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';

const OPTIONS = ['foo', 'bar', 'baz'];
const LABEL = 'Boop';

module('Integration | Component | Select', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.set('options', OPTIONS);
    this.set('label', LABEL);
    this.set('name', 'foo');
  });

  test('it renders', async function (assert) {
    await render(hbs`<Select @options={{this.options}} @label={{this.label}} @name={{this.name}}/>`);
    assert.dom('[data-test-select-label]').hasText('Boop');
    assert.dom('[data-test-select="foo"]').exists();
  });

  test('it renders when options is an array of strings', async function (assert) {
    await render(hbs`<Select @options={{this.options}} @label={{this.label}} @name={{this.name}}/>`);

    assert.dom('[data-test-select="foo"]').hasValue('foo');
    assert.strictEqual(this.element.querySelector('[data-test-select="foo"]').options.length, 3);
  });

  test('it renders when options is an array of objects', async function (assert) {
    const objectOptions = [
      { value: 'berry', label: 'Berry' },
      { value: 'cherry', label: 'Cherry' },
    ];
    this.set('options', objectOptions);
    await render(hbs`<Select @options={{this.options}} @label={{this.label}} @name={{this.name}}/>`);

    assert.dom('[data-test-select="foo"]').hasValue('berry');
    assert.strictEqual(this.element.querySelector('[data-test-select="foo"]').options.length, 2);
  });

  test('it renders when options is an array of custom objects', async function (assert) {
    const objectOptions = [
      { day: 'mon', fullDay: 'Monday' },
      { day: 'tues', fullDay: 'Tuesday' },
    ];
    const selectedValue = objectOptions[1].day;
    this.setProperties({
      options: objectOptions,
      valueAttribute: 'day',
      labelAttribute: 'fullDay',
      selectedValue: selectedValue,
    });

    await render(
      hbs`
        <Select
          @options={{this.options}}
          @label={{this.label}}
          @name={{this.name}}
          @valueAttribute={{this.valueAttribute}}
          @labelAttribute={{this.labelAttribute}}
          @selectedValue={{this.selectedValue}}/>`
    );

    assert.dom('[data-test-select="foo"]').hasValue('tues', 'sets selectedValue by default');
    assert
      .dom(this.element.querySelector('[data-test-select="foo"]').options[1])
      .hasText('Tuesday', 'uses the labelAttribute to determine the label');
  });

  test('it calls onChange when an item is selected', async function (assert) {
    this.set('onChange', sinon.spy());
    await render(
      hbs`<Select @label={{this.label}} @options={{this.options}} @name={{this.name}} @onChange={{this.onChange}}/>`
    );
    await fillIn('[data-test-select="foo"]', 'bar');

    assert.ok(this.onChange.calledOnce);
  });
});
