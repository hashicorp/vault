/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import Sinon from 'sinon';

module('Integration | Component | checkbox-grid', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.name = 'fooBar';
    this.label = 'Foo bar';
    this.fields = [
      { key: 'abc', label: 'All Bears Cry' },
      { key: 'def', label: 'Dark Eel Feelings' },
    ];

    this.onChange = Sinon.spy();
  });

  test('it renders with minimum inputs', async function (assert) {
    const changeSpy = Sinon.spy();
    this.set('onChange', changeSpy);
    await render(
      hbs`<CheckboxGrid @name={{this.name}} @label={{this.label}} @fields={{this.fields}} @onChange={{this.onChange}} />`
    );

    assert.dom('[data-test-checkbox]').exists({ count: 2 }, 'One checkbox is rendered for each field');
    assert.dom('[data-test-checkbox]').isNotChecked('no fields are checked by default');
    await click('[data-test-checkbox="abc"]');
    assert.ok(changeSpy.calledOnceWithExactly('fooBar', ['abc']));
  });

  test('it renders with values set', async function (assert) {
    const changeSpy = Sinon.spy();
    this.set('onChange', changeSpy);
    this.set('currentValue', ['abc']);
    await render(
      hbs`<CheckboxGrid @name={{this.name}} @label={{this.label}} @fields={{this.fields}} @onChange={{this.onChange}} @value={{this.currentValue}} />`
    );

    assert.dom('[data-test-checkbox]').exists({ count: 2 }, 'One checkbox is rendered for each field');
    assert.dom('[data-test-checkbox="abc"]').isChecked('abc field is checked on load');
    assert.dom('[data-test-checkbox="def"]').isNotChecked('def field is unchecked on load');
    await click('[data-test-checkbox="abc"]');
    assert.ok(changeSpy.calledOnceWithExactly('fooBar', []), 'Sends correct payload when unchecking');
    await click('[data-test-checkbox="def"]');
    await click('[data-test-checkbox="abc"]');
    assert.ok(
      changeSpy.calledWithExactly('fooBar', ['def', 'abc']),
      'sends correct payload with multiple checked'
    );
  });
});
