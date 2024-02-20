/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, find, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import sinon from 'sinon';

module('Integration | Component | flash-toast', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.flash = {
      type: 'info',
      message: 'The bare minimum flash message',
    };
    this.closeSpy = sinon.spy();
  });

  test('it renders', async function (assert) {
    await render(hbs`<FlashToast @flash={{this.flash}} @close={{this.closeSpy}} />`);

    assert.dom('[data-test-flash-message-body]').hasText('The bare minimum flash message');
    assert.dom('[data-test-flash-toast]').hasClass('hds-alert--color-highlight');
    await click('button');
    assert.ok(this.closeSpy.calledOnce, 'close action was called');
  });

  [
    { type: 'info', title: 'Info', color: 'hds-alert--color-highlight' },
    { type: 'success', title: 'Success', color: 'hds-alert--color-success' },
    { type: 'warning', title: 'Warning', color: 'hds-alert--color-warning' },
    { type: 'danger', title: 'Error', color: 'hds-alert--color-critical' },
    { type: 'foobar', title: 'Foobar', color: 'hds-alert--color-neutral' },
  ].forEach(({ type, title, color }) => {
    test(`it has correct title and color for type: ${type}`, async function (assert) {
      this.flash.type = type;
      await render(hbs`<FlashToast @flash={{this.flash}} @close={{this.closeSpy}} />`);

      assert.dom('[data-test-flash-toast-title]').hasText(title, 'title is correct');
      assert.dom('[data-test-flash-toast]').hasClass(color, 'color is correct');
    });
  });

  test('it renders messages with whitespaces correctly', async function (assert) {
    this.flash.message = `multi-

line msg`;

    await render(hbs`<FlashToast @flash={{this.flash}} @close={{this.closeSpy}} />`);
    const dom = find('[data-test-flash-message-body]');
    const lineHeight = 20;
    assert.true(dom.clientHeight > lineHeight, 'renders message on multiple lines');
  });
});
