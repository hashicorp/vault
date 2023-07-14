import { module, test } from 'qunit';
import sinon from 'sinon';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, render, typeIn } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

const selectors = {
  input: '[data-test-shamir-key-input]',
  inputLabel: '[data-test-shamir-key-label]',
  submitButton: '[data-test-shamir-submit]',
  otpInfo: '[data-test-otp-info]',
  otpCode: '[data-test-otp]',
};

module('Integration | Component | shamir/form', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.submitSpy = sinon.spy();
  });

  test('it does calls callback only if key value present', async function (assert) {
    await render(hbs`
      <Shamir::Form @onSubmit={{this.submitSpy}} @progress={{2}} @threshold={{4}} />
    `);
    assert.dom(selectors.submitButton).hasText('Submit', 'Submit button has default text');
    await click(selectors.submitButton);
    assert.ok(this.submitSpy.notCalled, 'onSubmit was not called');
    await typeIn(selectors.input, 'this-is-the-key');
    assert.dom(selectors.input).hasValue('this-is-the-key', 'input value set');
    assert.dom(selectors.inputLabel).hasText('Shamir key portion', 'label has default text');
    await click(selectors.submitButton);
    assert.ok(
      this.submitSpy.calledOnceWith({ key: 'this-is-the-key' }),
      'onSubmit called with correct params'
    );
    assert.dom(selectors.input).hasValue('', 'key value reset after submit');

    // Template block usage:
    await render(hbs`
    <Shamir::Form @onSubmit={{this.submitSpy}} @progress={{2}} @threshold={{4}} @buttonText="Do the thing" @inputLabel="Unseal key">
      <div data-test-block-content>Hello</div>
    </Shamir::Form>
    `);

    assert.dom('[data-test-block-content]').hasText('Hello', 'renders block content');
    assert.dom(selectors.submitButton).hasText('Do the thing', 'uses passed button text');
    assert.dom(selectors.inputLabel).hasText('Unseal key', 'uses passed inputLabel');
    assert.dom(selectors.otpInfo).doesNotExist('no OTP info shown');
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

    assert.dom(selectors.otpInfo).exists('shows OTP info');
    assert.dom(selectors.otpCode).hasText('this-is-otp', 'shows OTP code');
    assert.dom('[data-test-block-content]').hasText('Hello', 'renders block content');
  });
});
