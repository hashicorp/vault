/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import hbs from 'htmlbars-inline-precompile';
import { findAll, render } from '@ember/test-helpers';
import { capitalize } from '@ember/string';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | auth | fields', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.loginFields = ['username', 'role', 'token', 'password'];
    this.renderComponent = () => {
      return render(hbs`<Auth::Fields @loginFields={{this.loginFields}} />`);
    };
  });

  test('it renders field labels', async function (assert) {
    await this.renderComponent();
    const labels = findAll('label');
    this.loginFields.forEach((field) => {
      const label = labels.find((l) => l.innerText === capitalize(field));
      assert.dom(label).exists(`${field}: it renders capitalized field label`);
    });
  });

  test('it renders password input types for token and password fields', async function (assert) {
    await this.renderComponent();
    assert.dom(GENERAL.inputByAttr('token')).hasAttribute('type', 'password');
    assert.dom(GENERAL.inputByAttr('password')).hasAttribute('type', 'password');
  });

  test('it renders text input types for other fields', async function (assert) {
    await this.renderComponent();
    assert.dom(GENERAL.inputByAttr('username')).hasAttribute('type', 'text');
    assert.dom(GENERAL.inputByAttr('role')).hasAttribute('type', 'text');
  });

  test('it renders expected autocomplete values', async function (assert) {
    await this.renderComponent();
    const expectedValues = {
      username: 'username',
      role: 'role',
      token: 'off',
      password: 'current-password',
    };
    this.loginFields.forEach((field) => {
      const expected = expectedValues[field];
      assert
        .dom(GENERAL.inputByAttr(field))
        .hasAttribute('autocomplete', expected, `${field}: it renders autocomplete value "${expected}"`);
    });
  });
});
