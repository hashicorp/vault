/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, fillIn, click } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | filter-input', function (hooks) {
  setupRenderingTest(hooks);

  test('it should render attributes on input element', async function (assert) {
    assert.expect(3);

    this.onClick = () => assert.ok(true, 'on click modifier passed to input element');

    await render(
      hbs`<FilterInput aria-label="test-component" placeholder="Filter roles" @value="foo" {{on "click" this.onClick}} />`
    );
    await click('[data-test-filter-input]');
    assert
      .dom('[data-test-filter-input]')
      .hasAttribute('placeholder', 'Filter roles', 'Placeholder passed to input element');
    assert.dom('[data-test-filter-input]').hasValue('foo', 'Value passed to input element');
  });

  test('it should send input event', async function (assert) {
    assert.expect(1);

    this.onInput = (value) => {
      assert.strictEqual(value, 'foo', 'onInput event sent with value');
    };

    await render(hbs`<FilterInput aria-label="test-component" @wait={{0}} @onInput={{this.onInput}} />`);
    await fillIn('[data-test-filter-input]', 'foo');
  });
});
