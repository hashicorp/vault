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
import { authRequest, fillInLoginFields } from 'vault/tests/helpers/auth/auth-helpers';

module('Integration | Component | auth | page', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.router = this.owner.lookup('service:router');
    this.auth = this.owner.lookup('service:auth');
    this.cluster = { id: '1' };
    this.authQp = 'token';
    // this.nsQp = 'admin';
    // this.providerQp = 'token';
    this.onAuthSuccess = sinon.spy();
    this.onNamespaceUpdate = sinon.spy();

    this.renderComponent = async () => {
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

  // TODO build out for all auth types
  const REQUEST_DATA = {
    token: {
      loginData: { token: 'mysupersecuretoken' },
      url: () => 'auth/token/lookup-self',
    },
    username: {
      loginData: { username: 'matilda', password: 'password' },
      url: ({ path, username }) => `/auth/${path}/login/${username}`,
    },
    github: {
      loginData: { token: 'mysupersecuretoken' },
      url: ({ path }) => `auth/${path}/login`,
    },
    oidc: {
      loginData: {},
      url: ({ path }) => `auth/${path}/oidc/auth_url`,
    },
    saml: {
      loginData: {},
      url: ({ path }) => `auth/${path}/sso_service_url`,
    },
  };
  const AUTH_METHOD_TEST_CASES = [
    { authType: 'userpass', options: REQUEST_DATA.username },
    { authType: 'ldap', options: REQUEST_DATA.username },
    { authType: 'okta', options: REQUEST_DATA.username },
    { authType: 'radius', options: REQUEST_DATA.username },
    { authType: 'github', options: REQUEST_DATA.github },
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
});
