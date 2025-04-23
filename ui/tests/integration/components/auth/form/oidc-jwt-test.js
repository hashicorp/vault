/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import hbs from 'htmlbars-inline-precompile';
import { find, render } from '@ember/test-helpers';
import sinon from 'sinon';
import { setupMirage } from 'ember-cli-mirage/test-support';
import testHelper from './test-helper';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | auth | form | oidc-jwt', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.expectedFields = ['role'];

    this.authenticateStub = sinon.stub(this.owner.lookup('service:auth'), 'authenticate');
    this.cluster = { id: 1 };
    this.onError = sinon.spy();
    this.onSuccess = sinon.spy();
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
      this.expectedSubmit = {
        default: { path: 'oidc', role: 'some-dev' },
        custom: { path: 'custom-oidc', role: 'some-dev' },
      };
    });

    testHelper(test);
  });

  module('jwt', function (hooks) {
    hooks.beforeEach(function () {
      this.authType = 'jwt';
      this.expectedSubmit = {
        default: { path: 'jwt', role: 'some-dev' },
        custom: { path: 'custom-jwt', role: 'some-dev' },
      };
    });

    testHelper(test);
  });
});
