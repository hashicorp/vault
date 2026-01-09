/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, find, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import sinon from 'sinon';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

const SELECTORS = {
  element: '[data-test-flash-toast ]',
  title: '[data-test-flash-toast-title]',
  message: '[data-test-flash-message-body]',
};

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

    assert.dom(SELECTORS.message).hasText('The bare minimum flash message');
    assert.dom(SELECTORS.element).hasClass('hds-alert--color-highlight');
    assert.dom('a').doesNotExist('link does not render');
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

      assert.dom(SELECTORS.title).hasText(title, 'title is correct');
      assert.dom(SELECTORS.element).hasClass(color, 'color is correct');
    });
  });

  test('it renders messages with whitespaces correctly', async function (assert) {
    this.flash.message = `multi-

line msg`;

    await render(hbs`<FlashToast @flash={{this.flash}} @close={{this.closeSpy}} />`);
    const dom = find(SELECTORS.message);
    const lineHeight = 20;
    assert.true(dom.clientHeight > lineHeight, 'renders message on multiple lines');
  });

  test('it renders custom title when provided', async function (assert) {
    this.flash.title = 'Tada a flash!';
    await render(hbs`<FlashToast @flash={{this.flash}} @close={{this.closeSpy}} />`);
    assert.dom(SELECTORS.title).hasText('Tada a flash!');
  });

  test('it renders link when provided', async function (assert) {
    this.flash.link = {
      text: 'A snazzy link',
      route: 'vault.cluster.policy.show',
    };

    await render(hbs`<FlashToast @flash={{this.flash}} @close={{this.closeSpy}} />`);
    assert
      .dom('a')
      .exists()
      .hasText('A snazzy link')
      .hasClass('hds-link-standalone--color-secondary', 'it renders secondary color')
      .hasClass('hds-link-standalone--icon-position-trailing', 'it renders trailing icon');
    assert.dom(GENERAL.icon('arrow-right')).exists('it renders default icon');
  });

  test('it renders custom icon for link', async function (assert) {
    this.flash.link = { text: 'Click me', href: 'example.com', icon: 'rocket' };
    await render(hbs`<FlashToast @flash={{this.flash}} @close={{this.closeSpy}} />`);
    assert.dom(GENERAL.icon('rocket')).exists();
  });
});
