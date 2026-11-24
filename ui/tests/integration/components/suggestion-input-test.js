/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, click, fillIn, typeIn } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';

const {
  searchSelect: { option, options },
  suggestion,
} = GENERAL;

module('Integration | Component | suggestion-input', function (hooks) {
  setupRenderingTest(hooks);

  module('KV type', function (hooks) {
    hooks.beforeEach(function () {
      this.label = 'Input something';
      this.mountPath = 'secret/';
      this.type = 'kv';
      this.apiStub = sinon
        .stub(this.owner.lookup('service:api').secrets, 'kvV2List')
        .resolves({ keys: ['foo/', 'my-secret'] });

      this.renderComponent = () =>
        render(hbs`
        <SuggestionInput
          @type={{this.type}}
          @label={{this.label}}
          @subText="Suggestions will display"
          @value={{this.secretPath}}
          @mountPath={{this.mountPath}}
          @onChange={{fn (mut this.secretPath)}}
        />
      `);
    });

    test('it should render label and sub text', async function (assert) {
      await this.renderComponent();
      assert.dom('label').hasText('Input something', 'Label renders');
      assert.dom('[data-test-label-subtext]').hasText('Suggestions will display', 'Label subtext renders');

      this.label = null;
      await this.renderComponent();
      assert.dom('label').doesNotExist('Label is hidden when arg is not provided');
    });

    test('it should disable input when mountPath is not provided', async function (assert) {
      this.mountPath = null;
      await this.renderComponent();
      assert.dom(suggestion.input('kv')).isDisabled('Input is disabled when mountPath is not provided');
    });

    test('it should fetch suggestions for initial mount path', async function (assert) {
      await this.renderComponent();
      assert.dom(options).doesNotExist('Suggestions are hidden initially');
      await click(suggestion.input('kv'));
      ['foo/', 'my-secret'].forEach((secret, index) => {
        assert.dom(option(index)).hasText(secret, 'Suggestion renders for initial mount path');
      });
    });

    test('it should fetch secrets and update suggestions on mountPath change', async function (assert) {
      this.apiStub.resolves({ keys: ['test1'] });
      this.mountPath = 'foo/';
      await this.renderComponent();
      await click(suggestion.input('kv'));
      assert.dom(option()).hasText('test1', 'Suggestions are fetched and render on mountPath change');
    });

    test('it should filter current result set', async function (assert) {
      await this.renderComponent();
      await click(suggestion.input('kv'));
      await typeIn(suggestion.input('kv'), 'sec');
      assert.dom(options).exists({ count: 1 }, 'Correct number of options render based on input value');
      assert.dom(option()).hasText('my-secret', 'Result set is filtered');
    });

    test('it should replace filter terms with full path to secret', async function (assert) {
      await this.renderComponent();
      await fillIn(suggestion.input('kv'), 'sec');
      await click(option());
      assert.dom(suggestion.input('kv')).hasValue('my-secret', 'Partial term replaced with selected secret');

      await fillIn(suggestion.input('kv'), '');
      this.apiStub.resolves({ keys: ['secret-nested', 'bar', 'baz'] });
      await click(option());
      await fillIn(suggestion.input('kv'), 'nest');
      await click(option());
      assert
        .dom(suggestion.input('kv'))
        .hasValue('foo/secret-nested', 'Partial term in nested path replaced with selected secret');
    });

    test('it should fetch secrets at nested paths', async function (assert) {
      await this.renderComponent();
      this.apiStub.resolves({ keys: ['bar/'] });
      await click(suggestion.input('kv'));
      await click(option());
      assert.dom(suggestion.input('kv')).hasValue('foo/', 'Input value updates on select');
      assert.dom(option()).hasText('bar/', 'Suggestions are fetched at new path');

      this.apiStub.resolves({ keys: ['baz/'] });
      await click(option());
      assert.dom(suggestion.input('kv')).hasValue('foo/bar/', 'Input value updates on select');
      assert.dom(option()).hasText('baz/', 'Suggestions are fetched at new path');

      this.apiStub.resolves({ keys: ['nested-secret'] });
      await typeIn(suggestion.input('kv'), 'baz/');
      assert.dom(option()).hasText('nested-secret', 'Suggestions are fetched at new path');
    });

    test('it should only render dropdown when suggestions exist', async function (assert) {
      await this.renderComponent();
      await click(suggestion.input('kv'));
      assert.dom(options).exists({ count: 2 }, 'Suggestions render');

      await fillIn(suggestion.input('kv'), 'bar');
      assert.dom(options).doesNotExist('Drop down is hidden when there are no suggestions');

      await fillIn(suggestion.input('kv'), '');
      assert.dom(options).exists({ count: 2 }, 'Suggestions render');

      this.apiStub.resolves({ keys: [] });
      await click(option());
      assert.dom(options).doesNotExist('Drop down is hidden when there are no suggestions');
    });
  });

  module('Database type', function (hooks) {
    hooks.beforeEach(function () {
      this.label = 'Select a static role';
      this.mountPath = 'database';
      this.type = 'database';
      this.apiStub = sinon
        .stub(this.owner.lookup('service:api').secrets, 'databaseListStaticRoles')
        .resolves({ keys: ['role-1', 'role-2', 'test-role'] });

      this.renderComponent = () =>
        render(hbs`
        <SuggestionInput
          @type={{this.type}}
          @label={{this.label}}
          @subText="Select a role"
          @value={{this.roleName}}
          @mountPath={{this.mountPath}}
          @onChange={{fn (mut this.roleName)}}
        />
      `);
    });

    test('it should render label and sub text', async function (assert) {
      await this.renderComponent();
      assert.dom('label').hasText('Select a static role', 'Label renders');
      assert.dom('[data-test-label-subtext]').hasText('Select a role', 'Label subtext renders');
    });

    test('it should disable input when mountPath is not provided', async function (assert) {
      this.mountPath = null;
      await this.renderComponent();
      assert.dom(suggestion.input('database')).isDisabled('Input is disabled when mountPath is not provided');
    });

    test('it should fetch database static roles for initial mount', async function (assert) {
      await this.renderComponent();
      assert.dom(options).doesNotExist('Suggestions are hidden initially');
      await click(suggestion.input('database'));
      ['static-roles/role-1', 'static-roles/role-2', 'static-roles/test-role'].forEach((role, index) => {
        assert.dom(option(index)).hasText(role, 'Role suggestion renders');
      });
    });

    test('it should filter current result set', async function (assert) {
      await this.renderComponent();
      await click(suggestion.input('database'));
      await typeIn(suggestion.input('database'), 'test');
      assert.dom(options).exists({ count: 1 }, 'Correct number of options render based on input value');
      assert.dom(option()).hasText('static-roles/test-role', 'Result set is filtered');
    });

    test('it should set selected role as value', async function (assert) {
      await this.renderComponent();
      await fillIn(suggestion.input('database'), 'test');
      await click(option());
      assert
        .dom(suggestion.input('database'))
        .hasValue('static-roles/test-role', 'Selected role is set as value');
    });

    test('it should only render dropdown when suggestions exist', async function (assert) {
      await this.renderComponent();
      await click(suggestion.input('database'));
      assert.dom(options).exists({ count: 3 }, 'Suggestions render');

      await fillIn(suggestion.input('database'), 'nonexistent');
      assert.dom(options).doesNotExist('Drop down is hidden when there are no suggestions');

      await fillIn(suggestion.input('database'), '');
      assert.dom(options).exists({ count: 3 }, 'Suggestions render again');
    });

    test('it should handle empty role list', async function (assert) {
      this.apiStub.resolves({ keys: [] });
      await this.renderComponent();
      await click(suggestion.input('database'));
      assert.dom(options).doesNotExist('No suggestions when API returns empty list');
    });
  });
});
