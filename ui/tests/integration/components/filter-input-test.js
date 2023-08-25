/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, fillIn } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | filter-input', function (hooks) {
  setupRenderingTest(hooks);

  test('it should render placeholder and send input event', async function (assert) {
    assert.expect(2);

    this.onInput = (value) => {
      assert.strictEqual(value, 'foo', 'onInput event sent with value');
    };

    await render(hbs`<FilterInput @placeholder="Filter roles" @onInput={{this.onInput}} />`);

    assert
      .dom('[data-test-filter-input]')
      .hasAttribute('placeholder', 'Filter roles', 'Placeholder set on input element');

    await fillIn('[data-test-filter-input]', 'foo');
  });
});
