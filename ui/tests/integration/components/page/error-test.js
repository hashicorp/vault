/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

module('Integration | Component | page/error', function (hooks) {
  setupRenderingTest(hooks);

  test('it should render 404 error', async function (assert) {
    this.error = {
      httpStatus: 404,
      path: '/v1/kubernetes/config',
    };

    await render(hbs`<Page::Error @error={{this.error}} />`);

    assert.dom('h1').hasText('404 Not Found', 'Error title renders');
    assert
      .dom('p')
      .hasText(`Sorry, we were unable to find any content at ${this.error.path}.`, 'Error message renders');
  });

  test('it should render 403 error', async function (assert) {
    this.error = {
      httpStatus: 403,
      path: '/v1/kubernetes/config',
    };

    await render(hbs`<Page::Error @error={{this.error}} />`);

    assert.dom('h1').hasText('Not authorized', 'Error title renders');
    assert
      .dom('p')
      .hasText(`You are not authorized to access content at ${this.error.path}.`, 'Error message renders');
  });

  test('it should render general error', async function (assert) {
    this.error = {
      message: 'An unexpected error occurred',
      errors: ['This is one thing that went wrong', 'Unfortunately something else went wrong too'],
    };

    await render(hbs`<Page::Error @error={{this.error}} />`);

    assert.dom('h1').hasText('Error', 'Error title renders');
    assert.dom('[data-test-page-error-message]').hasText(this.error.message, 'Error message renders');
    this.error.errors.forEach((error, index) => {
      assert
        .dom(`[data-test-page-error-details="${index}"]`)
        .hasText(this.error.errors[index], 'Error detail renders');
    });
  });
});
