/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import hbs from 'htmlbars-inline-precompile';
import { click, fillIn, find, render, settled, waitUntil } from '@ember/test-helpers';
import { _cancelTimers as cancelTimers } from '@ember/runloop';
import { callbackData } from 'vault/tests/helpers/oidc-window-stub';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { overrideResponse } from 'vault/tests/helpers/stubs';
import { setupMirage } from 'ember-cli-mirage/test-support';
import * as parseURL from 'core/utils/parse-url';
import sinon from 'sinon';
import authFormTestHelper from './auth-form-test-helper';
import { RESPONSE_STUBS } from 'vault/tests/helpers/auth/response-stubs';
import { DOMAIN_PROVIDER_MAP, ERROR_JWT_LOGIN } from 'vault/utils/auth-form-helpers';
import { dasherize } from '@ember/string';

/* 
The OIDC and JWT mounts call the same endpoint (see docs https://developer.hashicorp.com/vault/docs/auth/jwt )
because of this the same component is used to render both method types. 
The module name refers to the selected auth type and there is test coverage in each to cover
1. auth url request (fetching role) is made when expected
2. JWT token login flow
2. OIDC exchange/login situation
*/

const authUrlRequestTests = (test) => {
  test('it requests auth_url when it initially renders', async function (assert) {
    assert.expect(2);
    this.server.post(`/auth/${this.authType}/oidc/auth_url`, (_, req) => {
      const { role } = JSON.parse(req.requestBody);
      assert.true(true, 'it makes request to auth_url');
      assert.strictEqual(role, '', 'role is empty');
      return { data: { auth_url: '123-example.com' } };
    });
    await this.renderComponent();
  });

  test('it re-requests auth_url when input changes: role', async function (assert) {
    // request assertions should be hit twice, once on initial render and again on role change
    assert.expect(4);
    let count = 0;
    this.server.post(`/auth/${this.authType}/oidc/auth_url`, (_, req) => {
      count++;
      const { role } = JSON.parse(req.requestBody);
      assert.true(true, 'it makes request to auth_url');
      const expectedRole = count === 1 ? '' : 'myrole';
      assert.strictEqual(role, expectedRole, 'payload has expected role');
      return { data: { auth_url: '123-example.com' } };
    });
    await this.renderComponent({ yieldBlock: true });
    await fillIn(GENERAL.inputByAttr('role'), 'myrole');
  });

  test('it re-requests auth_url when input changes: path', async function (assert) {
    assert.expect(2);
    let firstRequest, secondRequest;
    this.server.post(`/auth/${this.authType}/oidc/auth_url`, () => {
      firstRequest = true;
      return { data: { auth_url: '123-example.com' } };
    });
    this.server.post(`/auth/mypath/oidc/auth_url`, () => {
      secondRequest = true;
      return { data: { auth_url: '123-example.com' } };
    });
    await this.renderComponent({ yieldBlock: true });
    await fillIn(GENERAL.inputByAttr('path'), 'mypath');

    // asserting this way instead of inside the request to ensure each endpoint is hit.
    // (asserting within each request would rely on assertion count and could result in a false positive)
    assert.true(firstRequest, 'it makes FIRST request to auth_url with default path');
    assert.true(secondRequest, 'it makes SECOND request to auth_url with custom path');
  });
};

const jwtLoginTests = (test) => {
  test('it renders sign in button text', async function (assert) {
    await this.renderComponent();
    assert.dom(GENERAL.submitButton).hasText('Sign in');
  });

  test('it does NOT re-request the auth_url when jwt token changes', async function (assert) {
    assert.expect(1); // the assertion in the stubbed request should not be hit
    await this.renderComponent();
    this.server.post(`/auth/${this.authType}/oidc/auth_url`, () => {
      // we can't throw an error here because the component catches error handling.
      // setting this assertion to fail intentionally because this request should not have been made
      assert.false(true, 'request made to auth_url and it should not have been requested');
    });
    await fillIn(GENERAL.inputByAttr('jwt'), 'mytoken');
    assert.dom(GENERAL.inputByAttr('jwt')).hasValue('mytoken');
  });
};

