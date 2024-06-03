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

module('Integration | Component | tools/random', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.onClear = sinon.spy();
    this.onChange = sinon.spy();
    this.random_bytes = null;
    this.bytes = '32';
    this.format = 'base64';
    this.renderComponent = async () => {
      await render(hbs`
    <Tools::Random
      @onClear={{this.onClear}}
      @onChange={{this.onChange}}
      @random_bytes={{this.random_bytes}}
      @errors={{this.errors}}
      @format={{this.format}}
      @bytes={{this.bytes}}
    />`);
    };
  });

  test('it renders defaults', async function (assert) {
    await this.renderComponent();

    assert.dom('h1').hasText('Random Bytes', 'Title renders');
    assert.dom('#bytes').hasValue('32');
    assert.dom('#format').hasValue('base64');
    assert.dom(TS.toolsInput('random-bytes')).doesNotExist();
    assert.dom(TS.button('Done')).doesNotExist();
  });

  test('it renders errors', async function (assert) {
    this.errors = ['Something is wrong!'];
    await this.renderComponent();
    assert.dom(GENERAL.messageError).hasText('Error Something is wrong!', 'Error renders');
  });

  test('it renders random bytes view', async function (assert) {
    this.random_bytes = '1aG60EHnSz8zn+6qqSA+ulUSJylhjCZP30t7n21Qjro=';
    await this.renderComponent();

    assert.dom('label').hasText('Random bytes');
    assert.dom('#bytes').doesNotExist();
    assert.dom('#format').doesNotExist();
    assert.dom(TS.toolsInput('random-bytes')).hasText(this.random_bytes);
    await click(TS.button('Done'));
    assert.true(this.onClear.calledOnce, 'onClear is called');
  });

  test('it calls onChange when inputs change', async function (assert) {
    await this.renderComponent();
    await fillIn(TS.toolsInput('bytes'), '43');
    assert.strictEqual(this.bytes, '43', 'bytes update when input changes');

    await fillIn('#format', 'hex');
    assert.propEqual(this.onChange.lastCall.args, ['format', 'hex'], 'onChange is called with format');
  });
});
