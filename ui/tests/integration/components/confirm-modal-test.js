/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import sinon from 'sinon';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | confirm-modal', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.onConfirm = sinon.stub();
    this.onClose = sinon.stub();
  });

  test('it renders a reasonable default', async function (assert) {
    await render(hbs`<ConfirmModal @onConfirm={{this.onConfirm}} @onClose={{this.onClose}} />`);

    assert
      .dom('[data-test-confirm-modal]')
      .hasClass('hds-modal--color-warning', 'renders warning modal color');
    assert.dom(GENERAL.confirmButton).hasClass('hds-button--color-primary', 'renders primary confirm button');
    assert.dom('[data-test-confirm-action-title]').hasText('Are you sure?', 'renders default title');
    assert
      .dom('[data-test-confirm-action-message]')
      .hasText('You will not be able to recover it later.', 'renders default body text');
  });

  test('it renders a custom title and message', async function (assert) {
    await render(
      hbs`<ConfirmModal @onConfirm={{this.onConfirm}} @onClose={{this.onClose}} @confirmTitle="Fancy Title" @confirmMessage="Riveting message" />`
    );

    assert.dom('[data-test-confirm-action-title]').hasText('Fancy Title', 'renders custom title');
    assert.dom('[data-test-confirm-action-message]').hasText('Riveting message', 'renders custom body text');
  });

  test('it renders a disabled message', async function (assert) {
    await render(
      hbs`<ConfirmModal @onConfirm={{this.onConfirm}} @onClose={{this.onClose}} @disabledMessage="Nope" />`
    );

    assert.dom('[data-test-confirm-action-title]').hasText('Not allowed', 'renders disabled title');

    assert.dom('[data-test-confirm-action-message]').hasText('Nope', 'renders disabled message');

    assert.dom(GENERAL.confirmButton).doesNotExist('confirm button is not rendered when action is disabled');
  });

  test('it calls onConfirm when the modal is confirmed', async function (assert) {
    await render(hbs`<ConfirmModal @onConfirm={{this.onConfirm}} @onClose={{this.onClose}} />`);

    await click(GENERAL.confirmButton);
    assert.true(this.onConfirm.called);
  });

  test('it renders a type-to-confirm input when @confirmText is passed', async function (assert) {
    await render(
      hbs`<ConfirmModal @onConfirm={{this.onConfirm}} @onClose={{this.onClose}} @confirmText="delete-me" />`
    );

    assert.dom(GENERAL.confirmTextInput).exists('type-to-confirm input is rendered');
    assert.dom(GENERAL.confirmWarning).doesNotExist('warning is not shown before interact');
  });

  test('it renders a custom @confirmLabel', async function (assert) {
    await render(
      hbs`<ConfirmModal @onConfirm={{this.onConfirm}} @onClose={{this.onClose}} @confirmText="delete-me" @confirmLabel="Confirm deletion" />`
    );

    assert
      .dom('label[id^="label-confirm-modal-input"]')
      .containsText('Confirm deletion', 'renders custom confirm label');
  });

  test('it shows a warning after clicking confirm with wrong input', async function (assert) {
    await render(
      hbs`<ConfirmModal @onConfirm={{this.onConfirm}} @onClose={{this.onClose}} @confirmText="delete-me" />`
    );

    await fillIn(GENERAL.confirmTextInput, 'wrong');
    await click(GENERAL.confirmButton);

    assert.dom(GENERAL.confirmTextInput).exists('type-to-confirm input is rendered');
    assert.dom(GENERAL.confirmWarning).exists('warning shown after failed confirm');
    assert.false(this.onConfirm.called, 'onConfirm is not called with wrong input');
  });

  test('it shows a warning after clicking confirm with empty input', async function (assert) {
    await render(
      hbs`<ConfirmModal @onConfirm={{this.onConfirm}} @onClose={{this.onClose}} @confirmText="delete-me" />`
    );

    await click(GENERAL.confirmButton);

    assert.dom(GENERAL.confirmWarning).exists('warning shown when input is empty');
    assert.false(this.onConfirm.called, 'onConfirm is not called with empty input');
  });

  test('warning clears when user retypes after a failed confirm', async function (assert) {
    await render(
      hbs`<ConfirmModal @onConfirm={{this.onConfirm}} @onClose={{this.onClose}} @confirmText="delete-me" />`
    );

    await click(GENERAL.confirmButton);
    assert.dom(GENERAL.confirmWarning).exists('warning shown after empty confirm');

    await fillIn(GENERAL.confirmTextInput, 'del');
    assert.dom(GENERAL.confirmWarning).doesNotExist('warning clears on retype');
  });

  test('it calls onConfirm when correct text is entered', async function (assert) {
    await render(
      hbs`<ConfirmModal @onConfirm={{this.onConfirm}} @onClose={{this.onClose}} @confirmText="delete-me" />`
    );

    await fillIn(GENERAL.confirmTextInput, 'delete-me');
    await click(GENERAL.confirmButton);

    assert.dom(GENERAL.confirmWarning).doesNotExist('no warning when input matches');
    assert.true(this.onConfirm.called, 'onConfirm is called when input matches');
  });
});
