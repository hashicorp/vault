/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

const SELECTORS = {
  dropdown: '[data-test-copy-menu-trigger]',
  copyButton: '[data-test-copy-button]',
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
    assert.expect(5);
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
    assert
      .dom(SELECTORS.copyButton)
      .hasAttribute('data-test-copy-button', `${this.data}`, 'it renders copyable data');

    await click(SELECTORS.wrapButton);
    await click(SELECTORS.dropdown);
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

  test('it wrapped data', async function (assert) {
    assert.expect(1);
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
    assert
      .dom(`${SELECTORS.masked} ${SELECTORS.copyButton}`)
      .hasAttribute('data-test-copy-button', this.wrappedData, 'it renders wrapped data');
  });
});
