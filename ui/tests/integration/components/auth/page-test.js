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
  const USERNAME_PASSWORD_METHODS = ['userpass', 'ldap', 'okta', 'radius'];
  // const TOKEN_METHODS = ['token', 'github'];
  // BASE_LOGIN_METHODS.filter(
  //   (m) => m.formAttributes.includes('username') && m.formAttributes.includes('password')
  // ).map((m) => m.type);

  for (const authType of USERNAME_PASSWORD_METHODS) {
    test(`${authType} it calls onAuthSuccess on login for default path`, async function (assert) {
      assert.expect(1);
      const loginData = { username: 'matilda', password: 'password' };
      const methodData = { authType, authMountPath: authType };

      authRequest(this, { ...methodData, username: loginData.username });

      await this.renderComponent();
      await fillIn(AUTH_FORM.method, authType);
      await fillInLoginFields(loginData, { authType });
      await click(AUTH_FORM.login);

      const [actual] = this.onAuthSuccess.lastCall.args;
      const expected = {
        namespace: '',
        token: `vault-${authType}☃1`,
        isRoot: false,
      };
      assert.propEqual(actual, expected, `onAuthSuccess called with: ${JSON.stringify(actual)}`);
    });

    test(`${authType}: it calls onAuthSuccess on login for custom path`, async function (assert) {
      assert.expect(1);
      const customPath = `${authType}-custom`;
      const loginData = { username: 'matilda', password: 'password', 'auth-form-mount-path': customPath };
      const methodData = { authType, authMountPath: customPath };

      authRequest(this, { ...methodData, username: loginData.username });

      await this.renderComponent();
      await fillIn(AUTH_FORM.method, authType);
      await fillInLoginFields(loginData, { authType, toggleOptions: true });
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
      assert.expect(5);
      const loginData = { username: 'matilda', password: 'password' };
      const methodData = { authType, authMountPath: authType, isMfa: true };

      authRequest(this, { ...methodData, username: loginData.username });

      await this.renderComponent();
      await fillIn(AUTH_FORM.method, authType);
      await fillInLoginFields(loginData, { authType });
      await click(AUTH_FORM.login);

      assert
        .dom(MFA_SELECTORS.mfaForm)
        .hasText(
          'Multi-factor authentication is enabled for your account. Enter your authentication code to log in. TOTP passcode Verify'
        );
      await click(GENERAL.backButton);
      assert.dom(AUTH_FORM.form).exists('clicking back returns to auth form');
      assert.dom(GENERAL.selectByAttr('auth-method')).hasValue(authType, 'preserves method type on back');

      assert.dom(AUTH_FORM.input('username')).hasValue('', 'clears username on back');
      assert.dom(AUTH_FORM.input('password')).hasValue('', 'clears password on back');
    });

    test(`${authType}: it should display mfa requirement for custom path`, async function (assert) {
      assert.expect(5);
      const customPath = `${authType}-custom`;
      const loginData = { username: 'matilda', password: 'password', 'auth-form-mount-path': customPath };
      const methodData = { authType, authMountPath: customPath, isMfa: true };

      authRequest(this, { ...methodData, username: loginData.username });

      await this.renderComponent();
      await fillIn(AUTH_FORM.method, authType);
      await fillInLoginFields(loginData, { authType, toggleOptions: true });
      await click(AUTH_FORM.login);

      assert
        .dom(MFA_SELECTORS.mfaForm)
        .hasText(
          'Multi-factor authentication is enabled for your account. Enter your authentication code to log in. TOTP passcode Verify'
        );
      await click(GENERAL.backButton);
      assert.dom(AUTH_FORM.form).exists('clicking back returns to auth form');
      assert.dom(GENERAL.selectByAttr('auth-method')).hasValue(authType, 'preserves method type on back');

      assert.dom(AUTH_FORM.input('username')).hasValue('', 'clears username on back');
      assert.dom(AUTH_FORM.input('password')).hasValue('', 'clears password on back');
    });
  }
});
