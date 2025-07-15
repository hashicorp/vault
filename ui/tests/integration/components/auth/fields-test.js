/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import hbs from 'htmlbars-inline-precompile';
import { find, render } from '@ember/test-helpers';
import { capitalize } from '@ember/string';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | auth | fields', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.loginFields = [
      { name: 'username' },
      { name: 'role', helperText: 'Wow neat role!' },
      { name: 'token', label: 'Super secret token' },
      { name: 'password' },
    ];
    this.renderComponent = () => {
      return render(hbs`<Auth::Fields @loginFields={{this.loginFields}} />`);
    };
  });

  test('it renders field name as input label if "label" key is not specified', async function (assert) {
    await this.renderComponent();
    for (const field of ['username', 'password', 'role']) {
      const id = find(GENERAL.inputByAttr(field)).id;
      assert
        .dom(`#label-${id}`)
        .hasText(capitalize(field), `${field} it renders name if "label" key is not present`);
    }
  });

  test('it does NOT render "helperText" if not present', async function (assert) {
    await this.renderComponent();
    for (const field of ['username', 'password', 'token']) {
      const id = find(GENERAL.inputByAttr(field)).id;
      assert
        .dom(`#helper-text-${id}`)
        .doesNotExist(`${field}: it does not render helperText if key is not present`);
    }
  });

  test('it renders "helperText" if specified', async function (assert) {
    await this.renderComponent();
    const id = find(GENERAL.inputByAttr('role')).id;
    assert.dom(`#helper-text-${id}`).hasText('Wow neat role!');
  });

  test('it renders "label" if specified', async function (assert) {
    await this.renderComponent();
    const id = find(GENERAL.inputByAttr('token')).id;
    assert.dom(`#label-${id}`).hasText('Super secret token', 'it renders "label" instead of "name"');
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
      username: 'off',
      role: 'off',
      token: 'off',
      password: 'off',
    };

    for (const field of this.loginFields) {
      const { name } = field;
      const expected = expectedValues[name];
      assert
        .dom(GENERAL.inputByAttr(name))
        .hasAttribute('autocomplete', expected, `${name}: it renders autocomplete value "${expected}"`);
    }
  });
});
