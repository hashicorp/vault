/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import hbs from 'htmlbars-inline-precompile';
import { find, render } from '@ember/test-helpers';
import sinon from 'sinon';
import testHelper from './test-helper';
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

  module('github', function (hooks) {
    hooks.beforeEach(function () {
      this.authType = 'github';
      this.expectedFields = ['token'];
      this.renderComponent = () => {
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
      this.renderComponent = () => {
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
      this.renderComponent = () => {
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
      this.renderComponent = () => {
        return render(hbs`
          <Auth::Form::Token       
            @authType={{this.authType}}
            @cluster={{this.cluster}}
            @onError={{this.onError}}
            @onSuccess={{this.onSuccess}}
          />`);
      };
    });

    testHelper(test);
  });

  module('userpass', function (hooks) {
    hooks.beforeEach(function () {
      this.authType = 'userpass';
      this.expectedFields = ['username', 'password'];
      this.renderComponent = () => {
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
