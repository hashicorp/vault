/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, focus, triggerKeyEvent, typeIn, fillIn, click } from '@ember/test-helpers';
import { create } from 'ember-cli-page-object';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import maskedInput from 'vault/tests/pages/components/masked-input';

const component = create(maskedInput);

module('Integration | Component | masked input', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function (assert) {
    await render(hbs`<MaskedInput />`);
    assert.dom('[data-test-masked-input]').exists('shows masked input');
    assert.ok(component.textareaIsPresent);
    assert.dom('[data-test-textarea]').hasClass('masked-font', 'it renders an input with obscure font');
    assert.notOk(component.copyButtonIsPresent, 'does not render copy button by default');
    assert.notOk(component.downloadButtonIsPresent, 'does not render download button by default');

    await component.toggleMasked();
    assert.dom('.masked-value').doesNotHaveClass('masked-font', 'it unmasks when show button is clicked');
    await component.toggleMasked();
    assert.dom('.masked-value').hasClass('masked-font', 'it remasks text when button is clicked');
  });

  test('it renders correctly when displayOnly', async function (assert) {
    this.set('value', 'value');
    await render(hbs`<MaskedInput @displayOnly={{true}} @value={{this.value}} />`);

    assert.dom('.masked-value').hasClass('masked-font', 'value has obscured font');
    assert.notOk(component.textareaIsPresent, 'it does not render a textarea when displayOnly is true');
  });

  test('it renders a copy button when allowCopy is true', async function (assert) {
    await render(hbs`<MaskedInput @allowCopy={{true}} />`);
    assert.ok(component.copyButtonIsPresent);
  });

  test('it renders a download button when allowDownload is true', async function (assert) {
    await render(hbs`<MaskedInput @allowDownload={{true}} /> <div id="modal-wormhole"></div>
`);
    assert.ok(component.downloadIconIsPresent);

    await click('[data-test-download-icon]');
    assert.ok(component.downloadButtonIsPresent, 'clicking download icon opens modal with download button');

    assert;
  });

  test('it shortens all outputs when displayOnly and masked', async function (assert) {
    this.set('value', '123456789-123456789-123456789');
    await render(hbs`<MaskedInput @value={{this.value}} @displayOnly={{true}} />`);
    const maskedValue = document.querySelector('.masked-value').innerText;
    assert.strictEqual(maskedValue.length, 11);

    await component.toggleMasked();
    const unMaskedValue = document.querySelector('.masked-value').innerText;
    assert.strictEqual(unMaskedValue.length, this.value.length);
  });

  test('it does not unmask text on focus', async function (assert) {
    this.set('value', '123456789-123456789-123456789');
    await render(hbs`<MaskedInput @value={{this.value}} />`);
    assert.dom('.masked-value').hasClass('masked-font');
    await focus('.masked-value');
    assert.dom('.masked-value').hasClass('masked-font');
  });

  test('it calls onChange events with the correct values', async function (assert) {
    const changeSpy = sinon.spy();
    this.set('value', 'before');
    this.set('onChange', changeSpy);
    await render(hbs`<MaskedInput @value={{this.value}} @onChange={{this.onChange}} />`);
    await fillIn('[data-test-textarea]', 'after');
    assert.true(changeSpy.calledWith('after'));
  });

  test('it calls onKeyUp events with the correct values', async function (assert) {
    const keyupSpy = sinon.spy();
    this.set('value', '');
    this.set('onKeyUp', keyupSpy);
    await render(hbs`<MaskedInput @name="foo" @value={{this.value}} @onKeyUp={{this.onKeyUp}} />`);
    await typeIn('[data-test-textarea]', 'baz');
    assert.true(keyupSpy.calledThrice, 'calls for each letter of typing');
    assert.true(keyupSpy.firstCall.calledWithExactly('foo', 'b'));
    assert.true(keyupSpy.secondCall.calledWithExactly('foo', 'ba'));
    assert.true(keyupSpy.thirdCall.calledWithExactly('foo', 'baz'));
  });

  test('it does not remove value on tab', async function (assert) {
    this.set('value', 'hello');
    await render(hbs`<MaskedInput @value={{this.value}} />`);
    await triggerKeyEvent('[data-test-textarea]', 'keydown', 9);
    await component.toggleMasked();
    const unMaskedValue = document.querySelector('.masked-value').value;
    assert.strictEqual(unMaskedValue, this.value);
  });

  test('it should render stringify toggle in download modal', async function (assert) {
    assert.expect(3);

    // this looks wonky but need a new line in there to test stringify adding escape character
    this.value = `bar
`;

    const downloadStub = sinon.stub(this.owner.lookup('service:download'), 'miscExtension');
    downloadStub.callsFake((fileName, value) => {
      const firstCall = downloadStub.callCount === 1;
      const assertVal = firstCall ? this.value : JSON.stringify(this.value);
      assert.strictEqual(assertVal, value, `Value is ${firstCall ? 'not ' : ''}stringified`);
      return true;
    });

    await render(hbs`
      <MaskedInput
        @name="key"
        @value={{this.value}}
        @displayOnly={{true}}
        @allowDownload={{true}}
      />
    `);

    await click('[data-test-download-icon]');
    assert.dom('[data-test-stringify-toggle]').isNotChecked('Stringify toggle off as default');
    await click('[data-test-download-button]');

    await click('[data-test-download-icon]');
    await click('[data-test-stringify-toggle]');
    await click('[data-test-download-button]');
  });
});