const oidcLoginTests = (test) => {
  // true success has to be asserted in acceptance tests because it's not possible to mock a trusted message event
  test('it opens the popup window on submit', async function (assert) {
    this.server.post(`/auth/${this.authType}/oidc/auth_url`, () => {
      return { data: { auth_url: '123-example.com' } };
    });
    sinon.replaceGetter(window, 'screen', () => ({ height: 600, width: 500 }));
    await this.renderComponent();
    await fillIn(GENERAL.inputByAttr('role'), 'test');
    await click(GENERAL.submitButton);
    await waitUntil(() => {
      return this.windowStub.calledOnce;
    });

    const [authURL, windowName, windowDimensions] = this.windowStub.lastCall.args;

    assert.strictEqual(authURL, '123-example.com', 'window stub called with auth_url');
    assert.strictEqual(windowName, 'vaultOIDCWindow', 'window stub called with name');
    assert.strictEqual(
      windowDimensions,
      'width=500,height=600,resizable,scrollbars=yes,top=0,left=0',
      'window stub called with dimensions'
    );
    sinon.restore();
  });

  /* Tests for auth_url error handling on submit */
  test('it fires onError callback on submit when auth_url request fails with 400', async function (assert) {
    this.server.post('/auth/:path/oidc/auth_url', () => overrideResponse(400));
    await this.renderComponent();
    await click(GENERAL.submitButton);

    const [actual] = this.onError.lastCall.args;
    assert.strictEqual(actual, 'Authentication failed: Invalid role. Please try again.');
  });

  test('it fires onError callback on submit when auth_url request fails with 403', async function (assert) {
    this.server.post('/auth/:path/oidc/auth_url', () => overrideResponse(403));
    await this.renderComponent();
    await click(GENERAL.submitButton);

    const [actual] = this.onError.lastCall.args;
    assert.strictEqual(actual, 'Authentication failed: Error fetching role: permission denied');
  });

  test('it fires onError callback on submit when auth_url request is successful but missing auth_url', async function (assert) {
    this.server.post('/auth/:path/oidc/auth_url', () => ({ data: { auth_url: '' } }));
    await this.renderComponent();
    await click(GENERAL.submitButton);

    const [actual] = this.onError.lastCall.args;
    assert.strictEqual(
      actual,
      'Authentication failed: Missing auth_url. Please check that allowed_redirect_uris for the role include this mount path.',
      'it calls onError'
    );
  });
  /* END auth_url error handling */

  /* test for prepareForOIDC logic */
  test('fails silently when event is not trusted', async function (assert) {
    assert.expect(2);
    this.server.post(`/auth/${this.authType}/oidc/auth_url`, () => {
      return { data: { auth_url: '123-example.com' } };
    });
    // prevent test incorrectly passing because the event isn't triggered at all
    // by also asserting that the message event fires
    const messageData = callbackData();
    const assertEvent = (event) => {
      assert.propEqual(event.data, messageData, 'message event fires');
    };
    window.addEventListener('message', assertEvent);

    await this.renderComponent();
    await fillIn(GENERAL.inputByAttr('role'), 'test');
    await click(GENERAL.submitButton);
    await waitUntil(() => {
      return this.windowStub.calledOnce;
    });
    // mocking a message event is always untrusted (there is no way to override isTrusted on the window object)
    window.dispatchEvent(new MessageEvent('message', { data: messageData }));

    cancelTimers();
    await settled();
    assert.false(this.onSuccess.called, 'onSuccess is not called');

    // Cleanup
    window.removeEventListener('message', assertEvent);
  });

  // not the greatest test because this assertion would also pass if the event.origin === window.origin.
  // because event.isTrusted is always false (another condition checked by the component)
  // but this is good enough because the origin logic is checked first in the conditional.
  test('it fails silently when event origin does not match window origin', async function (assert) {
    assert.expect(3);
    this.server.post(`/auth/${this.authType}/oidc/auth_url`, () => {
      return { data: { auth_url: '123-example.com' } };
    });
    // prevent test incorrectly passing because the event isn't triggered at all
    // by also asserting that the message event fires
    const message = { data: callbackData(), origin: 'http://hackerz.com' };
    const assertEvent = (event) => {
      assert.propEqual(event.data, message.data, 'message has expected data');
      assert.strictEqual(event.origin, message.origin, 'message has expected origin');
    };
    window.addEventListener('message', assertEvent);

    await this.renderComponent();
    await fillIn(GENERAL.inputByAttr('role'), 'test');
    await click(GENERAL.submitButton);
    await waitUntil(() => {
      return this.windowStub.calledOnce;
    });

    window.dispatchEvent(new MessageEvent('message', message));
    cancelTimers();
    await settled();
    assert.false(this.onSuccess.called, 'onSuccess is not called');

    // Cleanup
    window.removeEventListener('message', assertEvent);
  });
  /* end of tests for prepareForOIDC logic */
};

