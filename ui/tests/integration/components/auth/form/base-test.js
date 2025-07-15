/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import hbs from 'htmlbars-inline-precompile';
import { find, render } from '@ember/test-helpers';
import sinon from 'sinon';
import authFormTestHelper from './auth-form-test-helper';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { AUTH_METHOD_LOGIN_DATA } from 'vault/tests/helpers/auth/auth-helpers';
import { RESPONSE_STUBS } from 'vault/tests/helpers/auth/response-stubs';

// These auth types all use the default methods in auth/form/base
// Any auth types with custom logic should be in a separate test file, i.e. okta

module('Integration | Component | auth | form | base', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.cluster = { id: 1 };
    this.onError = sinon.spy();
    this.onSuccess = sinon.spy();
    const api = this.owner.lookup('service:api');
    this.setup = (authType, loginMethod) => {
      this.authType = authType;
      this.authenticateStub = sinon.stub(api.auth, loginMethod);
      this.authResponse = RESPONSE_STUBS[authType];
      this.loginData = AUTH_METHOD_LOGIN_DATA[authType];
    };
  });

  module('github', function (hooks) {
    hooks.beforeEach(function () {
      this.setup('github', 'githubLogin');
      this.assertSubmit = (assert, loginRequestArgs, loginData) => {
        const [path, { token }] = loginRequestArgs;
        // if path is included in loginData, a custom path was submitted
        const expectedPath = loginData?.path || this.authType;
        assert.strictEqual(path, expectedPath, 'it calls githubLogin with expected path');
        assert.strictEqual(token, loginData.token, 'it calls githubLogin with token');
      };
      this.renderComponent = ({ yieldBlock = false } = {}) => {
        if (yieldBlock) {
          return render(hbs`
            <Auth::Form::Github 
              @authType={{this.authType}} 
              @cluster={{this.cluster}}
              @onError={{this.onError}}
              @onSuccess={{this.onSuccess}}
            >
             <:advancedSettings>
             <label for="path">Mount path</label>
             <input data-test-input="path" id="path" name="path" type="text" /> 
             </:advancedSettings>
            </Auth::Form::Github>`);
        }
        return render(hbs`
          <Auth::Form::Github       
            @authType={{this.authType}}
            @cluster={{this.cluster}}
            @onError={{this.onError}}
            @onSuccess={{this.onSuccess}}
          />`);
      };
    });

    hooks.afterEach(function () {
      this.authenticateStub.restore();
    });

    authFormTestHelper(test);

    test('it renders custom label', async function (assert) {
      await this.renderComponent();
      const id = find(GENERAL.inputByAttr('token')).id;
      assert.dom(`#label-${id}`).hasText('Github token');
    });
  });

  module('ldap', function (hooks) {
    hooks.beforeEach(function () {
      this.setup('ldap', 'ldapLogin');
      this.assertSubmit = (assert, loginRequestArgs, loginData) => {
        const [username, path, { password }] = loginRequestArgs;
        // if path is included in loginData, a custom path was submitted
        const expectedPath = loginData?.path || this.authType;
        assert.strictEqual(path, expectedPath, 'it calls ldapLogin with expected path');
        assert.strictEqual(username, loginData.username, 'it calls ldapLogin with username');
        assert.strictEqual(password, loginData.password, 'it calls ldapLogin with password');
      };
      this.renderComponent = ({ yieldBlock = false } = {}) => {
        if (yieldBlock) {
          return render(hbs`
            <Auth::Form::Ldap 
              @authType={{this.authType}} 
              @cluster={{this.cluster}}
              @onError={{this.onError}}
              @onSuccess={{this.onSuccess}}
            >
             <:advancedSettings>
             <label for="path">Mount path</label>
             <input data-test-input="path" id="path" name="path" type="text" /> 
             </:advancedSettings>
            </Auth::Form::Ldap>`);
        }
        return render(hbs`
          <Auth::Form::Ldap       
            @authType={{this.authType}}
            @cluster={{this.cluster}}
            @onError={{this.onError}}
            @onSuccess={{this.onSuccess}}
          />`);
      };
    });

    hooks.afterEach(function () {
      this.authenticateStub.restore();
    });

    authFormTestHelper(test);
  });

  module('radius', function (hooks) {
    hooks.beforeEach(function () {
      this.setup('radius', 'radiusLoginWithUsername');
      this.assertSubmit = (assert, loginRequestArgs, loginData) => {
        const [username, path, { password }] = loginRequestArgs;
        // if path is included in loginData, a custom path was submitted
        const expectedPath = loginData?.path || this.authType;
        assert.strictEqual(username, loginData.username, 'it calls radiusLoginWithUsername with username');
        assert.strictEqual(path, expectedPath, 'it calls radiusLoginWithUsername with expected path');
        assert.strictEqual(password, loginData.password, 'it calls radiusLoginWithUsername with password');
      };
      this.renderComponent = ({ yieldBlock = false } = {}) => {
        if (yieldBlock) {
          return render(hbs`
            <Auth::Form::Radius 
              @authType={{this.authType}} 
              @cluster={{this.cluster}}
              @onError={{this.onError}}
              @onSuccess={{this.onSuccess}}
            >
             <:advancedSettings>
              <label for="path">Mount path</label>
              <input data-test-input="path" id="path" name="path" type="text" /> 
             </:advancedSettings>
            </Auth::Form::Radius>`);
        }
        return render(hbs`
          <Auth::Form::Radius       
            @authType={{this.authType}}
            @cluster={{this.cluster}}
            @onError={{this.onError}}
            @onSuccess={{this.onSuccess}}
          />`);
      };
    });

    hooks.afterEach(function () {
      this.authenticateStub.restore();
    });

    authFormTestHelper(test);
  });

  module('token', function (hooks) {
    hooks.beforeEach(function () {
      this.setup('token', 'tokenLookUpSelf');
      this.assertSubmit = (assert, loginRequestArgs) => {
        const [{ headers }] = loginRequestArgs;
        assert.strictEqual(headers['X-Vault-Token'], 'mysupersecuretoken', 'token is submitted as header');
      };
      this.renderComponent = ({ yieldBlock = false } = {}) => {
        if (yieldBlock) {
          return render(hbs`
            <Auth::Form::Token 
              @authType={{this.authType}} 
              @cluster={{this.cluster}}
              @onError={{this.onError}}
              @onSuccess={{this.onSuccess}}
            >
             <:advancedSettings>
                <label for="path">Mount path</label>
                <input data-test-input="path" id="path" name="path" type="text" /> 
             </:advancedSettings>
            </Auth::Form::Token>`);
        }
        return render(hbs`
          <Auth::Form::Token       
            @authType={{this.authType}}
            @cluster={{this.cluster}}
            @onError={{this.onError}}
            @onSuccess={{this.onSuccess}}
          />`);
      };
    });

    hooks.afterEach(function () {
      this.authenticateStub.restore();
    });

    authFormTestHelper(test);
  });

  module('userpass', function (hooks) {
    hooks.beforeEach(function () {
      this.setup('userpass', 'userpassLogin');
      this.assertSubmit = (assert, loginRequestArgs, loginData) => {
        const [username, path, { password }] = loginRequestArgs;
        // if path is included in loginData, a custom path was submitted
        const expectedPath = loginData?.path || this.authType;
        assert.strictEqual(path, expectedPath, 'it calls userpassLogin with expected path');
        assert.strictEqual(username, loginData.username, 'it calls userpassLogin with username');
        assert.strictEqual(password, loginData.password, 'it calls userpassLogin with password');
      };
      this.renderComponent = ({ yieldBlock = false } = {}) => {
        if (yieldBlock) {
          return render(hbs`
            <Auth::Form::Userpass 
              @authType={{this.authType}} 
              @cluster={{this.cluster}}
              @onError={{this.onError}}
              @onSuccess={{this.onSuccess}}
            >
             <:advancedSettings>
              <label for="path">Mount path</label>
              <input data-test-input="path" id="path" name="path" type="text" /> 
             </:advancedSettings>
            </Auth::Form::Userpass>`);
        }
        return render(hbs`
          <Auth::Form::Userpass       
            @authType={{this.authType}}
            @cluster={{this.cluster}}
            @onError={{this.onError}}
            @onSuccess={{this.onSuccess}}
          />`);
      };
    });

    hooks.afterEach(function () {
      this.authenticateStub.restore();
    });

    authFormTestHelper(test);
  });
});
