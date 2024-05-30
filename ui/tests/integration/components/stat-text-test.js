/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, settled } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

module('Integration | Component | StatText', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function (assert) {
    await render(hbs`<StatText />`);

    assert.dom('[data-test-stat-text]').exists('renders element');
  });

  test('it renders passed in attributes', async function (assert) {
    this.set('label', 'A Label');
    this.set('value', 'A value');
    this.set('size', 'l');
    this.set('subText', 'This is my description');

    await render(
      hbs`<StatText @label={{this.label}} @size={{this.size}} @value={{this.value}} @subText={{this.subText}}/>`
    );
    assert.dom('.stat-label').hasText(this.label, 'renders label');
    assert.dom('.stat-text').hasText(this.subText, 'renders subtext');
    assert.dom('.stat-value').hasText(this.value, 'renders a non-integer value');

    this.set('value', 604099);
    await settled();

    const formattedNumber = '604,099';
    assert.dom('.stat-value').hasText(formattedNumber, 'renders correctly formatted integer value');
  });
});
