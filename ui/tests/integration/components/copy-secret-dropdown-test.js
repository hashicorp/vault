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
  masked: '[data-test-masked-input]',
};

module('Integration | Component | copy-secret-dropdown', function (hooks) {
  setupRenderingTest(hooks);
  hooks.beforeEach(function () {
    this.onWrap = sinon.stub().returns('wrapped-data');
    this.onClose = () => {};
    this.data = `{ foo: 'bar' }`;
    this.wrappedData = 'wrapped-data';
    this.clipboardSpy = sinon.stub(navigator.clipboard, 'writeText').resolves();
  });

  hooks.afterEach(function () {
    sinon.restore(); // resets all stubs, including clipboard
  });

  test('it renders and fires callback functions', async function (assert) {
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
    assert.dom(GENERAL.copyButton).hasText('Copy JSON');
    assert.dom(GENERAL.button('wrap')).hasText('Wrap secret');

    await click(GENERAL.button('wrap'));
    await click(SELECTORS.dropdown);
  });

  test('it copies JSON', async function (assert) {
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
    assert.true(this.clipboardSpy.calledOnce, 'Clipboard was called once');
    assert.strictEqual(this.clipboardSpy.firstCall.args[0], this.data, 'copy value is the json secret data');
  });

  test('it renders loading wrap button', async function (assert) {
    await render(
      hbs`<ToolbarActions>
  <CopySecretDropdown
    @clipboardText={{this.data}}
    @onWrap={{this.onWrap}}
    @isWrapping={{true}}
    @onClose={{this.onClose}}
  />
</ToolbarActions>`
    );
    await click(SELECTORS.dropdown);
    assert.dom(`[data-test-icon="loading"]`).exists('renders loading icon');
  });

  test('it wraps data', async function (assert) {
    this.data = '';
    await render(
      hbs`<ToolbarActions>
  <CopySecretDropdown
    @clipboardText={{this.data}}
    @onWrap={{this.onWrap}}
    @isWrapping={{false}}
    @onClose={{this.onClose}}
  />
</ToolbarActions>`
    );

    await click(SELECTORS.dropdown);
    await click(GENERAL.button('wrap'));
    assert.true(this.onWrap.calledOnce, 'onWrap was called');
  });
});
