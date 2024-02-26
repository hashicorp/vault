/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

module('Integration | Component | EnableInput', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders and enables yielded input', async function (assert) {
    assert.expect(4);
    await render(hbs`
    <EnableInput @label="my-attr-name">
      <label for="foo">Foo</label>
      <Input data-test-yielded-input @type='text' id="foo" />
    </EnableInput>
      `);
    assert.dom('input').hasAttribute('readonly', '', 'input is read only');
    assert.dom('input').hasValue('**********', 'disabled input renders asterisks');
    await click('button');
    assert
      .dom('[data-test-yielded-input]')
      .doesNotHaveAttribute('readonly', 'toggles to enabled, yielded input');
    assert.dom('button').doesNotExist('button disappears when input is enabled');
  });

  test('it renders passed attribute', async function (assert) {
    assert.expect(6);
    this.attr = {
      name: 'specialClientCredentials',
      type: 'string',
      options: {
        subText: 'This value is protected and not returned from the API. Enable input to update value.',
      },
    };
    this.model = { specialClientCredentials: '' };
    await render(hbs`
    <EnableInput @attr={{this.attr}} >
      <FormField @attr={{this.attr}} @model={{this.model}} />
    </EnableInput>
      `);

    assert
      .dom(`[data-test-input="${this.attr.name}"]`)
      .hasAttribute('readonly', '', 'renders readonly ReadonlyFormField');
    assert
      .dom(`[data-test-input="${this.attr.name}"]`)
      .hasValue('**********', 'readonly input renders asterisks');
    assert.dom('[data-test-readonly-label]').hasText('Special client credentials');
    assert.dom('p.sub-text').hasText(this.attr.options.subText);
    await click('button');
    assert
      .dom(`[data-test-field="${this.attr.name}"] input`)
      .doesNotHaveAttribute('readonly', 'toggles to enabled, yielded form field component');
    assert.dom('button').doesNotExist('button disappears when input is enabled');
  });
});
