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
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { AUTH_FORM, MFA_SELECTORS } from 'vault/tests/helpers/auth/auth-form-selectors';
import { fillInLoginFields } from 'vault/tests/helpers/auth/auth-helpers';
import { setupTotpMfaResponse } from 'vault/tests/helpers/auth/mfa-helpers';

// in the real world more info is returned by auth requests
// only including pertinent data for testing
const authRequest = (context, options) => {
  const { isMfa = false, authMountPath = '', url = '' } = options;
  return context.server.post(url, () => {
    if (isMfa) {
      return {
        warnings: [
          'A login request was issued that is subject to MFA validation. Please make sure to validate the login by sending another request to mfa/validate endpoint.',
        ],
        ...setupTotpMfaResponse(authMountPath),
      };
    }
    return {
      auth: { policies: ['default'] },
    };
  });
};

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
    token: {
      loginData: { token: 'mysupersecuretoken' },
      url: () => '/auth/token/lookup-self',
    },
    username: {
      loginData: { username: 'matilda', password: 'password' },
      url: ({ path, username }) => `/auth/${path}/login/${username}`,
    },
    github: {
      loginData: { token: 'mysupersecuretoken' },
      url: ({ path }) => `/auth/${path}/login`,
    },
    // oidc: {
    //   loginData: { role: 'some-dev' },
    //   url: ({ path }) => `/auth/${path}/oidc/auth_url`,
    //   responseType: 'oidc',
    // },
    // saml: {
    //   loginData: { role: 'some-dev' },
    //   url: ({ path }) => `/auth/${path}/sso_service_url`,
    // },
  };

  const AUTH_METHOD_TEST_CASES = [
    { authType: 'github', options: REQUEST_DATA.github },
    //input username + password
    { authType: 'userpass', options: REQUEST_DATA.username },
    { authType: 'ldap', options: REQUEST_DATA.username },
    { authType: 'okta', options: REQUEST_DATA.username },
    { authType: 'radius', options: REQUEST_DATA.username },
    // TODO CMB add these tests cases when login logic is standardized (currently login success/mfa is tested by the individual components)
    // { authType: 'oidc', options: REQUEST_DATA.oidc },
    // { authType: 'jwt', options: REQUEST_DATA.oidc },
    // { authType: 'saml', options: REQUEST_DATA.saml },
  ];

  for (const { authType, options } of AUTH_METHOD_TEST_CASES) {
    test(`${authType}: it calls onAuthSuccess on submit for default path`, async function (assert) {
      assert.expect(1);
      const { loginData, url } = options;
      const requestUrl = url({ path: authType, username: loginData?.username });
      authRequest(this, { url: requestUrl });

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
      authRequest(this, { url: requestUrl });

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

    test(`${authType}: it should display mfa requirement for default path`, async function (assert) {
      const { loginData, url } = options;
      assert.expect(3 + Object.keys(loginData).length);

      const requestUrl = url({ path: authType, username: loginData?.username });
      // authMountPath necessary to return mfa_constraints
      authRequest(this, { isMfa: true, authMountPath: authType, url: requestUrl });

      await this.renderComponent();
      await fillIn(AUTH_FORM.method, authType);
      await fillInLoginFields(loginData);
      await click(AUTH_FORM.login);

      assert
        .dom(MFA_SELECTORS.mfaForm)
        .hasText(
          'Multi-factor authentication is enabled for your account. Enter your authentication code to log in. TOTP passcode Verify'
        );
      await click(GENERAL.backButton);
      assert.dom(AUTH_FORM.form).exists('clicking back returns to auth form');
      assert.dom(GENERAL.selectByAttr('auth-method')).hasValue(authType, 'preserves method type on back');

      for (const field of Object.keys(loginData)) {
        assert.dom(AUTH_FORM.input(field)).hasValue('', `${field} input clears on back`);
      }
    });

    test(`${authType}: it should display mfa requirement for custom path`, async function (assert) {
      const customPath = `${authType}-custom`;
      const { loginData, url } = options;
      assert.expect(3 + Object.keys(loginData).length);

      const loginDataWithPath = { ...loginData, 'auth-form-mount-path': customPath };
      // pass custom path to request URL
      const requestUrl = url({ path: customPath, username: loginData?.username });
      // authMountPath necessary to return mfa_constraints
      authRequest(this, { isMfa: true, authMountPath: customPath, url: requestUrl });

      await this.renderComponent();
      await fillIn(AUTH_FORM.method, authType);
      await fillInLoginFields(loginDataWithPath, { toggleOptions: true });
      await click(AUTH_FORM.login);

      assert
        .dom(MFA_SELECTORS.mfaForm)
        .hasText(
          'Multi-factor authentication is enabled for your account. Enter your authentication code to log in. TOTP passcode Verify'
        );
      await click(GENERAL.backButton);
      assert.dom(AUTH_FORM.form).exists('clicking back returns to auth form');
      assert.dom(GENERAL.selectByAttr('auth-method')).hasValue(authType, 'preserves method type on back');

      for (const field of Object.keys(loginData)) {
        assert.dom(AUTH_FORM.input(field)).hasValue('', `${field} input clears on back`);
      }
    });
  }

  // token makes a GET request so test separately
  test('token: it calls onAuthSuccess on submit', async function (assert) {
    assert.expect(1);
    const authType = 'token';
    const { loginData, url } = REQUEST_DATA.token;
    this.server.get(url(), () => {
      return { data: { policies: ['default'] } };
    });

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
});
