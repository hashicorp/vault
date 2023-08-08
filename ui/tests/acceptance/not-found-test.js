/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { visit } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import authPage from 'vault/tests/pages/auth';
import Ember from 'ember';

let adapterException;

module('Acceptance | not-found', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    adapterException = Ember.Test.adapter.exception;
    Ember.Test.adapter.exception = () => {};
    return authPage.login();
  });

  hooks.afterEach(function () {
    Ember.Test.adapter.exception = adapterException;
  });

  test('top-level not-found', async function (assert) {
    await visit('/404');
    assert
      .dom('[data-test-error-description]')
      .hasText(
        'Sorry, we were unable to find any content at that URL. Double check it or go back home.',
        'renders cluster error template'
      );
  });

  test('vault route not-found', async function (assert) {
    await visit('/vault/404');
    assert.dom('[data-test-not-found]').exists('renders the not found component');
  });

  test('cluster route not-found', async function (assert) {
    await visit('/vault/secrets/secret/404/show');
    assert.dom('[data-test-not-found]').exists('renders the not found component');
  });

  test('secret not-found', async function (assert) {
    await visit('/vault/secrets/secret/show/404');
    assert.dom('[data-test-secret-not-found]').exists('renders the message about the secret not being found');
  });
});
