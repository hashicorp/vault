/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | UnsavedChangesModal', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.unsavedChanges = this.owner.lookup('service:unsavedChanges');
    this.unsavedChanges.showModal = true;
    const initialState = { title: 'Title 1', description: 'Old description' };
    const currentState = { title: 'Title 1', description: 'New description' };
    this.unsavedChanges.initialState = initialState;
    this.unsavedChanges.currentState = currentState;
    this.save = () => 'saved!';
    this.discard = () => 'discarded!';
  });

  test('it shows unsaved changes modal', async function (assert) {
    await render(hbs`<UnsavedChangesModal @onSave={{this.save}} @onDiscard={{this.discard}} />`);
    assert.dom(GENERAL.modal.header('unsaved-changes')).hasText('Unsaved changes');
    assert
      .dom(GENERAL.modal.body('unsaved-changes'))
      .hasText(`You've made changes to the following: Description Would you like to apply them?`);
  });

  test('it shows unsaved changes modal with custom changedFields', async function (assert) {
    await render(
      hbs`<UnsavedChangesModal @onSave={{this.save}} @onDiscard={{this.discard}} @changedFields={{this.changedFields}}/>`
    );
    assert.dom(GENERAL.modal.header('unsaved-changes')).hasText('Unsaved changes');
    assert
      .dom(GENERAL.modal.body('unsaved-changes'))
      .hasText(`You've made changes to the following: Description Would you like to apply them?`);
  });
});
