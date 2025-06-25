/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { AUTH_FORM } from 'vault/tests/helpers/auth/auth-form-selectors';
import { click, fillIn, waitUntil } from '@ember/test-helpers';
import { ERROR_JWT_LOGIN } from 'vault/components/auth/form/oidc-jwt';
import { fillInLoginFields } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { module, test } from 'qunit';
import { overrideResponse } from 'vault/tests/helpers/stubs';
import { RESPONSE_STUBS, TOKEN_DATA } from 'vault/tests/helpers/auth/response-stubs';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { setupRenderingTest } from 'ember-qunit';
import { triggerMessageEvent, windowStub } from 'vault/tests/helpers/oidc-window-stub';
import setupTestContext from './setup-test-context';
import sinon from 'sinon';

const methodAuthenticationTests = (test) => {
  test('it sets token data on login for default path', async function (assert) {
    assert.expect(5);
    // Setup
    this.stubRequests();
    // Render and log in
    await this.renderComponent();
    await fillIn(AUTH_FORM.selectMethod, this.authType);
    await fillInLoginFields(this.loginData);
    if (this.authType === 'oidc') {
      triggerMessageEvent(this.path);
    }
    await click(GENERAL.submitButton);
    await waitUntil(() => this.setTokenDataSpy.calledOnce);
    const [tokenName, persistedTokenData] = this.setTokenDataSpy.lastCall.args;

    const expectedData = {
      ...TOKEN_DATA[this.authType],
      // there are other tests that confirm this calculation happens as expected, just copy value from spy
      tokenExpirationEpoch: persistedTokenData.tokenExpirationEpoch,
    };

    assert.strictEqual(tokenName, this.tokenName, 'setTokenData is called with expected token name');
    assert.propEqual(persistedTokenData, expectedData, 'setTokenData is called with expected data');

    // propEqual failures are challenging to parse in CI so pulling out a couple of important attrs
    const { token, displayName, entity_id } = expectedData;
    assert.strictEqual(persistedTokenData.token, token, 'setTokenData has expected token');
    assert.strictEqual(persistedTokenData.displayName, displayName, 'setTokenData has expected display name');
    assert.strictEqual(persistedTokenData.entity_id, entity_id, 'setTokenData has expected entity_id');
  });

  test('it calls onAuthSuccess on submit for custom path', async function (assert) {
    assert.expect(1);
    // Setup
    this.path = `${this.authType}-custom`;
    this.loginData = { ...this.loginData, path: this.path };
    this.stubRequests();
    // Render and log in
    await this.renderComponent();
    await fillIn(AUTH_FORM.selectMethod, this.authType);
    // toggle mount path input to specify custom path
    await fillInLoginFields(this.loginData, { toggleOptions: true });
    if (this.authType === 'oidc') {
      triggerMessageEvent(this.path);
    }
    await click(GENERAL.submitButton);

    await waitUntil(() => this.onAuthSuccess.calledOnce);
    const [actual] = this.onAuthSuccess.lastCall.args;
    const expected = { namespace: '', token: this.tokenName, isRoot: false };
    assert.propEqual(actual, expected, `onAuthSuccess called with: ${JSON.stringify(actual)}`);
  });
};

