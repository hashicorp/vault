/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Helper | format-duration', function (hooks) {
  setupRenderingTest(hooks);

  test('it supports strings and formats seconds', async function (assert) {
    await render(hbs`<p data-test-format-duration>Date: {{format-duration '3606'}}</p>`);

    assert
      .dom('[data-test-format-duration]')
      .includesText('1 hour 6 seconds', 'it renders the duration in hours and seconds');
  });

  test('it is able to format seconds and days', async function (assert) {
    await render(hbs`<p data-test-format-duration>Date: {{format-duration '93606000'}}</p>`);

    assert
      .dom('[data-test-format-duration]')
      .includesText(
        '2 years 11 months 18 days 9 hours 40 minutes',
        'it renders with years months and days and hours and minutes'
      );
  });

  test('it is able to format numbers', async function (assert) {
    this.set('number', 60);
    await render(hbs`<p data-test-format-duration>Date: {{format-duration this.number}}</p>`);

    assert
      .dom('[data-test-format-duration]')
      .includesText('1 minute', 'it renders duration when a number is passed in.');
  });

  test('it renders the input if time not found', async function (assert) {
    this.set('number', 'arg');

    await render(hbs`<p data-test-format-duration>Date: {{format-duration this.number}}</p>`);
    assert.dom('[data-test-format-duration]').hasText('Date: arg');
  });

  test('it renders no value if nullable true', async function (assert) {
    this.set('number', 0);

    await render(hbs`<p data-test-format-duration>Date: {{format-duration this.number nullable=true}}</p>`);
    assert.dom('[data-test-format-duration]').hasText('Date:');
  });
});
