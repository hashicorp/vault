/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { fillIn, find, render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { AUTH_FORM } from 'vault/tests/helpers/auth/auth-form-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | auth | namespace input', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.disabled = false;
    this.oidcProviderQueryParam = '';
    this.handleNamespaceUpdate = sinon.spy();
    this.shouldRefocusNamespaceInput = false;

    this.renderComponent = () => {
      return render(hbs`
      <Auth::NamespaceInput
        @disabled={{this.disabled}}
        @handleNamespaceUpdate={{this.handleNamespaceUpdate}}
        @namespaceQueryParam={{this.namespaceQueryParam}}
        @shouldRefocusNamespaceInput={{this.shouldRefocusNamespaceInput}}
      />`);
    };
  });

  test('it fires @handleNamespaceUpdate callback', async function (assert) {
    assert.expect(1);
    await this.renderComponent();
    await fillIn(GENERAL.inputByAttr('namespace'), 'ns-1');
    const [actual] = this.handleNamespaceUpdate.lastCall.args;
    assert.strictEqual(actual, 'ns-1', `handleNamespaceUpdate called with: ${actual}`);
  });

  test('it disables the input if @disabled is true', async function (assert) {
    this.disabled = true;
    await this.renderComponent();
    assert.dom(GENERAL.inputByAttr('namespace')).isDisabled();
  });

  test('it does not focus the input if @shouldRefocusNamespaceInput is false', async function (assert) {
    await this.renderComponent();
    const element = find(GENERAL.inputByAttr('namespace'));
    assert.notStrictEqual(document.activeElement, element, 'the namespace input is NOT focused');
  });

  test('it focuses the input if @shouldRefocusNamespaceInput is true', async function (assert) {
    this.shouldRefocusNamespaceInput = true;
    await this.renderComponent();
    const element = find(GENERAL.inputByAttr('namespace'));
    assert.strictEqual(document.activeElement, element, 'the namespace input is focused');
  });

  module('HVD managed', function (hooks) {
    hooks.beforeEach(function () {
      this.owner.lookup('service:flags').featureFlags = ['VAULT_CLOUD_ADMIN_NAMESPACE'];
    });

    test('it sets namespace', async function (assert) {
      this.namespaceQueryParam = 'admin/west-coast';
      await this.renderComponent();
      assert.dom(AUTH_FORM.managedNsRoot).hasValue('/admin');
      assert.dom(AUTH_FORM.managedNsRoot).hasAttribute('readonly');
      assert.dom(GENERAL.inputByAttr('namespace')).hasValue('/west-coast');
    });

    test('it calls onNamespaceUpdate', async function (assert) {
      assert.expect(2);
      this.namespaceQueryParam = 'admin';
      await this.renderComponent();

      assert.dom(GENERAL.inputByAttr('namespace')).hasValue('');
      await fillIn(GENERAL.inputByAttr('namespace'), 'ns-1');
      const [actual] = this.handleNamespaceUpdate.lastCall.args;
      assert.strictEqual(actual, 'ns-1', `handleNamespaceUpdate called with: ${actual}`);
    });
  });
});