module('Integration | Component | auth | page | method authentication', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    setupTestContext(this);
    this.auth = this.owner.lookup('service:auth');
    this.setTokenDataSpy = sinon.spy(this.auth, 'setTokenData');
  });

  hooks.afterEach(function () {
    window.localStorage.clear();
  });

  module('github', function (hooks) {
    hooks.beforeEach(async function () {
      this.authType = 'github';
      this.loginData = { token: 'mysupersecuretoken' };
      this.path = this.authType;
      this.response = RESPONSE_STUBS.github;
      this.tokenName = 'vault-github☃1';
      this.stubRequests = () => {
        this.server.post(`/auth/${this.path}/login`, () => this.response);
      };
    });

    methodAuthenticationTests(test);
  });

  module('jwt', function (hooks) {
    hooks.beforeEach(async function () {
      this.authType = 'jwt';
      this.loginData = { role: 'some-dev', jwt: 'jwttoken' };
      this.path = this.authType;
      this.response = RESPONSE_STUBS.jwt.login;
      this.tokenName = 'vault-jwt☃1';
      this.routerStub = sinon.stub(this.owner.lookup('service:router'), 'urlFor').returns('123-example.com');
      // Requests are stubbed in the order they are hit
      this.stubRequests = () => {
        // passing a dynamic path so that even the OIDC form renders the JWT token input
        // (there is test coverage elsewhere to assert switching between methods updates the form)
        this.server.post('/auth/:path/oidc/auth_url', () =>
          overrideResponse(400, { errors: [ERROR_JWT_LOGIN] })
        );
        this.server.post(`/auth/${this.path}/login`, () => this.response);
        this.server.get(`/auth/token/lookup-self`, () => RESPONSE_STUBS.jwt['lookup-self']);
      };
    });

    hooks.afterEach(function () {
      this.routerStub.restore();
    });

    methodAuthenticationTests(test);
  });

  module('ldap', function (hooks) {
    hooks.beforeEach(async function () {
      this.authType = 'ldap';
      this.loginData = { username: 'matilda', password: 'password' };
      this.path = this.authType;
      this.response = RESPONSE_STUBS.ldap;
      this.tokenName = 'vault-ldap☃1';

      this.stubRequests = () => {
        this.server.post(`/auth/${this.path}/login/${this.loginData.username}`, () => this.response);
      };
    });

    methodAuthenticationTests(test);
  });

  module('oidc', function (hooks) {
    hooks.beforeEach(async function () {
      this.authType = 'oidc';
      this.loginData = { role: 'some-dev' };
      this.path = this.authType;
      this.response = RESPONSE_STUBS.oidc['oidc/callback'];
      this.tokenName = 'vault-token☃1';
      // Requests are stubbed in the order they are hit
      this.stubRequests = () => {
        this.server.post(`/auth/${this.path}/oidc/auth_url`, () => {
          return { data: { auth_url: 'http://dev-foo-bar.com' } };
        });
        this.server.get(`/auth/${this.path}/oidc/callback`, () => this.response);
        this.server.get(`/auth/token/lookup-self`, () => RESPONSE_STUBS.oidc['lookup-self']);
      };

      // additional OIDC setup
      this.routerStub = sinon.stub(this.owner.lookup('service:router'), 'urlFor').returns('123-example.com');
      this.windowStub = windowStub();
    });

    hooks.afterEach(function () {
      this.routerStub.restore();
      this.windowStub.restore();
    });

    methodAuthenticationTests(test);
  });

  module('okta', function (hooks) {
    hooks.beforeEach(async function () {
      this.authType = 'okta';
      this.loginData = { username: 'matilda', password: 'password' };
      this.path = this.authType;
      this.response = RESPONSE_STUBS.okta;
      this.tokenName = 'vault-okta☃1';
      this.stubRequests = () => {
        this.server.post(`/auth/${this.path}/login/${this.loginData.username}`, () => this.response);
      };
    });

    methodAuthenticationTests(test);
  });

  module('radius', function (hooks) {
    hooks.beforeEach(async function () {
      this.authType = 'radius';
      this.loginData = { username: 'matilda', password: 'password' };
      this.path = this.authType;
      this.response = RESPONSE_STUBS.radius;
      this.tokenName = 'vault-radius☃1';
      this.stubRequests = () => {
        this.server.post(`/auth/${this.path}/login/${this.loginData.username}`, () => this.response);
      };
    });

    methodAuthenticationTests(test);
  });

  module('token', function (hooks) {
    hooks.beforeEach(async function () {
      this.authType = 'token';
      this.tokenName = 'vault-token☃1';
      this.server.get('/auth/token/lookup-self', () => RESPONSE_STUBS.token);
    });

    test('it sets token data and calls onAuthSuccess', async function (assert) {
      assert.expect(6);
      await this.renderComponent();
      await fillIn(AUTH_FORM.selectMethod, this.authType);
      await fillInLoginFields({ token: 'mysupersecuretoken' });
      await click(GENERAL.submitButton);

      const [actual] = this.onAuthSuccess.lastCall.args;
      const expected = { namespace: '', token: this.tokenName, isRoot: false };
      assert.propEqual(actual, expected, `onAuthSuccess called with: ${JSON.stringify(actual)}`);

      const [tokenName, persistedTokenData] = this.setTokenDataSpy.lastCall.args;
      const expectedTokenData = {
        ...TOKEN_DATA[this.authType],
        // there are other tests that confirm this calculation happens as expected, just copy value from spy
        tokenExpirationEpoch: persistedTokenData.tokenExpirationEpoch,
      };
      assert.strictEqual(tokenName, this.tokenName, 'setTokenData is called with expected token name');
      assert.propEqual(persistedTokenData, expectedTokenData, 'setTokenData is called with expected data');

      // propEqual failures are challenging to parse in CI so pulling out a couple of important attrs
      const { token, displayName, entity_id } = expectedTokenData;
      assert.strictEqual(persistedTokenData.token, token, 'setTokenData has expected token');
      assert.strictEqual(
        persistedTokenData.displayName,
        displayName,
        'setTokenData has expected display name'
      );
      assert.strictEqual(persistedTokenData.entity_id, entity_id, 'setTokenData has expected entity_id');
    });
  });

  module('userpass', function (hooks) {
    hooks.beforeEach(async function () {
      this.authType = 'userpass';
      this.loginData = { username: 'matilda', password: 'password' };
      this.path = this.authType;
      this.response = RESPONSE_STUBS.userpass;
      this.tokenName = 'vault-userpass☃1';
      this.stubRequests = () => {
        this.server.post(`/auth/${this.path}/login/${this.loginData.username}`, () => this.response);
      };
    });

    methodAuthenticationTests(test);
  });

  // ENTERPRISE METHODS
  module('saml', function (hooks) {
    hooks.beforeEach(async function () {
      this.version.type = 'enterprise';
      this.authType = 'saml';
      this.path = this.authType;
      this.loginData = { role: 'some-dev' };
      this.response = RESPONSE_STUBS.saml['saml/token'];
      this.tokenName = 'vault-token☃1';
      // Requests are stubbed in the order they are hit
      this.stubRequests = () => {
        this.server.put(`/auth/${this.path}/sso_service_url`, () => ({
          data: {
            sso_service_url: 'test/fake/sso/route',
            token_poll_id: '1234',
          },
        }));
        this.server.put(`/auth/${this.path}/token`, () => this.response);
        this.server.get(`/auth/token/lookup-self`, () => RESPONSE_STUBS.saml['lookup-self']);
      };
      this.windowStub = windowStub();
    });

    hooks.afterEach(function () {
      this.windowStub.restore();
    });

    methodAuthenticationTests(test);
  });
});
