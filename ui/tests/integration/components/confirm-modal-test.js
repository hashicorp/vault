/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import sinon from 'sinon';

module('Integration | Component | confirm-modal', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.onConfirm = sinon.spy();
    this.onClose = sinon.spy();
  });

  test('it renders a reasonable default', async function (assert) {
    await render(hbs`<ConfirmModal @onConfirm={{this.onConfirm}} @onClose={{this.onClose}} />`);
    assert
      .dom('[data-test-confirm-modal]')
      .hasClass('hds-modal--color-warning', 'renders warning modal color');
    assert
      .dom('[data-test-confirm-button]')
      .hasClass('hds-button--color-primary', 'renders primary confirm button');
    assert.dom('[data-test-confirm-action-title]').hasText('Are you sure?', 'renders default title');
    assert
      .dom('[data-test-confirm-action-message]')
      .hasText('You will not be able to recover it later.', 'renders default body text');
    await click('[data-test-confirm-cancel-button]');
    assert.ok(this.onClose.called, 'calls the onClose action when Cancel is clicked');
    await click('[data-test-confirm-button]');
    assert.ok(this.onConfirm.called, 'calls the onConfirm action when Confirm is clicked');
  });
});
