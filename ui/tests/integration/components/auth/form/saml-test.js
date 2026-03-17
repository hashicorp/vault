/**
 * Copyright IBM Corp. 2016, 2025
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
import * as uuid from 'core/utils/uuid';
import { getErrorResponse } from 'vault/tests/helpers/api/error-response';
import { RESPONSE_STUBS } from 'vault/tests/helpers/auth/response-stubs';
import { AUTH_METHOD_LOGIN_DATA } from 'vault/tests/helpers/auth/auth-helpers';
import authFormTestHelper from './auth-form-test-helper';

module('Integration | Component | auth | form | saml', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.authType = 'saml';
    this.expectedFields = ['role'];

    this.tokenPollId = '4fe2ec01-1f56-b665-0ba2-09c7bca10ae8';
    this.cluster = { id: 1 };
    this.onError = sinon.spy();
    this.handleAuthResponse = sinon.spy();
    // Window stub
    this.windowStub = windowStub();
    sinon.replaceGetter(window, 'screen', () => ({ height: 600, width: 500 }));
    // Auth request stub
    const api = this.owner.lookup('service:api');
    this.authenticateStub = sinon.stub(api.auth, 'samlWriteToken');
    this.authResponse = RESPONSE_STUBS.saml['saml/token'];
    this.loginData = AUTH_METHOD_LOGIN_DATA.saml;
    // stub uuid so verifier can be asserted
    this.verifier = uuid.default();
    sinon.stub(uuid, 'default').returns(this.verifier);
    // role request
    this.server.post('/auth/:path/sso_service_url', () => {
      return {
        data: {
          sso_service_url: 'https://my-single-sign-on-url.com',
          token_poll_id: this.tokenPollId,
        },
      };
    });

    this.assertSubmit = (assert, loginRequestArgs, loginData) => {
      const [path, { client_verifier, token_poll_id }] = loginRequestArgs;
      // if path is included in loginData, a custom path was submitted
      const expectedPath = loginData?.path || this.authType;
      assert.strictEqual(path, expectedPath, 'it calls samlWriteToken with expected path');
      assert.strictEqual(client_verifier, this.verifier, 'it calls samlWriteToken with verifier');
      assert.strictEqual(token_poll_id, this.tokenPollId, 'it calls samlWriteToken with tokenPollId');
    };

    this.renderComponent = ({ yieldBlock = false } = {}) => {
      if (yieldBlock) {
        return render(hbs`
          <Auth::Form::Saml 
            @authType={{this.authType}} 
            @cluster={{this.cluster}}
            @onError={{this.onError}}
            @handleAuthResponse={{this.handleAuthResponse}}
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
        @handleAuthResponse={{this.handleAuthResponse}}
      />`);
    };
  });

  hooks.afterEach(function () {
    this.windowStub.restore();
    this.authenticateStub.restore();
  });

  authFormTestHelper(test);

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
    this.server.post('/auth/saml/sso_service_url', (_, req) => {
      const { acs_url, role } = JSON.parse(req.requestBody);
      assert.strictEqual(acs_url, `${window.origin}/v1/auth/saml/callback`, 'it builds acs_url for payload');
      assert.strictEqual(role, '', 'role has no value');
      return {
        data: {
          sso_service_url: 'https://my-single-sign-on-url.com',
          token_poll_id: this.tokenPollId,
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
    this.server.post('/auth/saml/sso_service_url', (_, req) => {
      const { acs_url, role } = JSON.parse(req.requestBody);
      assert.strictEqual(acs_url, `${window.origin}/v1/auth/saml/callback`, 'it builds acs_url for payload');
      assert.strictEqual(role, 'some-dev', 'payload contains role');
      return {
        data: {
          sso_service_url: 'https://my-single-sign-on-url.com',
          token_poll_id: this.tokenPollId,
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

  test('it re-requests samlWriteToken if 401 is returned', async function (assert) {
    assert.expect(1);
    this.authenticateStub.onFirstCall().rejects(getErrorResponse({}, 401));
    this.authenticateStub.onSecondCall().rejects(getErrorResponse({}, 401));
    // MAX_TRIES in the component is set to 3 for tests
    this.authenticateStub.onThirdCall().resolves(this.authResponse);
    await this.renderComponent();
    await click(GENERAL.submitButton);
    assert.strictEqual(this.authenticateStub.callCount, 3, 'it polls token request until request resolves');
  });

  test('it calls onError if sso_service_url request fails', async function (assert) {
    // role request
    this.server.post('/auth/saml/sso_service_url', () => overrideResponse(403));
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
    this.authenticateStub.rejects(getErrorResponse({ errors: ['uh oh!'] }, 500));
    await this.renderComponent();
    await click(GENERAL.submitButton);
    const [actual] = this.onError.lastCall.args;
    assert.strictEqual(actual, 'Authentication failed: uh oh!', 'onError called with expected failure');
  });
});
