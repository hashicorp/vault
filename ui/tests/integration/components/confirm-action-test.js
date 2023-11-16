/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';

module('Integration | Component | confirm-action', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function (assert) {
    const confirmAction = sinon.spy();
    this.set('onConfirm', confirmAction);
    await render(hbs`
      <ConfirmAction
        @buttonText="DELETE"
        @onConfirmAction={{this.onConfirm}}
      />
      `);

    assert.dom('[data-test-confirm-action-trigger]').hasText('DELETE', 'renders button text');
    await click('[data-test-confirm-action-trigger]');
    assert.dom('[data-test-confirm-action-title]').hasText('Are you sure?', 'renders default title');
    assert
      .dom('[data-test-confirm-action-message]')
      .hasText('You will not be able to recover it later.', 'renders default body text');
    await click('[data-test-confirm-cancel-button]');
    assert.false(confirmAction.called, 'does not call the action when Cancel is clicked');
    await click('[data-test-confirm-action-trigger]');
    await click('[data-test-confirm-button]');
    assert.true(confirmAction.called, 'calls the action when Confirm is clicked');
    assert.dom('[data-test-confirm-action-title]').doesNotExist('modal closes after confirm is clicked');
  });

  test('it renders loading state', async function (assert) {
    const confirmAction = sinon.spy();
    this.set('onConfirm', confirmAction);
    await render(hbs`
      <ConfirmAction
        @buttonText="Open!"
        @onConfirmAction={{this.onConfirm}}
        @isRunning={{true}}
      />
      `);

    await click('[data-test-confirm-action-trigger]');

    assert.dom('[data-test-confirm-button]').isDisabled('disables confirm button when loading');
    assert.dom('[data-test-confirm-button] [data-test-icon="loading"]').exists('it renders loading icon');
  });

  // handle passed in args confirmTitle and modal color logic with button
});
