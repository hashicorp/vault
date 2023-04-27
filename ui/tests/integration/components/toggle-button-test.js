/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | toggle-button', function (hooks) {
  setupRenderingTest(hooks);

  test('toggle functionality', async function (assert) {
    await render(hbs`
      <ToggleButton
        @isOpen={{this.isOpen}}
        @openLabel={{this.openLabel}}
        @closedLabel={{this.closedLabel}}
        @onClick={{fn (mut this.isOpen)}}
        data-test-toggle-button
      />
    `);

    assert.dom('button').hasText('More options', 'renders default closedLabel');
    await click('button');
    assert.true(this.isOpen, 'it updates the value on click');
    assert.dom('button').hasText('Hide options', 'renders default openLabel');
    await click('button');
    assert.false(this.isOpen, 'it updates the value on click');

    this.setProperties({
      openLabel: 'Close the options!',
      closedLabel: 'Open the options!',
    });

    assert.dom('button').hasText('Open the options!', 'renders passed closedLabel');
    await click('button');
    assert.dom('button').hasText('Close the options!', 'renders passed openLabel');
    assert
      .dom('button')
      .hasAttribute('data-test-toggle-button', '', 'Attributes are spread on the button element');
  });
});
