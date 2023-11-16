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

  test('it closes the confirmation modal on successful delete', async function (assert) {
    const confirmAction = sinon.spy();
    this.set('onConfirm', confirmAction);
    await render(hbs`
      <ConfirmAction
        @buttonText="DELETE"
        @onConfirmAction={{this.onConfirm}}
      />
      `);
    await click('[data-test-confirm-action-trigger]');
    await click('[data-test-confirm-cancel-button]');
    // assert that after CANCEL the icon button is pointing down.
    assert.dom('[data-test-icon="chevron-down"]').exists('Icon is pointing down after clicking cancel');
    // open the modal again to test the DELETE action
    await click('[data-test-confirm-action-trigger="true"]');
    await click('[data-test-confirm-button="true"]');
    assert
      .dom('[data-test-icon="chevron-down"]')
      .exists('Icon is pointing down after executing the Delete action');
    assert.true(confirmAction.called, 'calls the action when Delete is pressed');
    assert
      .dom('[data-test-confirm-action-title]')
      .doesNotExist('it has closed the confirm content and does not show the title');
  });
});
