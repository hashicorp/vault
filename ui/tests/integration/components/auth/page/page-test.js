/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, fillIn, waitFor } from '@ember/test-helpers';
import { AUTH_FORM } from 'vault/tests/helpers/auth/auth-form-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { CSP_ERROR } from 'vault/components/auth/page';
import setupTestContext from './setup-test-context';

/*
The AuthPage parents much of the authentication workflow and so it can be used to test lots of auth functionality.
This file tests the base component functionality. The other files test method authentication, listing visibility, 
login settings (enterprise feature), and mfa.
*/
module('Integration | Component | auth | page', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    setupTestContext(this);
  });

  test('it renders error on CSP violation', async function (assert) {
    assert.expect(2);
    this.cluster.standby = true;
    await this.renderComponent();
    assert.dom(GENERAL.pageError.error).doesNotExist();
    this.owner.lookup('service:csp-event').handleEvent({ violatedDirective: 'connect-src' });
    await waitFor(GENERAL.pageError.error);
    assert.dom(GENERAL.pageError.error).hasText(CSP_ERROR);
  });

  test('it renders splash logo and disables namespace input when oidc provider query param is present', async function (assert) {
    this.oidcProviderQueryParam = 'myprovider';
    this.version.features = ['Namespaces'];
    await this.renderComponent();
    assert.dom(AUTH_FORM.logo).exists();
    assert.dom(GENERAL.inputByAttr('namespace')).isDisabled();
    assert
      .dom(AUTH_FORM.helpText)
      .hasText(
        'Once you log in, you will be redirected back to your application. If you require login credentials, contact your administrator.'
      );
  });

  test('it calls onNamespaceUpdate', async function (assert) {
    assert.expect(1);
    this.version.features = ['Namespaces'];
    await this.renderComponent();
    await fillIn(GENERAL.inputByAttr('namespace'), 'mynamespace');
    const [actual] = this.onNamespaceUpdate.lastCall.args;
    assert.strictEqual(actual, 'mynamespace', `onNamespaceUpdate called with: ${actual}`);
  });

  test('it passes query param to namespace input', async function (assert) {
    this.version.features = ['Namespaces'];
    this.namespaceQueryParam = 'ns-1';
    await this.renderComponent();
    assert.dom(GENERAL.inputByAttr('namespace')).hasValue(this.namespaceQueryParam);
  });

  test('it does not render the namespace input on community', async function (assert) {
    this.version.type = 'community';
    this.version.features = [];
    await this.renderComponent();
    assert.dom(GENERAL.inputByAttr('namespace')).doesNotExist();
  });

  test('it does not render the namespace input on enterprise without the "Namespaces" feature', async function (assert) {
    this.version.type = 'enterprise';
    this.version.features = [];
    await this.renderComponent();
    assert.dom(GENERAL.inputByAttr('namespace')).doesNotExist();
  });

  test('it selects type in the dropdown if direct link just has type', async function (assert) {
    this.directLinkData = { type: 'oidc' };
    await this.renderComponent();
    assert.dom(AUTH_FORM.tabBtn('oidc')).doesNotExist('tab does not render');
    assert.dom(GENERAL.selectByAttr('auth type')).hasValue('oidc', 'dropdown has type selected');
    assert.dom(AUTH_FORM.authForm('oidc')).exists();
    assert.dom(GENERAL.inputByAttr('role')).exists();
    await click(AUTH_FORM.advancedSettings);
    assert.dom(GENERAL.inputByAttr('path')).exists({ count: 1 });
    assert.dom(GENERAL.backButton).doesNotExist();
    assert
      .dom(GENERAL.button('Sign in with other methods'))
      .doesNotExist('"Sign in with other methods" does not render');
  });
});
