/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import apiErrorMessage from 'vault/utils/api-error-message';

module('Unit | Util | api-error-message', function (hooks) {
  hooks.beforeEach(function () {
    this.apiError = {
      errors: ['first error', 'second error'],
      message: 'there were some errors',
    };
    this.getErrorContext = () => ({ response: new Response(JSON.stringify(this.apiError)) });
  });

  test('it should return errors from ErrorContext', async function (assert) {
    const message = await apiErrorMessage(this.getErrorContext());
    assert.strictEqual(message, 'first error, second error');
  });

  test('it should return message from ErrorContext when errors are empty', async function (assert) {
    this.apiError.errors = [];
    const message = await apiErrorMessage(this.getErrorContext());
    assert.strictEqual(message, 'there were some errors');
  });

  test('it should return fallback message for ErrorContext without errors or message', async function (assert) {
    this.apiError = {};
    const message = await apiErrorMessage(this.getErrorContext());
    assert.strictEqual(message, 'An error occurred, please try again');
  });

  test('it should return message from Error', async function (assert) {
    const error = new Error('some js type error');
    const message = await apiErrorMessage(error);
    assert.strictEqual(message, error.message);
  });

  test('it should return default fallback', async function (assert) {
    const message = await apiErrorMessage('some random error');
    assert.strictEqual(message, 'An error occurred, please try again');
  });

  test('it should return custom fallback message', async function (assert) {
    const fallback = 'Everything is broken, sorry';
    const message = await apiErrorMessage('some random error', fallback);
    assert.strictEqual(message, fallback);
  });
});
