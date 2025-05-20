/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';

const SELECTORS = {
  dropdown: '[data-test-copy-menu-trigger]',
  copyButton: GENERAL.copyButton,
  clipboard: 'data-test-copy-button',
  wrapButton: '[data-test-wrap-button]',
  masked: '[data-test-masked-input]',
};

module('Integration | Component | copy-secret-dropdown', function (hooks) {
  setupRenderingTest(hooks);
  hooks.beforeEach(function () {
    this.onWrap = () => {};
    this.onClose = () => {};
  });

  test('it renders and fires callback functions', async function (assert) {
    this.data = `{ foo: 'bar' }`;
    this.onWrap = () => assert.ok(true, 'onWrap callback fires Wrap secret is clicked');
    this.onClose = () => assert.ok(true, 'onClose callback fires when dropdown closes');

    await render(
      hbs`<ToolbarActions>
  <CopySecretDropdown
    @clipboardText={{this.data}}
    @onWrap={{this.onWrap}}
    @onClose={{this.onClose}}
  />
</ToolbarActions>`
    );

    await click(SELECTORS.dropdown);
    assert.dom(SELECTORS.copyButton).hasText('Copy JSON');
    assert.dom(SELECTORS.wrapButton).hasText('Wrap secret');

    await click(SELECTORS.wrapButton);
    await click(SELECTORS.dropdown);
  });

  test('it copies JSON', async function (assert) {
    // Sinon spy for clipboard
    const clipboardSpy = sinon.stub(navigator.clipboard, 'writeText').resolves();

    this.data = `{ foo: 'bar' }`;
    this.onWrap = () => assert.ok(true, 'onWrap callback fires Wrap secret is clicked');
    this.onClose = () => assert.ok(true, 'onClose callback fires when dropdown closes');

    await render(
      hbs`<ToolbarActions>
  <CopySecretDropdown
    @clipboardText={{this.data}}
    @onWrap={{this.onWrap}}
    @onClose={{this.onClose}}
  />
</ToolbarActions>`
    );

    await click(SELECTORS.dropdown);
    await click(GENERAL.copyButton);
    assert.true(clipboardSpy.calledOnce, 'Clipboard was called once');
    assert.strictEqual(clipboardSpy.firstCall.args[0], this.data, 'copy value is the json secret data');
    // Restore original clipboard
    clipboardSpy.restore(); // cleanup
  });

  test('it renders loading wrap button', async function (assert) {
    assert.expect(2);
    await render(
      hbs`<ToolbarActions>
  <CopySecretDropdown
    @clipboardText={{this.data}}
    @onWrap={{this.onWrap}}
    @isWrapping={{true}}
    @wrappedData={{this.wrappedData}}
    @onClose={{this.onClose}}
  />
</ToolbarActions>`
    );

    await click(SELECTORS.dropdown);
    assert.dom(`${SELECTORS.wrapButton} [data-test-icon="loading"]`).exists('renders loading icon');
    assert.dom(SELECTORS.wrapButton).isDisabled();
  });

  test('it wraps data', async function (assert) {
    // Sinon spy for clipboard
    const clipboardSpy = sinon.stub(navigator.clipboard, 'writeText').resolves();
    this.wrappedData = 'my-token';

    await render(
      hbs`<ToolbarActions>
  <CopySecretDropdown
    @clipboardText={{this.data}}
    @onWrap={{this.onWrap}}
    @isWrapping={{false}}
    @wrappedData={{this.wrappedData}}
    @onClose={{this.onClose}}
  />
</ToolbarActions>`
    );

    await click(SELECTORS.dropdown);
    await click(SELECTORS.wrapButton);
    assert.true(clipboardSpy.calledOnce, 'Clipboard was called once');
    assert.strictEqual(clipboardSpy.firstCall.args[0], this.wrappedData, 'copy value is wrapped secret data');
    // Restore original clipboard
    clipboardSpy.restore(); // cleanup
  });
});
