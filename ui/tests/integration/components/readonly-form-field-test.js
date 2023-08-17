/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import EmberObject from '@ember/object';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

const minimumAttr = {
  name: 'my-input',
  type: 'text',
};
const customLabelAttr = {
  name: 'test-input',
  type: 'text',
  options: {
    subText: 'Subtext here',
    label: 'Custom-label',
  },
};

module('Integration | Component | readonly-form-field', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function (assert) {
    this.set('attr', EmberObject.create(minimumAttr));
    await render(hbs`<ReadonlyFormField @attr={{this.attr}} @value="value" />`);
    assert
      .dom('[data-test-readonly-label]')
      .includesText('My input', 'formats the attr name when no label provided');
    assert.dom(`[data-test-input="${minimumAttr.name}"]`).hasValue('value', 'Uses the value as passed');
    assert.dom(`[data-test-input="${minimumAttr.name}"]`).hasAttribute('readonly');
  });

  test('it renders with options', async function (assert) {
    this.set('attr', customLabelAttr);
    await render(hbs`<ReadonlyFormField @attr={{this.attr}} @value="another value" />`);
    assert
      .dom('[data-test-readonly-label]')
      .includesText('Custom-label', 'Uses the provided label as passed');
    assert.dom('.sub-text').includesText('Subtext here', 'Renders subtext');
    assert
      .dom(`[data-test-input="${customLabelAttr.name}"]`)
      .hasValue('another value', 'Uses the value as passed');
    assert.dom(`[data-test-input="${customLabelAttr.name}"]`).hasAttribute('readonly');
  });
});
