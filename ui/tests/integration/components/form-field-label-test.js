/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { click } from '@ember/test-helpers';

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
    assert.dom('[data-test-help-text]').doesNotExist('Help text hidden when not provided');
    assert.dom('.sub-text').doesNotExist('Sub text hidden when not provided');
    this.setProperties({
      helpText: 'More info',
      subText: 'Some description',
    });
    await click('[data-test-tool-tip-trigger]');
    assert.dom('[data-test-help-text]').hasText(this.helpText, 'Help text renders in tooltip');
    assert.dom('.sub-text').hasText(this.subText, 'Sub text renders');
    assert.dom('a').doesNotExist('docLink hidden when not provided');
    this.set('docLink', '/doc/path');
    assert.dom('.sub-text').includesText('See our documentation for help', 'Doc link text renders');
    assert
      .dom('a')
      .hasAttribute('href', 'https://developer.hashicorp.com' + this.docLink, 'Doc link renders');
  });
});
