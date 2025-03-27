/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, fillIn, render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { AUTH_FORM } from 'vault/tests/helpers/auth/auth-form-selectors';
import { fillInLoginFields } from 'vault/tests/helpers/auth/auth-helpers';

module('Integration | Component | auth | page', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.version = this.owner.lookup('service:version');
    this.cluster = { id: '1' };
    this.onAuthSuccess = sinon.spy();
    this.onNamespaceUpdate = sinon.spy();

    this.renderComponent = () => {
      return render(hbs`
        <Auth::Page
          @authMethodQueryParam={{this.authQp}}
          @cluster={{this.cluster}}
          @namespaceQueryParam={{this.nsQp}}
          @oidcProviderQueryParam={{this.providerQp}}
          @onAuthSuccess={{this.onAuthSuccess}}
          @onNamespaceUpdate={{this.onNamespaceUpdate}}
          @wrappedToken={{this.wrappedToken}}
        />
        `);
    };

    // in the real world more info is returned by auth requests
    // only including pertinent data for testing
    this.authRequest = (url) => this.server.post(url, () => ({ auth: { policies: ['default'] } }));
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
    assert.dom(AUTH_FORM.namespaceInput).isDisabled();
  });

  test('it calls onNamespaceUpdate', async function (assert) {
    assert.expect(1);
    this.version.features = ['Namespaces'];
    await this.renderComponent();
    await fillIn(AUTH_FORM.namespaceInput, 'mynamespace');
    const [actual] = this.onNamespaceUpdate.lastCall.args;
    assert.strictEqual(actual, 'mynamespace', `onNamespaceUpdate called with: ${actual}`);
  });

  test('it calls onNamespaceUpdate for HVD managed clusters', async function (assert) {
    assert.expect(2);
    this.version.features = ['Namespaces'];
    this.owner.lookup('service:flags').featureFlags = ['VAULT_CLOUD_ADMIN_NAMESPACE'];
    this.nsQp = 'admin';
    await this.renderComponent();

    assert.dom(AUTH_FORM.namespaceInput).hasValue('');
    await fillIn(AUTH_FORM.namespaceInput, 'mynamespace');
    const [actual] = this.onNamespaceUpdate.lastCall.args;
    assert.strictEqual(actual, 'mynamespace', `onNamespaceUpdate called with: ${actual}`);
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
      await fillIn(AUTH_FORM.method, authType);
      await fillInLoginFields(loginData);
      await click(AUTH_FORM.login);
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
      const loginDataWithPath = { ...loginData, 'auth-form-mount-path': customPath };
      // pass custom path to request URL
      const requestUrl = url({ path: customPath, username: loginData?.username });
      this.authRequest(requestUrl);

      await this.renderComponent();
      await fillIn(AUTH_FORM.method, authType);
      // toggle mount path input to specify custom path
      await fillInLoginFields(loginDataWithPath, { toggleOptions: true });
      await click(AUTH_FORM.login);

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
    await fillIn(AUTH_FORM.method, 'token');
    await fillInLoginFields({ token: 'mysupersecuretoken' });
    await click(AUTH_FORM.login);
    const [actual] = this.onAuthSuccess.lastCall.args;
    const expected = {
      namespace: '',
      token: `vault-token☃1`,
      isRoot: false,
    };
    assert.propEqual(actual, expected, `onAuthSuccess called with: ${JSON.stringify(actual)}`);
  });
});
