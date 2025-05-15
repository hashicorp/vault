/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { getErrorResponse } from 'vault/tests/helpers/api/error-response';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

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

  test('it should handle api client errors', async function (assert) {
    this.error = getErrorResponse();

    await render(hbs`<Page::Error @error={{this.error}} />`);

    assert
      .dom(GENERAL.pageError.errorTitle('404'))
      .hasText('404 Not Found', 'Error title renders based on status');
    assert
      .dom(GENERAL.pageError.errorSubtitle)
      .hasText(
        'Sorry, we were unable to find any content at /v1/test/error/parsing.',
        'Error subtitle renders'
      );

    const error = { errors: ['something bad happened'], message: 'bad things occurred' };
    this.error = getErrorResponse(error, 400);

    await render(hbs`<Page::Error @error={{this.error}} />`);

    assert.dom(GENERAL.pageError.errorTitle('400')).hasText('Error', 'Error title renders');
    assert.dom(GENERAL.pageError.errorMessage).hasText(error.message, 'Error message renders');
    assert.dom(GENERAL.pageError.errorDetails).hasText(error.errors[0], 'Error details render');
  });
});
