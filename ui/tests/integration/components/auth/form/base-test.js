/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import hbs from 'htmlbars-inline-precompile';
import { click, fillIn, find, render } from '@ember/test-helpers';
import sinon from 'sinon';
import testHelper from './auth-form-test-helper';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

// These auth types all use the default methods in auth/form/base
// Any auth types with custom logic should be in a separate test file, i.e. okta

module('Integration | Component | auth | form | base', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.authenticateStub = sinon.stub(this.owner.lookup('service:auth'), 'authenticate');
    this.cluster = { id: 1 };
    this.onError = sinon.spy();
    this.onSuccess = sinon.spy();
  });

  hooks.afterEach(function () {
    this.authenticateStub.restore();
  });

  module('github', function (hooks) {
    hooks.beforeEach(function () {
      this.authType = 'github';
      this.expectedFields = ['token'];
      this.expectedSubmit = {
        default: { path: 'github', token: 'mysupersecuretoken' },
        custom: { path: 'custom-github', token: 'mysupersecuretoken' },
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

    testHelper(test);

    test('it renders custom label', async function (assert) {
      await this.renderComponent();
      const id = find(GENERAL.inputByAttr('token')).id;
      assert.dom(`#label-${id}`).hasText('Github token');
    });
  });

  module('ldap', function (hooks) {
    hooks.beforeEach(function () {
      this.authType = 'ldap';
      this.expectedFields = ['username', 'password'];
      this.expectedSubmit = {
        default: { password: 'password', path: 'ldap', username: 'matilda' },
        custom: { password: 'password', path: 'custom-ldap', username: 'matilda' },
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

    testHelper(test);
  });

  module('radius', function (hooks) {
    hooks.beforeEach(function () {
      this.authType = 'radius';
      this.expectedFields = ['username', 'password'];
      this.expectedSubmit = {
        default: { password: 'password', path: 'radius', username: 'matilda' },
        custom: { password: 'password', path: 'custom-radius', username: 'matilda' },
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

    testHelper(test);
  });

  module('token', function (hooks) {
    hooks.beforeEach(function () {
      this.authType = 'token';
      this.expectedFields = ['token'];
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

    testHelper(test, { standardSubmit: false });

    test('it submits form data with defaults', async function (assert) {
      await this.renderComponent();
      await fillIn(GENERAL.inputByAttr('token'), 'mytoken');
      await click(GENERAL.submitButton);
      const [actual] = this.authenticateStub.lastCall.args;
      assert.propEqual(
        actual.data,
        { token: 'mytoken' },
        'auth service "authenticate" method is called with form data'
      );
    });

    test('it submits form data from yielded inputs', async function (assert) {
      await this.renderComponent({ yieldBlock: true });
      await fillIn(GENERAL.inputByAttr('token'), 'mytoken');
      // token doesn't support custom paths, so testing path is not sent
      await fillIn(GENERAL.inputByAttr('path'), `path-${this.authType}`);
      await click(GENERAL.submitButton);
      const [actual] = this.authenticateStub.lastCall.args;
      assert.propEqual(
        actual.data,
        { token: 'mytoken' },
        'auth service "authenticate" method is called without "path"'
      );
    });
  });

  module('userpass', function (hooks) {
    hooks.beforeEach(function () {
      this.authType = 'userpass';
      this.expectedFields = ['username', 'password'];
      this.expectedSubmit = {
        default: { password: 'password', path: 'userpass', username: 'matilda' },
        custom: { password: 'password', path: 'custom-userpass', username: 'matilda' },
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

    testHelper(test);
  });
});
