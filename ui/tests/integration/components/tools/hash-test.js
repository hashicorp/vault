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

module('Integration | Component | tools/hash', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.onClear = sinon.spy();
    this.onChange = sinon.spy();
    this.sum = null;
    this.hashData = '';
    this.algorithm = 'sha2-256';
    this.format = 'base64';
    this.renderComponent = async () => {
      await render(hbs`
    <Tools::Hash
      @onClear={{this.onClear}}
      @onChange={{this.onChange}}
      @errors={{this.errors}}
      @sum={{this.sum}}
      @algorithm={{this.algorithm}}
      @format={{this.format}}
      @hashData={{this.hashData}}
    />`);
    };
  });

  test('it renders defaults', async function (assert) {
    await this.renderComponent();

    assert.dom('h1').hasText('Hash Data', 'Title renders');
    assert.dom('#algorithm').hasValue('sha2-256');
    assert.dom('#format').hasValue('base64');
    assert.dom(TS.toolsInput('sum')).doesNotExist();
    assert.dom(TS.button('Done')).doesNotExist();
  });

  test('it renders errors', async function (assert) {
    this.errors = ['Something is wrong!'];
    await this.renderComponent();
    assert.dom(GENERAL.messageError).hasText('Error Something is wrong!', 'Error renders');
  });

  test('it renders sum view', async function (assert) {
    this.sum = '2a70350c9a44171d6b1180c6be5cbb2ee3f79d532c8a1dd9ef2e8e08e752a3babb';
    await this.renderComponent();

    assert.dom('h1').hasText('Hash Data');
    assert.dom('label').hasText('Sum');
    assert.dom('#algorithm').doesNotExist();
    assert.dom('#format').doesNotExist();
    assert.dom(TS.toolsInput('sum')).hasText(this.sum);
    await click(TS.button('Done'));
    assert.true(this.onClear.calledOnce, 'onClear is called');
  });

  test('it calls onChange when inputs change', async function (assert) {
    await this.renderComponent();
    await fillIn(TS.toolsInput('hash-input'), 'foo');
    assert.propEqual(this.onChange.lastCall.args, ['hashData', 'foo'], 'onChange is called with hash input');

    await fillIn('#algorithm', 'sha2-224');
    assert.propEqual(
      this.onChange.lastCall.args,
      ['algorithm', 'sha2-224'],
      'onChange is called with algorithm'
    );
    await fillIn('#format', 'hex');
    assert.propEqual(this.onChange.lastCall.args, ['format', 'hex'], 'onChange is called with format');
  });
});
