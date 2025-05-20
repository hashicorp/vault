/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, fillIn, render, waitFor } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { AUTH_FORM } from 'vault/tests/helpers/auth/auth-form-selectors';
import { fillInLoginFields, VISIBLE_MOUNTS } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { CSP_ERROR } from 'vault/components/auth/page';

module('Integration | Component | auth | page', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.version = this.owner.lookup('service:version');
    this.cluster = { id: '1' };
    this.onAuthSuccess = sinon.spy();
    this.onNamespaceUpdate = sinon.spy();
    this.visibleAuthMounts = false;
    this.directLinkData = null;

    this.renderComponent = () => {
      return render(hbs`
        <Auth::Page
          @cluster={{this.cluster}}
          @directLinkData={{this.directLinkData}}
          @namespaceQueryParam={{this.nsQp}}
          @oidcProviderQueryParam={{this.providerQp}}
          @onAuthSuccess={{this.onAuthSuccess}}
          @onNamespaceUpdate={{this.onNamespaceUpdate}}
          @visibleAuthMounts={{this.visibleAuthMounts}}
        />
        `);
    };

    // in the real world more info is returned by auth requests
    // only including pertinent data for testing
    this.authRequest = (url) => this.server.post(url, () => ({ auth: { policies: ['default'] } }));
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

  test('it renders splash logo when oidc provider query param is present', async function (assert) {
    this.providerQp = 'myprovider';
    await this.renderComponent();
    assert.dom(AUTH_FORM.logo).exists();
    assert
      .dom(AUTH_FORM.helpText)
      .hasText(
        'Once you log in, you will be redirected back to your application. If you require login credentials, contact your administrator.'
      );
  });

  test('it disables namespace input when oidc provider query param is present', async function (assert) {
    this.providerQp = 'myprovider';
    this.version.features = ['Namespaces'];
    await this.renderComponent();
    assert.dom(AUTH_FORM.logo).exists();
    assert.dom(GENERAL.inputByAttr('namespace')).isDisabled();
  });

  test('it calls onNamespaceUpdate', async function (assert) {
    assert.expect(1);
    this.version.features = ['Namespaces'];
    await this.renderComponent();
    await fillIn(GENERAL.inputByAttr('namespace'), 'mynamespace');
    const [actual] = this.onNamespaceUpdate.lastCall.args;
    assert.strictEqual(actual, 'mynamespace', `onNamespaceUpdate called with: ${actual}`);
  });

  test('it calls onNamespaceUpdate for HVD managed clusters', async function (assert) {
    assert.expect(2);
    this.version.features = ['Namespaces'];
    this.owner.lookup('service:flags').featureFlags = ['VAULT_CLOUD_ADMIN_NAMESPACE'];
    this.nsQp = 'admin';
    await this.renderComponent();

    assert.dom(GENERAL.inputByAttr('namespace')).hasValue('');
    await fillIn(GENERAL.inputByAttr('namespace'), 'mynamespace');
    const [actual] = this.onNamespaceUpdate.lastCall.args;
    assert.strictEqual(actual, 'mynamespace', `onNamespaceUpdate called with: ${actual}`);
  });

  module('listing visibility', function (hooks) {
    hooks.beforeEach(function () {
      this.visibleAuthMounts = VISIBLE_MOUNTS;
      window.localStorage.clear();
    });

    test('it formats tab data if visible auth mounts exist', async function (assert) {
      await this.renderComponent();
      const expectedTabs = [
        { type: 'userpass', display: 'Userpass' },
        { type: 'oidc', display: 'OIDC' },
        { type: 'token', display: 'Token' },
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

    test('it selects type in the dropdown if @directLinkData references NON visible type', async function (assert) {
      this.directLinkData = { type: 'ldap', isVisibleMount: false };
      await this.renderComponent();
      assert.dom(GENERAL.selectByAttr('auth type')).hasValue('ldap', 'dropdown has type selected');
      assert.dom(AUTH_FORM.authForm('ldap')).exists();
      assert.dom(GENERAL.inputByAttr('username')).exists();
      assert.dom(GENERAL.inputByAttr('password')).exists();
      await click(AUTH_FORM.advancedSettings);
      assert.dom(GENERAL.inputByAttr('path')).exists();
      assert.dom(AUTH_FORM.tabBtn('ldap')).doesNotExist('tab does not render');
      assert
        .dom(GENERAL.backButton)
        .exists('back button renders because listing_visibility="unauth" for other mounts');
      assert
        .dom(GENERAL.buttonByAttr('other-methods'))
        .doesNotExist('"Sign in with other methods" does not render');
    });

    test('it renders single mount view instead of tabs if @directLinkData data references a visible type', async function (assert) {
      this.directLinkData = { path: 'my-oidc/', type: 'oidc', isVisibleMount: true };
      await this.renderComponent();
      assert.dom(AUTH_FORM.tabBtn('oidc')).hasText('OIDC', 'it renders tab for type');
      assert.dom(GENERAL.inputByAttr('role')).exists();
      assert.dom(GENERAL.inputByAttr('path')).hasAttribute('type', 'hidden');
      assert.dom(GENERAL.inputByAttr('path')).hasValue('my-oidc/');
      assert.dom(GENERAL.buttonByAttr('other-methods')).exists('"Sign in with other methods" renders');
      assert.dom(GENERAL.selectByAttr('auth type')).doesNotExist();
      assert.dom(AUTH_FORM.advancedSettings).doesNotExist();
      assert.dom(GENERAL.backButton).doesNotExist();
    });
  });

  const REQUEST_DATA = {
    username: {
      loginData: { username: 'matilda', password: 'password' },
      url: ({ path, username }) => `/auth/${path}/login/${username}`,
    },
    github: {
      loginData: { token: 'mysupersecuretoken' },
      url: ({ path }) => `/auth/${path}/login`,
    },
  };

  // only testing methods that submit via AuthForm (and not separate, child component)
  const AUTH_METHOD_TEST_CASES = [
    { authType: 'github', options: REQUEST_DATA.github },
    { authType: 'userpass', options: REQUEST_DATA.username },
    { authType: 'ldap', options: REQUEST_DATA.username },
    { authType: 'okta', options: REQUEST_DATA.username },
    { authType: 'radius', options: REQUEST_DATA.username },
  ];

  for (const { authType, options } of AUTH_METHOD_TEST_CASES) {
    test(`${authType}: it calls onAuthSuccess on submit for default path`, async function (assert) {
      assert.expect(1);
      const { loginData, url } = options;
      const requestUrl = url({ path: authType, username: loginData?.username });
      this.authRequest(requestUrl);

      await this.renderComponent();
      await fillIn(AUTH_FORM.selectMethod, authType);
      // await fillIn(AUTH_FORM.selectMethod, authType);
      await fillInLoginFields(loginData);
      await click(GENERAL.saveButton);
      const [actual] = this.onAuthSuccess.lastCall.args;
      const expected = {
        namespace: '',
        token: `vault-${authType}☃1`,
        isRoot: false,
      };
      assert.propEqual(actual, expected, `onAuthSuccess called with: ${JSON.stringify(actual)}`);
    });

    test(`${authType}: it calls onAuthSuccess on submit for custom path`, async function (assert) {
      assert.expect(1);
      const customPath = `${authType}-custom`;
      const { loginData, url } = options;
      const loginDataWithPath = { ...loginData, path: customPath };
      // pass custom path to request URL
      const requestUrl = url({ path: customPath, username: loginData?.username });
      this.authRequest(requestUrl);

      await this.renderComponent();
      await fillIn(AUTH_FORM.selectMethod, authType);
      // await fillIn(AUTH_FORM.selectMethod, authType);
      // toggle mount path input to specify custom path
      await fillInLoginFields(loginDataWithPath, { toggleOptions: true });
      await click(GENERAL.saveButton);

      const [actual] = this.onAuthSuccess.lastCall.args;
      const expected = {
        namespace: '',
        token: `vault-${authType}☃1`,
        isRoot: false,
      };
      assert.propEqual(actual, expected, `onAuthSuccess called with: ${JSON.stringify(actual)}`);
    });
  }

  // token makes a GET request so test separately
  test('token: it calls onAuthSuccess on submit', async function (assert) {
    assert.expect(1);
    this.server.get('/auth/token/lookup-self', () => {
      return { data: { policies: ['default'] } };
    });

    await this.renderComponent();
    await fillIn(AUTH_FORM.selectMethod, 'token');
    // await fillIn(AUTH_FORM.selectMethod, 'token');
    await fillInLoginFields({ token: 'mysupersecuretoken' });
    await click(GENERAL.saveButton);
    const [actual] = this.onAuthSuccess.lastCall.args;
    const expected = {
      namespace: '',
      token: `vault-token☃1`,
      isRoot: false,
    };
    assert.propEqual(actual, expected, `onAuthSuccess called with: ${JSON.stringify(actual)}`);
  });
});
