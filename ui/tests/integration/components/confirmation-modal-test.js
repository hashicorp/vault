/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import sinon from 'sinon';
import { click, fillIn, render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | confirmation-modal', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders with disabled confirmation button until input matches', async function (assert) {
    const confirmAction = sinon.spy();
    const closeAction = sinon.spy();
    this.set('onConfirm', confirmAction);
    this.set('onClose', closeAction);
    await render(hbs`
            <ConfirmationModal
        @title="Confirmation Modal"
        @isActive={{true}}
        @onConfirm={{this.onConfirm}}
        @onClose={{this.onClose}}
        @buttonText="Plz Continue"
        @confirmText="Destructive Thing"
      />
    `);

    assert.dom('[data-test-confirm-button]').isDisabled();
    assert.dom('#confirmation-modal').exists('modal is active');
    assert.dom('[data-test-confirm-button]').hasText('Plz Continue', 'Confirm button has specified value');
    assert
      .dom('[data-test-confirmation-modal-title] [data-test-icon="alert-triangle"]')
      .exists('title has with warning icon');
    await fillIn('[data-test-confirmation-modal-input="Confirmation Modal"]', 'Destructive Thing');
    assert.dom('[data-test-confirm-button="Confirmation Modal"]').isNotDisabled();

    await click('[data-test-cancel-button]');
    assert.true(closeAction.called, 'executes passed in onClose function');
    await click('[data-test-confirm-button]');
    assert.true(confirmAction.called, 'executes passed in onConfirm function');
  });
});
