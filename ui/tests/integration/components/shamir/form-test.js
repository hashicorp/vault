/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import sinon from 'sinon';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, render, settled, typeIn } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { SHAMIR_FORM } from 'vault/tests/helpers/components/shamir-selectors';

module('Integration | Component | shamir/form', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.submitSpy = sinon.spy();
  });

  test('it does calls callback only if key value present', async function (assert) {
    await render(hbs`
      <Shamir::Form @onSubmit={{this.submitSpy}} @progress={{0}} @threshold={{3}} />
    `);
    assert.dom(SHAMIR_FORM.submitButton).hasText('Submit', 'Submit button has default text');
    await click(SHAMIR_FORM.submitButton);
    assert.dom(SHAMIR_FORM.progress).doesNotExist('Hides progress bar if none made');
    assert.ok(this.submitSpy.notCalled, 'onSubmit was not called');
    await typeIn(SHAMIR_FORM.input, 'this-is-the-key');
    assert.dom(SHAMIR_FORM.input).hasValue('this-is-the-key', 'input value set');
    assert.dom(SHAMIR_FORM.inputLabel).hasText('Shamir key portion', 'label has default text');
    await click(SHAMIR_FORM.submitButton);
    assert.ok(
      this.submitSpy.calledOnceWith({ key: 'this-is-the-key' }),
      'onSubmit called with correct params'
    );
    assert.dom(SHAMIR_FORM.input).hasValue('', 'key value reset after submit');

    await render(hbs`
    <Shamir::Form @onSubmit={{this.submitSpy}} @progress={{0}} @threshold={{3}} @alwaysShowProgress={{true}} @buttonText="Do the thing" @inputLabel="Unseal key">
      <div data-test-block-content>Hello</div>
    </Shamir::Form>
    `);

    assert.dom('[data-test-block-content]').hasText('Hello', 'renders block content');
    assert.dom(SHAMIR_FORM.submitButton).hasText('Do the thing', 'uses passed button text');
    assert.dom(SHAMIR_FORM.inputLabel).hasText('Unseal key', 'uses passed inputLabel');
    assert.dom(SHAMIR_FORM.otpInfo).doesNotExist('no OTP info shown');
    assert
      .dom(SHAMIR_FORM.progress)
      .hasText('0/3 keys provided', 'displays textual progress when alwaysShowProgress=true');
  });

  test('it shows OTP info if provided', async function (assert) {
    await render(hbs`
      <Shamir::Form
        @onSubmit={{this.submitSpy}}
        @progress={{2}}
        @threshold={{4}}
        @otp="this-is-otp"
        @inputLabel="Please input key"
      >
        <div data-test-block-content>Hello</div>
      </Shamir::Form>
    `);

    assert.dom(SHAMIR_FORM.otpInfo).exists('shows OTP info');
    assert.dom(SHAMIR_FORM.otpCode).hasText('this-is-otp', 'shows OTP code');
    assert.dom(SHAMIR_FORM.progress).hasText('2/4 keys provided', 'displays textual progress');
    assert.dom('[data-test-block-content]').hasText('Hello', 'renders block content');
  });

  test('renders errors provided', async function (assert) {
    this.set('errors', ['first error', 'this is fine']);
    await render(hbs`
      <Shamir::Form
        @onSubmit={{this.submitSpy}}
        @errors={{this.errors}}
      />
    `);
    assert.dom(SHAMIR_FORM.error).exists({ count: 2 }, 'renders errors');

    this.set('errors', []);
    await settled();
    assert.dom(SHAMIR_FORM.error).doesNotExist('errors cleared');
  });
});
