/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, findAll } from '@ember/test-helpers';
import sinon from 'sinon';
import hbs from 'htmlbars-inline-precompile';

const handler = (data, e) => {
  if (e && e.preventDefault) e.preventDefault();
  return;
};

module('Integration | Component | toggle', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders with name as default label', async function (assert) {
    this.set('handler', sinon.spy(handler));

    await render(hbs`<Toggle
      @onChange={{this.handler}}
      @name="thing"
    />`);

    assert.dom(findAll('label')[0]).hasText('thing');

    await render(hbs`
      <Toggle
        @onChange={{this.handler}}
        @name="thing"
      >
        <span id="test-value" class="has-text-grey">template block text</span>
      </Toggle>
    `);
    assert.dom('[data-test-toggle-label="thing"]').exists('toggle label exists');
    assert.dom('#test-value').hasText('template block text', 'yielded text renders');
  });

  test('it has the correct classes', async function (assert) {
    this.set('handler', sinon.spy(handler));
    await render(hbs`
      <Toggle
        @onChange={{this.handler}}
        @name="thing"
        @size="small"
      >
        template block text
      </Toggle>
    `);
    assert.dom('.toggle.is-small').exists('toggle has is-small class');
  });

  test('it sets the id of the input correctly', async function (assert) {
    this.set('handler', sinon.spy(handler));
    await render(hbs`
    <Toggle
      @onChange={{this.handler}}
      @name="my toggle #_has we!rd chars"
    >
      Label
    </Toggle>
    `);
    assert.dom('#toggle-mytoggle_haswerdchars').exists('input has correct ID');
    assert
      .dom('label')
      .hasAttribute('for', 'toggle-mytoggle_haswerdchars', 'label has correct for attribute');
  });
});