module('Integration | Component | auth | form | oidc-jwt', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.cluster = { id: 1 };
    this.onError = sinon.spy();
    this.onSuccess = sinon.spy();
    this.routerStub = sinon.stub(this.owner.lookup('service:router'), 'urlFor').returns('123-example.com');
    const api = this.owner.lookup('service:api');

    this.oidcLoginSetup = () => {
      // authentication request called on submit for oidc login
      this.authenticateStub = sinon.stub(api.auth, 'jwtOidcCallback');
      this.authResponse = RESPONSE_STUBS.oidc['oidc/callback'];
      this.windowStub = sinon.stub(window, 'open');
    };

    this.jwtLoginSetup = () => {
      // stubbing this request shows JWT token input and does not perform OIDC
      this.server.post(`/auth/:path/oidc/auth_url`, () => {
        return overrideResponse(400, { errors: [ERROR_JWT_LOGIN] });
      });
      // authentication request called on submit for jwt tokens
      this.authenticateStub = sinon.stub(api.auth, 'jwtLogin');
      this.authResponse = RESPONSE_STUBS.jwt.login;
    };

    this.assertSubmit = (assert, loginRequestArgs, loginData) => {
      const [path, payload] = loginRequestArgs;
      // if path is included in loginData, a custom path was submitted
      const expectedPath = loginData?.path || this.authType;
      assert.strictEqual(path, expectedPath, `auth request made with path: ${expectedPath}`);

      // iterate through each item in the payload and check its value
      for (const field in payload) {
        const actualValue = payload[field];
        const expectedValue = loginData[field];
        assert.strictEqual(actualValue, expectedValue, `payload includes field: ${field}`);
      }
    };

    this.renderComponent = ({ yieldBlock = false } = {}) => {
      if (yieldBlock) {
        return render(hbs`
          <Auth::Form::OidcJwt 
            @authType={{this.authType}} 
            @cluster={{this.cluster}}
            @onError={{this.onError}}
            @onSuccess={{this.onSuccess}}
          >
            <:advancedSettings>
              <label for="path">Mount path</label>
              <input data-test-input="path" id="path" name="path" type="text" /> 
            </:advancedSettings>
          </Auth::Form::OidcJwt>`);
      }
      return render(hbs`
        <Auth::Form::OidcJwt       
        @authType={{this.authType}}
        @cluster={{this.cluster}}
        @onError={{this.onError}}
        @onSuccess={{this.onSuccess}}
        />
        `);
    };
  });

  hooks.afterEach(function () {
    this.routerStub.restore();
  });

  /* TESTS FOR BASE COMPONENT FUNCTIONALITY
   These tests intentionally do not set authType as they are asserting type agnostic functionality.
   This means the /auth_url request is missing the :path param and the console will output this request error:
   403 http://localhost:7357/v1/auth//oidc/auth_url
   */
  test('it renders helper text', async function (assert) {
    await this.renderComponent();
    const id = find(GENERAL.inputByAttr('role')).id;
    assert
      .dom(`#helper-text-${id}`)
      .hasText('Vault will use the default role to sign in if this field is left blank.');
  });

  for (const domain in DOMAIN_PROVIDER_MAP) {
    const provider = DOMAIN_PROVIDER_MAP[domain];

    test(`${provider}: it renders provider icon and name`, async function (assert) {
      // parseUrl uses the actual window origin, so stub the util's return instead of authUrl
      const parseURLStub = sinon.stub(parseURL, 'default').returns({ hostname: domain });
      await this.renderComponent();
      assert.dom(GENERAL.submitButton).hasText(`Sign in with ${provider}`);

      // Right now there is a bug in HDS where the ping-identity icon name has a trailing whitespace.
      // This test should fail when upgrading to an HDS version with the corrected icon name and then we can remove this conditional.
      const iconName = domain === 'ping.com' ? 'ping-identity ' : dasherize(provider.toLowerCase());
      // convenience message for HDS upgrade failure, can be removed when we upgrade
      const message =
        iconName === 'ping-identity '
          ? `If you are attempting to upgrade @hashicorp/design-system-components and this test is failing, please remove the icon override for Ping Identity in oidc-jwt.ts`
          : `it renders icon for ${domain}`;
      assert.dom(GENERAL.icon(iconName)).exists(message);
      parseURLStub.restore();
    });
  }

  test('it does not return provider unless domain matches completely', async function (assert) {
    assert.expect(2);
    // parseUrl uses the actual window origin, so stub the return
    const parseURLStub = sinon
      .stub(parseURL, 'default')
      .returns({ hostname: `http://custom-auth0-provider.com` });
    await this.renderComponent();
    assert.dom(GENERAL.submitButton).hasText('Sign in with OIDC Provider');
    assert.dom(GENERAL.icon()).doesNotExist();
    parseURLStub.restore();
  });

  module('@authType: oidc', function (hooks) {
    hooks.beforeEach(function () {
      this.authType = 'oidc';
    });

    // base component functionality so outside login workflow modules
    authUrlRequestTests(test);

    module('login workflow: jwt token', function (hooks) {
      hooks.beforeEach(function () {
        this.jwtLoginSetup();
        this.loginData = { role: 'some-dev', jwt: 'some-jwt-token' };
      });

      hooks.afterEach(function () {
        this.authenticateStub.restore();
      });

      authFormTestHelper(test);

      jwtLoginTests(test);
    });

    module('login workflow: oidc', function (hooks) {
      hooks.beforeEach(function () {
        this.oidcLoginSetup();
        this.loginData = { role: 'some-dev' };
      });

      hooks.afterEach(function () {
        this.authenticateStub.restore();
        this.windowStub.restore();
      });

      oidcLoginTests(test);
    });
  });

  module('@authType: jwt', function (hooks) {
    hooks.beforeEach(function () {
      this.authType = 'jwt';
    });

    // base component functionality so outside login workflow modules
    authUrlRequestTests(test);

    module('login workflow: jwt token', function (hooks) {
      hooks.beforeEach(function () {
        this.jwtLoginSetup();
        this.loginData = { role: 'some-dev', jwt: 'some-jwt-token' };
      });

      hooks.afterEach(function () {
        this.authenticateStub.restore();
      });

      authFormTestHelper(test);

      jwtLoginTests(test);
    });

    module('login workflow: oidc', function (hooks) {
      hooks.beforeEach(function () {
        this.oidcLoginSetup();
        this.loginData = { role: 'some-dev' };
      });

      hooks.afterEach(function () {
        this.authenticateStub.restore();
        this.windowStub.restore();
      });

      oidcLoginTests(test);
    });
  });
});
