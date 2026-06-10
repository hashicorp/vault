/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';

module('Integration | Component | form/v2/error-alert', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders the default title and error message', async function (assert) {
    this.error = 'Something went wrong';

    await render(hbs`
      <Form::V2::ErrorAlert @error={{this.error}} />
    `);

    assert.dom('.hds-alert').exists('renders alert');
    assert.dom('.hds-alert__title').hasText('Submission error', 'renders default title');
    assert.dom('.hds-alert__description').hasText('Something went wrong', 'renders error message');
  });

  test('it renders a custom title', async function (assert) {
    this.error = 'Configuration failed';
    this.title = 'Configuration Error';

    await render(hbs`
      <Form::V2::ErrorAlert @error={{this.error}} @title={{this.title}} />
    `);

    assert.dom('.hds-alert__title').hasText('Configuration Error', 'renders custom title');
    assert.dom('.hds-alert__description').hasText('Configuration failed', 'renders error message');
  });

  test('it renders without an error message', async function (assert) {
    await render(hbs`
      <Form::V2::ErrorAlert />
    `);

    assert.dom('.hds-alert').exists('renders alert');
    assert.dom('.hds-alert__title').hasText('Submission error', 'renders default title');
    assert.dom('.hds-alert__description').hasText('', 'renders empty description');
  });
});
