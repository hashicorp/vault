/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, focus, triggerKeyEvent } from '@ember/test-helpers';
import { create } from 'ember-cli-page-object';
import hbs from 'htmlbars-inline-precompile';
import maskedInput from 'vault/tests/pages/components/masked-input';

const component = create(maskedInput);

module('Integration | Component | masked input', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function (assert) {
    await render(hbs`{{masked-input}}`);
    assert.dom('[data-test-masked-input]').exists('shows expiration beacon');
  });

  test('it renders a textarea', async function (assert) {
    await render(hbs`{{masked-input}}`);

    assert.ok(component.textareaIsPresent);
  });

  test('it renders an input with obscure font', async function (assert) {
    await render(hbs`{{masked-input}}`);

    assert.dom('[data-test-textarea]').hasClass('masked-font', 'loading class with correct font');
  });

  test('it renders obscure font when displayOnly', async function (assert) {
    this.set('value', 'value');
    await render(hbs`{{masked-input displayOnly=true value=this.value}}`);

    assert.dom('.masked-value').hasClass('masked-font', 'loading class with correct font');
  });

  test('it does not render a textarea when displayOnly is true', async function (assert) {
    await render(hbs`{{masked-input displayOnly=true}}`);

    assert.notOk(component.textareaIsPresent);
  });

  test('it renders a copy button when allowCopy is true', async function (assert) {
    await render(hbs`{{masked-input allowCopy=true}}`);

    assert.ok(component.copyButtonIsPresent);
  });

  test('it does not render a copy button when allowCopy is false', async function (assert) {
    await render(hbs`{{masked-input allowCopy=false}}`);

    assert.notOk(component.copyButtonIsPresent);
  });

  test('it unmasks text when button is clicked', async function (assert) {
    this.set('value', 'value');
    await render(hbs`{{masked-input value=this.value}}`);
    await component.toggleMasked();

    assert.dom('.masked-value').doesNotHaveClass('masked-font');
  });

  test('it remasks text when button is clicked', async function (assert) {
    this.set('value', 'value');
    await render(hbs`{{masked-input value=this.value}}`);

    await component.toggleMasked();
    await component.toggleMasked();

    assert.dom('.masked-value').hasClass('masked-font');
  });

  test('it shortens all outputs when displayOnly and masked', async function (assert) {
    this.set('value', '123456789-123456789-123456789');
    await render(hbs`{{masked-input value=this.value displayOnly=true}}`);
    const maskedValue = document.querySelector('.masked-value').innerText;
    assert.strictEqual(maskedValue.length, 11);

    await component.toggleMasked();
    const unMaskedValue = document.querySelector('.masked-value').innerText;
    assert.strictEqual(unMaskedValue.length, this.value.length);
  });

  test('it does not unmask text on focus', async function (assert) {
    this.set('value', '123456789-123456789-123456789');
    await render(hbs`{{masked-input value=this.value}}`);
    assert.dom('.masked-value').hasClass('masked-font');
    await focus('.masked-value');
    assert.dom('.masked-value').hasClass('masked-font');
  });

  test('it does not remove value on tab', async function (assert) {
    this.set('value', 'hello');
    await render(hbs`{{masked-input value=this.value}}`);
    await triggerKeyEvent('[data-test-textarea]', 'keydown', 9);
    await component.toggleMasked();
    const unMaskedValue = document.querySelector('.masked-value').value;
    assert.strictEqual(unMaskedValue, this.value);
  });
});
