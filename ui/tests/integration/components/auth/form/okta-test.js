/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import hbs from 'htmlbars-inline-precompile';
import { render } from '@ember/test-helpers';
import sinon from 'sinon';
import { setupMirage } from 'ember-cli-mirage/test-support';
import testHelper from './test-helper';

module('Integration | Component | auth | form | okta', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.authType = 'okta';
    this.expectedFields = ['username', 'password'];

    this.authenticateStub = sinon.stub(this.owner.lookup('service:auth'), 'authenticate');
    this.cluster = { id: 1 };
    this.onError = sinon.spy();
    this.onSuccess = sinon.spy();
    this.renderComponent = () => {
      return render(hbs`
      <Auth::Form::Okta       
        @authType={{this.authType}}
        @cluster={{this.cluster}}
        @onError={{this.onError}}
        @onSuccess={{this.onSuccess}}
      />`);
    };
  });

  testHelper(test);
});
