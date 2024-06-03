/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import sinon from 'sinon';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { TOOLS_SELECTORS as TS } from 'vault/tests/helpers/tools-selectors';

module('Integration | Component | tools/rewrap', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.onClear = sinon.spy();
    this.rewrap_token = null;
    this.token = null;
    this.errors = null;
    this.renderComponent = async () => {
      await render(hbs`
    <Tools::Rewrap
      @onClear={{this.onClear}}
      @rewrap_token={{this.rewrap_token}}
      @errors={{this.errors}}
      @token={{this.token}}
    />`);
    };
  });

  test('it renders defaults', async function (assert) {
    await this.renderComponent();

    assert.dom('h1').hasText('Rewrap Token', 'title renders');
    assert.dom(TS.submit).hasText('Rewrap token');
    assert.dom(TS.toolsInput('wrapping-token')).hasValue('');
    assert.dom(TS.toolsInput('rewrapped-token')).doesNotExist();
    assert.dom(TS.button('Done')).doesNotExist();
  });

  test('it renders errors', async function (assert) {
    this.errors = ['Something is wrong!'];
    await this.renderComponent();
    assert.dom(GENERAL.messageError).hasText('Error Something is wrong!', 'Error renders');
  });

  test('it renders random bytes view', async function (assert) {
    this.rewrap_token = 'blah.CAESIK19SQeLYUZ65lEGSMYYOMZFbUurY0ppT2RTMGpRa0JOSUFqUzJUaGNqdWUQ6ooG';
    await this.renderComponent();

    assert.dom('label').hasText('Rewrapped token');
    assert.dom(TS.toolsInput('wrapping-token')).doesNotExist();
    assert.dom(TS.toolsInput('rewrapped-token')).hasText(this.rewrap_token);

    await click(TS.button('Done'));
    assert.true(this.onClear.calledOnce, 'onClear is called');
  });

  test('it updates arg when input changes', async function (assert) {
    await this.renderComponent();
    await fillIn(TS.toolsInput('wrapping-token'), 'my-token');
    assert.strictEqual(this.token, 'my-token', 'token value updates when input changes');
  });
});
