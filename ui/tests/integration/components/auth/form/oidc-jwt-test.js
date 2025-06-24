/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import hbs from 'htmlbars-inline-precompile';
import { click, fillIn, find, render, settled, waitUntil } from '@ember/test-helpers';
import { _cancelTimers as cancelTimers } from '@ember/runloop';
import { AUTH_FORM } from 'vault/tests/helpers/auth/auth-form-selectors';
import { callbackData } from 'vault/tests/helpers/oidc-window-stub';
import { ERROR_JWT_LOGIN } from 'vault/components/auth/form/oidc-jwt';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { overrideResponse } from 'vault/tests/helpers/stubs';
import { setupMirage } from 'ember-cli-mirage/test-support';
import * as parseURL from 'core/utils/parse-url';
import sinon from 'sinon';
import testHelper from './auth-form-test-helper';

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

  test('it submits form data with defaults', async function (assert) {
    await this.renderComponent();

    await fillIn(GENERAL.inputByAttr('role'), 'some-dev');
    await fillIn(GENERAL.inputByAttr('jwt'), 'some-jwt-token');

    await click(GENERAL.submitButton);
    const [actual] = this.authenticateStub.lastCall.args;
    assert.propEqual(
      actual.data,
      this.expectedSubmit.default,
      'auth service "authenticate" method is called with form data'
    );
  });

  test('it submits form data from yielded inputs', async function (assert) {
    await this.renderComponent({ yieldBlock: true });
    await fillIn(GENERAL.inputByAttr('role'), 'some-dev');
    await fillIn(GENERAL.inputByAttr('jwt'), 'some-jwt-token');
    await fillIn(GENERAL.inputByAttr('path'), `custom-${this.authType}`);

    await click(GENERAL.submitButton);
    const [actual] = this.authenticateStub.lastCall.args;
    assert.propEqual(
      actual.data,
      this.expectedSubmit.custom,
      'auth service "authenticate" method is called with yielded form data'
    );
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
  test('it renders fields', async function (assert) {
    await this.renderComponent();
    assert.dom(AUTH_FORM.authForm(this.authType)).exists(`${this.authType}: it renders form component`);
    assert.dom(GENERAL.submitButton).hasText('Sign in with OIDC Provider');
    this.expectedFields.forEach((field) => {
      assert.dom(GENERAL.inputByAttr(field)).exists(`${this.authType}: it renders ${field}`);
    });
  });

  test('it renders provider icon and name', async function (assert) {
    const parseURLStub = sinon.stub(parseURL, 'default').returns({ hostname: 'auth0.com' });
    this.server.post(`/auth/${this.authType}/oidc/auth_url`, () => {
      return { data: { auth_url: '123.auth0.com' } };
    });
    await this.renderComponent();
    assert.dom(GENERAL.submitButton).hasText('Sign in with Auth0');
    assert.dom(GENERAL.icon('auth0')).exists();
    parseURLStub.restore();
  });

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

  // auth_url error handling on submit
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

  test('it fires onError callback on submit when auth_url request is successful but missing auth_url key', async function (assert) {
    this.server.post('/auth/:path/oidc/auth_url', () => ({ data: {} }));
    await this.renderComponent();
    await click(GENERAL.submitButton);

    const [actual] = this.onError.lastCall.args;
    assert.strictEqual(
      actual,
      'Authentication failed: Missing auth_url. Please check that allowed_redirect_uris for the role include this mount path.',
      'it calls onError'
    );
  });
  // end auth_url error handling

  // prepareForOIDC logic tests
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
};

module('Integration | Component | auth | form | oidc-jwt', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.authenticateStub = sinon.stub(this.owner.lookup('service:auth'), 'authenticate');
    this.cluster = { id: 1 };
    this.onError = sinon.spy();
    this.onSuccess = sinon.spy();

    // additional test setup for oidc/jwt business
    this.store = this.owner.lookup('service:store');
    this.routerStub = sinon.stub(this.owner.lookup('service:router'), 'urlFor').returns('123-example.com');

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
    this.authenticateStub.restore();
    this.routerStub.restore();
  });

  test('it renders helper text', async function (assert) {
    await this.renderComponent();
    const id = find(GENERAL.inputByAttr('role')).id;
    assert
      .dom(`#helper-text-${id}`)
      .hasText('Vault will use the default role to sign in if this field is left blank.');
  });

  module('oidc', function (hooks) {
    hooks.beforeEach(function () {
      this.authType = 'oidc';
      this.expectedFields = ['role'];
    });

    // base component functionality so outside login workflow modules
    authUrlRequestTests(test);

    module('login workflow: jwt token', function (hooks) {
      hooks.beforeEach(function () {
        // stubbing this request shows JWT token input and does not perform OIDC
        this.server.post(`/auth/:path/oidc/auth_url`, () => {
          return overrideResponse(400, { errors: [ERROR_JWT_LOGIN] });
        });
        this.expectedFields = ['role', 'jwt'];
        this.expectedSubmit = {
          default: { path: 'oidc', role: 'some-dev', jwt: 'some-jwt-token' },
          custom: { path: 'custom-oidc', role: 'some-dev', jwt: 'some-jwt-token' },
        };
      });

      testHelper(test, { standardSubmit: false });

      jwtLoginTests(test);
    });

    module('login workflow: oidc', function (hooks) {
      hooks.beforeEach(function () {
        // for oidc login workflow only
        this.windowStub = sinon.stub(window, 'open');
      });

      hooks.afterEach(function () {
        this.windowStub.restore();
      });

      oidcLoginTests(test);
    });
  });

  module('jwt', function (hooks) {
    hooks.beforeEach(function () {
      this.authType = 'jwt';
      this.expectedFields = ['role'];
    });

    // base component functionality so outside login workflow modules
    authUrlRequestTests(test);

    module('login workflow: jwt token', function (hooks) {
      hooks.beforeEach(function () {
        // stubbing this request shows JWT token input and does not perform OIDC
        this.server.post(`/auth/:path/oidc/auth_url`, () => {
          return overrideResponse(400, { errors: [ERROR_JWT_LOGIN] });
        });
        this.expectedFields = ['role', 'jwt'];
        this.expectedSubmit = {
          default: { path: 'jwt', role: 'some-dev', jwt: 'some-jwt-token' },
          custom: { path: 'custom-jwt', role: 'some-dev', jwt: 'some-jwt-token' },
        };
      });

      testHelper(test, { standardSubmit: false });

      jwtLoginTests(test);
    });

    module('login workflow: oidc', function (hooks) {
      hooks.beforeEach(function () {
        // for oidc login workflow only
        this.windowStub = sinon.stub(window, 'open');
      });

      hooks.afterEach(function () {
        this.windowStub.restore();
      });

      oidcLoginTests(test);
    });
  });
});
