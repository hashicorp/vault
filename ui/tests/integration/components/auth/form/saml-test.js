/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import hbs from 'htmlbars-inline-precompile';
import { click, fillIn, find, render } from '@ember/test-helpers';
import sinon from 'sinon';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { overrideResponse } from 'vault/tests/helpers/stubs';
import { windowStub } from 'vault/tests/helpers/oidc-window-stub';

module('Integration | Component | auth | form | saml', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.authType = 'saml';
    this.expectedFields = ['role'];

    this.authenticateStub = sinon.stub(this.owner.lookup('service:auth'), 'authenticate');
    this.store = this.owner.lookup('service:store');
    this.cluster = { id: 1 };
    this.onError = sinon.spy();
    this.onSuccess = sinon.spy();
    this.windowStub = windowStub();
    sinon.replaceGetter(window, 'screen', () => ({ height: 600, width: 500 }));

    // role request
    this.server.put('/auth/saml/sso_service_url', () => {
      return {
        data: {
          sso_service_url: 'https://my-single-sign-on-url.com',
          token_poll_id: '4fe2ec01-1f56-b665-0ba2-09c7bca10ae8',
        },
      };
    });
    // polling request
    this.server.put('/auth/saml/token', () => {
      return { auth: { client_token: 'my_token' } };
    });

    this.renderComponent = ({ yieldBlock = false } = {}) => {
      if (yieldBlock) {
        return render(hbs`
          <Auth::Form::Saml 
            @authType={{this.authType}} 
            @cluster={{this.cluster}}
            @onError={{this.onError}}
            @onSuccess={{this.onSuccess}}
          >
            <:advancedSettings>
              <label for="path">Mount path</label>
              <input data-test-input="path" id="path" name="path" type="text" /> 
            </:advancedSettings>
          </Auth::Form::Saml>`);
      }
      return render(hbs`
      <Auth::Form::Saml       
        @authType={{this.authType}}
        @cluster={{this.cluster}}
        @onError={{this.onError}}
        @onSuccess={{this.onSuccess}}
      />`);
    };
  });

  hooks.afterEach(function () {
    this.windowStub.restore();
    this.authenticateStub.restore();
  });

  test('it renders helper text', async function (assert) {
    await this.renderComponent();
    const id = find(GENERAL.inputByAttr('role')).id;
    assert
      .dom(`#helper-text-${id}`)
      .hasText('Vault will use the default role to sign in if this field is left blank.');
  });

  test('it renders warning if insecure context is detected', async function (assert) {
    sinon.replaceGetter(window, 'isSecureContext', () => false);

    await this.renderComponent();
    assert
      .dom('[data-test-saml-auth-not-allowed]')
      .hasText(
        'Insecure context detected Logging in with a SAML auth method requires a browser in a secure context. Read more about secure contexts.'
      );
  });

  test('it requests sso_service_url and opens popup on submit if role is empty', async function (assert) {
    assert.expect(6);
    this.server.put('/auth/saml/sso_service_url', (_, req) => {
      const { acs_url, role } = JSON.parse(req.requestBody);
      assert.strictEqual(acs_url, `${window.origin}/v1/auth/saml/callback`, 'it builds acs_url for payload');
      assert.strictEqual(role, '', 'role has no value');
      return {
        data: {
          sso_service_url: 'https://my-single-sign-on-url.com',
          token_poll_id: '4fe2ec01-1f56-b665-0ba2-09c7bca10ae8',
        },
      };
    });

    await this.renderComponent();
    await click(GENERAL.submitButton);

    const [sso_service_url, name, dimensions] = this.windowStub.lastCall.args;
    assert.strictEqual(
      sso_service_url,
      'https://my-single-sign-on-url.com',
      'it calls window opener with sso_service_url returned by role request'
    );
    assert.strictEqual(sso_service_url, 'https://my-single-sign-on-url.com');
    assert.strictEqual(name, 'vaultSAMLWindow', 'it calls window opener with expected name');
    assert.strictEqual(
      dimensions,
      'width=500,height=600,resizable,scrollbars=yes,top=0,left=0',
      'it calls window opener with expected dimensions'
    );
  });

  test('it requests sso_service_url with inputted role and default path', async function (assert) {
    assert.expect(6);
    this.server.put('/auth/saml/sso_service_url', (_, req) => {
      const { acs_url, role } = JSON.parse(req.requestBody);
      assert.strictEqual(acs_url, `${window.origin}/v1/auth/saml/callback`, 'it builds acs_url for payload');
      assert.strictEqual(role, 'some-dev', 'payload contains role');
      return {
        data: {
          sso_service_url: 'https://my-single-sign-on-url.com',
          token_poll_id: '4fe2ec01-1f56-b665-0ba2-09c7bca10ae8',
        },
      };
    });

    await this.renderComponent();
    await fillIn(GENERAL.inputByAttr('role'), 'some-dev');
    await click(GENERAL.submitButton);

    const [sso_service_url, name, dimensions] = this.windowStub.lastCall.args;
    assert.strictEqual(
      sso_service_url,
      'https://my-single-sign-on-url.com',
      'it calls window opener with sso_service_url returned by role request'
    );
    assert.strictEqual(sso_service_url, 'https://my-single-sign-on-url.com');
    assert.strictEqual(name, 'vaultSAMLWindow', 'it calls window opener with expected name');
    assert.strictEqual(
      dimensions,
      'width=500,height=600,resizable,scrollbars=yes,top=0,left=0',
      'it calls window opener with expected dimensions'
    );
  });

  test('it requests sso_service_url with custom path', async function (assert) {
    assert.expect(6);
    const path = 'custom-path';
    this.server.put(`/auth/${path}/sso_service_url`, (_, req) => {
      const { acs_url, role } = JSON.parse(req.requestBody);
      assert.strictEqual(
        acs_url,
        `${window.origin}/v1/auth/${path}/callback`,
        'it builds acs_url for payload'
      );
      assert.strictEqual(role, 'some-dev', 'payload contains role');
      return {
        data: {
          sso_service_url: 'https://my-single-sign-on-url.com',
          token_poll_id: '4fe2ec01-1f56-b665-0ba2-09c7bca10ae8',
        },
      };
    });

    await this.renderComponent({ yieldBlock: true });
    await fillIn(GENERAL.inputByAttr('role'), 'some-dev');
    await fillIn(GENERAL.inputByAttr('path'), path);
    await click(GENERAL.submitButton);

    const [sso_service_url, name, dimensions] = this.windowStub.lastCall.args;
    assert.strictEqual(
      sso_service_url,
      'https://my-single-sign-on-url.com',
      'it calls window opener with sso_service_url returned by role request'
    );
    assert.strictEqual(sso_service_url, 'https://my-single-sign-on-url.com');
    assert.strictEqual(name, 'vaultSAMLWindow', 'it calls window opener with expected name');
    assert.strictEqual(
      dimensions,
      'width=500,height=600,resizable,scrollbars=yes,top=0,left=0',
      'it calls window opener with expected dimensions'
    );
  });

  test('it polls token request', async function (assert) {
    assert.expect(2); // auth/saml/token url should be requested twice

    let count = 0;
    this.server.put('/auth/saml/token', () => {
      count++;
      const msg =
        count === 1
          ? 'it makes initial request to token url'
          : 'it re-requests token url if httpStatus was 401';
      assert.true(true, msg);

      if (count === 1) {
        return overrideResponse(401);
      } else {
        return { auth: { client_token: 'my_token' } };
      }
    });
    await this.renderComponent();
    await click(GENERAL.submitButton);
  });

  test('it calls auth service with token request callback data', async function (assert) {
    await this.renderComponent();
    await click(GENERAL.submitButton);

    const [actual] = this.authenticateStub.lastCall.args;
    assert.propEqual(
      actual.data,
      {
        token: 'my_token',
      },
      'auth service "authenticate" method is called token callback data'
    );
  });

  test('it calls onSuccess if auth service authentication is successful', async function (assert) {
    const expectedResponse = {
      namespace: '',
      token: 'my_token',
      isRoot: false,
    };
    // stub happy response
    this.authenticateStub.returns(expectedResponse);

    await this.renderComponent();
    await click(GENERAL.submitButton);
    const [actualResponse, methodData] = this.onSuccess.lastCall.args;
    assert.propEqual(actualResponse, expectedResponse, 'onSuccess is called with auth response');
    assert.strictEqual(methodData.path, undefined, 'onSuccess is called without path value');
    assert.strictEqual(methodData.selectedAuth, 'saml', 'onSuccess is called with selected auth type');
  });

  test('it calls onError if auth service authentication fails', async function (assert) {
    this.authenticateStub.throws('permission denied!!');
    await this.renderComponent();
    await click(GENERAL.submitButton);
    const [actual] = this.onError.lastCall.args;
    assert.strictEqual(
      actual,
      'Authentication failed: Sinon-provided permission denied!!',
      'onError called with auth service failure'
    );
  });

  test('it calls onError if sso_service_url request fails', async function (assert) {
    // role request
    this.server.put('/auth/saml/sso_service_url', () => overrideResponse(403));
    await this.renderComponent();
    await click(GENERAL.submitButton);
    const [actual] = this.onError.lastCall.args;
    assert.strictEqual(
      actual,
      'Authentication failed: permission denied',
      'onError called with sso_service_url failure'
    );
  });

  test('it calls onError if polling token errors in status code that is NOT 401', async function (assert) {
    this.server.put('/auth/saml/token', () => overrideResponse(500));
    await this.renderComponent();
    await click(GENERAL.submitButton);
    const [actual] = this.onError.lastCall.args;
    assert.strictEqual(
      actual,
      'Authentication failed: Ember Data Request PUT /v1/auth/saml/token returned a 500\nPayload (application/json)\n{}',
      'onError called with auth failure'
    );
  });
});
