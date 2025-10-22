/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { click } from '@ember/test-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | form-field-label', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function (assert) {
    this.setProperties({
      label: 'Test Label',
      helpText: null,
      subText: null,
      docLink: null,
    });

    await render(hbs`
      <FormFieldLabel
        @label={{this.label}}
        @helpText={{this.helpText}}
        @subText={{this.subText}}
        @docLink={{this.docLink}}
        for="some-input"
      />
    `);

    assert.dom('label').hasAttribute('for', 'some-input', 'Attributes passed to label element');
    assert.dom('label').hasText(this.label, 'Label text renders');
    assert.dom(GENERAL.tooltipText).doesNotExist('Help text hidden when not provided');
    assert.dom('.sub-text').doesNotExist('Sub text hidden when not provided');
    this.setProperties({
      helpText: 'More info',
      subText: 'Some description',
    });
    await click(GENERAL.tooltip('form field label'));
    assert.dom(GENERAL.tooltipText).hasText(this.helpText, 'Help text renders in tooltip');
    assert.dom('.sub-text').hasText(this.subText, 'Sub text renders');
    assert.dom('a').doesNotExist('docLink hidden when not provided');
    this.set('docLink', '/doc/path');
    assert.dom('.sub-text').includesText('See our documentation for help', 'Doc link text renders');
    assert
      .dom('a')
      .hasAttribute('href', 'https://developer.hashicorp.com' + this.docLink, 'Doc link renders');
  });
});
