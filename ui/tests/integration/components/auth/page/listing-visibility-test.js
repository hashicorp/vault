/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { AUTH_FORM } from 'vault/tests/helpers/auth/auth-form-selectors';
import { click, fillIn, waitFor } from '@ember/test-helpers';
import { fillInLoginFields, SYS_INTERNAL_UI_MOUNTS } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { module, test } from 'qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { setupRenderingTest } from 'ember-qunit';
import { setupTotpMfaResponse } from 'vault/tests/helpers/mfa/mfa-helpers';
import setupTestContext from './setup-test-context';
import sinon from 'sinon';

module('Integration | Component | auth | page | listing visibility', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    setupTestContext(this);
    this.visibleAuthMounts = SYS_INTERNAL_UI_MOUNTS;
    // extra setup for when the "oidc" is selected and the oidc-jwt component renders
    this.routerStub = sinon.stub(this.owner.lookup('service:router'), 'urlFor').returns('123-example.com');
  });

  hooks.afterEach(function () {
    this.routerStub.restore();
  });

  test('it formats and renders tabs if visible auth mounts exist', async function (assert) {
    await this.renderComponent();
    const expectedTabs = [
      { type: 'userpass', display: 'Userpass' },
      { type: 'oidc', display: 'OIDC' },
      { type: 'ldap', display: 'LDAP' },
    ];

    assert.dom(GENERAL.selectByAttr('auth type')).doesNotExist('dropdown does not render');
    // there are 4 mount paths returned in visibleAuthMounts above,
    // but two are of the same type so only expect 3 tabs
    assert.dom(AUTH_FORM.tabs).exists({ count: 3 }, 'it groups mount paths by type and renders 3 tabs');
    expectedTabs.forEach((m) => {
      assert.dom(AUTH_FORM.tabBtn(m.type)).exists(`${m.type} renders as a tab`);
      assert.dom(AUTH_FORM.tabBtn(m.type)).hasText(m.display, `${m.type} renders expected display name`);
    });
    assert
      .dom(AUTH_FORM.tabBtn('userpass'))
      .hasAttribute('aria-selected', 'true', 'it selects the first type by default');
  });

  test('it renders dropdown as alternate view', async function (assert) {
    await this.renderComponent();
    assert.dom(AUTH_FORM.tabs).exists({ count: 3 }, 'tabs render by default');
    assert.dom(GENERAL.backButton).doesNotExist();
    await click(GENERAL.button('Sign in with other methods'));
    assert
      .dom(GENERAL.button('Sign in with other methods'))
      .doesNotExist('button disappears after it is clicked');
    assert.dom(GENERAL.selectByAttr('auth type')).hasValue('userpass', 'dropdown has userpass selected');
    assert.dom(AUTH_FORM.advancedSettings).exists('toggle renders even though userpass has visible mounts');
    await click(AUTH_FORM.advancedSettings);
    assert.dom(GENERAL.inputByAttr('path')).exists({ count: 1 });
    assert.dom(GENERAL.inputByAttr('path')).hasValue('', 'it renders empty custom path input');

    await fillIn(GENERAL.selectByAttr('auth type'), 'oidc');
    assert.dom(AUTH_FORM.advancedSettings).exists('toggle renders even though oidc has a visible mount');
    await click(AUTH_FORM.advancedSettings);
    assert.dom(GENERAL.inputByAttr('path')).exists({ count: 1 });
    assert.dom(GENERAL.inputByAttr('path')).hasValue('', 'it renders empty custom path input');
    await click(GENERAL.backButton);
    assert.dom(GENERAL.backButton).doesNotExist('"Back" button does not render after it is clicked');
    assert.dom(AUTH_FORM.tabs).exists({ count: 3 }, 'clicking "Back" renders tabs again');
    assert
      .dom(GENERAL.button('Sign in with other methods'))
      .exists('"Sign in with other methods" renders again');
  });

  // integration tests for ?with= query param
  module('with a direct link', function (hooks) {
    hooks.beforeEach(function () {
      // if path exists, the mount has listing_visibility="unauth"
      this.directLinkIsVisibleMount = { path: 'my_oidc/', type: 'oidc' };
      this.directLinkIsJustType = { type: 'okta' };
    });

    test('it selects type in the dropdown if direct link is just type', async function (assert) {
      this.directLinkData = this.directLinkIsJustType;
      await this.renderComponent();
      assert.dom(AUTH_FORM.tabBtn('okta')).doesNotExist('tab does not render');
      assert.dom(GENERAL.selectByAttr('auth type')).hasValue('okta', 'dropdown has type selected');
      assert.dom(AUTH_FORM.authForm('okta')).exists();
      assert.dom(GENERAL.inputByAttr('username')).exists();
      assert.dom(GENERAL.inputByAttr('password')).exists();
      await click(AUTH_FORM.advancedSettings);
      assert.dom(GENERAL.inputByAttr('path')).exists({ count: 1 });
      assert
        .dom(GENERAL.button('Sign in with other methods'))
        .doesNotExist('"Sign in with other methods" does not render');
      assert.dom(GENERAL.backButton).exists('back button renders because tabs exist for other methods');
      await click(GENERAL.backButton);
      assert
        .dom(AUTH_FORM.tabBtn('userpass'))
        .hasAttribute('aria-selected', 'true', 'first tab is selected on back');
    });

    test('it renders single method view instead of tabs if direct link includes path', async function (assert) {
      this.directLinkData = this.directLinkIsVisibleMount;
      await this.renderComponent();
      assert.dom(AUTH_FORM.authForm('oidc')).exists;
      assert.dom(AUTH_FORM.tabBtn('oidc')).hasText('OIDC', 'it renders tab for type');
      assert.dom(AUTH_FORM.tabs).exists({ count: 1 }, 'only one tab renders');
      assert.dom(GENERAL.inputByAttr('role')).exists();
      assert.dom(GENERAL.inputByAttr('path')).hasAttribute('type', 'hidden');
      assert.dom(GENERAL.inputByAttr('path')).hasValue('my_oidc/');
      assert.dom(GENERAL.button('Sign in with other methods')).exists('"Sign in with other methods" renders');
      assert.dom(GENERAL.selectByAttr('auth type')).doesNotExist();
      assert.dom(AUTH_FORM.advancedSettings).doesNotExist();
      assert.dom(GENERAL.backButton).doesNotExist();
    });

    test('it prioritizes auth type from canceled mfa instead of direct link for path', async function (assert) {
      assert.expect(1);
      this.directLinkData = this.directLinkIsVisibleMount; // type is "oidc"
      const authType = 'okta'; // set to a type that differs from direct link
      this.server.post(`/auth/okta/login/matilda`, () => setupTotpMfaResponse(authType));
      await this.renderComponent();
      await click(GENERAL.button('Sign in with other methods'));
      await fillIn(AUTH_FORM.selectMethod, authType);
      await fillInLoginFields({ username: 'matilda', password: 'password' });
      await click(GENERAL.submitButton);
      await waitFor('[data-test-mfa-description]'); // wait until MFA validation renders
      await click(GENERAL.cancelButton);
      assert.dom(AUTH_FORM.selectMethod).hasValue(authType, 'Okta is selected in dropdown');
    });

    test('it prioritizes auth type from canceled mfa instead of direct link with type', async function (assert) {
      assert.expect(1);
      this.directLinkData = this.directLinkIsJustType; // type is "okta"
      const authType = 'userpass'; // set to a type that differs from direct link
      this.server.post(`/auth/userpass/login/matilda`, () => setupTotpMfaResponse(authType));
      await this.renderComponent();
      await fillIn(AUTH_FORM.selectMethod, authType);
      await fillInLoginFields({ username: 'matilda', password: 'password' });
      await click(GENERAL.submitButton);
      await waitFor('[data-test-mfa-description]'); // wait until MFA validation renders
      await click(GENERAL.cancelButton);
      assert.dom(AUTH_FORM.tabBtn('userpass')).hasAttribute('aria-selected', 'true');
    });
  });
});
