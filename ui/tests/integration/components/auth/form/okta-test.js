/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import hbs from 'htmlbars-inline-precompile';
import { click, render } from '@ember/test-helpers';
import sinon from 'sinon';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { AUTH_METHOD_LOGIN_DATA, fillInLoginFields } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { AUTH_FORM } from 'vault/tests/helpers/auth/auth-form-selectors';
import * as uuid from 'core/utils/uuid';
import { Response } from 'miragejs';
import authFormTestHelper from './auth-form-test-helper';
import { RESPONSE_STUBS } from 'vault/tests/helpers/auth/response-stubs';

module('Integration | Component | auth | form | okta', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.authType = 'okta';
    this.cluster = { id: 1 };
    this.onError = sinon.spy();
    this.onSuccess = sinon.spy();

    // stub uuid so auth/okta/verify request can be stubbed using mirage
    this.nonce = '12345';
    this.nonceStub = sinon.stub(uuid, 'default').returns(this.nonce);
    this.verifyResponse = { data: { correct_answer: 68 } };
    this.server.get(`/auth/:path/verify/${this.nonce}`, () => this.verifyResponse);

    // Auth request stub
    const api = this.owner.lookup('service:api');
    this.authenticateStub = sinon.stub(api.auth, 'oktaLogin');
    this.authResponse = RESPONSE_STUBS.okta;
    this.loginData = AUTH_METHOD_LOGIN_DATA.okta;
    // Resolve response by default, specific tests override this as needed
    this.authenticateStub.resolves(this.authResponse);

    this.assertSubmit = (assert, loginRequestArgs, loginData) => {
      const [username, path, { nonce, password }] = loginRequestArgs;
      // if path is included in loginData, a custom path was submitted
      const expectedPath = loginData?.path || this.authType;
      assert.strictEqual(path, expectedPath, 'it calls oktaLogin with expected path');
      assert.strictEqual(username, loginData.username, 'it calls oktaLogin with username');
      assert.strictEqual(password, loginData.password, 'it calls oktaLogin with password');
      assert.strictEqual(nonce, this.nonce, 'it calls oktaLogin with nonce');
    };

    this.renderComponent = ({ yieldBlock = false } = {}) => {
      if (yieldBlock) {
        return render(hbs`
          <Auth::Form::Okta 
            @authType={{this.authType}} 
            @cluster={{this.cluster}}
            @onError={{this.onError}}
            @onSuccess={{this.onSuccess}}
          >
            <:advancedSettings>
              <label for="path">Mount path</label>
              <input data-test-input="path" id="path" name="path" type="text" /> 
            </:advancedSettings>
          </Auth::Form::Okta>`);
      }
      return render(hbs`
      <Auth::Form::Okta       
        @authType={{this.authType}}
        @cluster={{this.cluster}}
        @onError={{this.onError}}
        @onSuccess={{this.onSuccess}}
      />`);
    };
  });

  hooks.afterEach(function () {
    this.authenticateStub.restore();
    this.nonceStub.restore();
  });

  authFormTestHelper(test);

  test('it displays okta number challenge answer', async function (assert) {
    await this.renderComponent();
    await fillInLoginFields(this.loginData);
    await click(GENERAL.submitButton);
    assert
      .dom('[data-test-okta-number-challenge]')
      .hasText(
        'To finish signing in, you will need to complete an additional MFA step. Okta verification Select the following number to complete verification: 68 Back to login'
      );
  });

  test('it returns to login when "Back to login" is clicked', async function (assert) {
    await this.renderComponent();
    await fillInLoginFields(this.loginData);
    await click(GENERAL.submitButton);
    assert.dom('[data-test-okta-number-challenge]').exists();
    await click(GENERAL.backButton);
    assert.dom(AUTH_FORM.authForm('okta')).exists('it returns to okta form');
    assert.dom('[data-test-okta-number-challenge]').doesNotExist();
    assert.dom(GENERAL.inputByAttr('username')).hasValue('', 'username is cleared');
    assert.dom(GENERAL.inputByAttr('password')).hasValue('', 'password is cleared');
  });

  test('it shows loading state when polling okta verify request', async function (assert) {
    assert.expect(2);
    let count = 0;
    this.server.get(`/auth/okta/verify/${this.nonce}`, () => {
      count++;
      assert.dom('[data-test-okta-number-challenge]').hasText(
        'To finish signing in, you will need to complete an additional MFA step. Please wait... Back to login',
        count === 1
          ? // the response hasn't returned anything yet, so the first assertion is just the initial state
            'it shows loading message before polling initiates'
          : // by now the response has returned a 404 so this asserts error handling works as expected
            'it shows loading message while response returns 404'
      );
      // okta/verify returns a 404 until the user interacts with okta via their configured MFA app.
      // to simulate this interaction we return data on the third request - which ends the polling.
      const response = count < 2 ? new Response(404) : this.verifyResponse;
      return response;
    });

    await this.renderComponent();
    await fillInLoginFields(this.loginData);
    await click(GENERAL.submitButton);
  });

  test('it renders error message when okta verify request errors', async function (assert) {
    this.server.get(`/auth/okta/verify/${this.nonce}`, () => new Response(500, {}, { errors: ['oh no!!'] }));
    await this.renderComponent();
    await fillInLoginFields(this.loginData);
    await click(GENERAL.submitButton);
    assert.dom(GENERAL.messageError).hasText('Error oh no!!');
  });
});
