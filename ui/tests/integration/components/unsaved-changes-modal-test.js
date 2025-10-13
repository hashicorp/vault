/**
 * Copyright (c) HashiCorp, Inc.
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
    this.unsavedChanges.changedFields = ['Description', 'Secrets duration'];
    this.save = () => 'saved!';
    this.discard = () => 'discarded!';
  });

  test('it shows unsaved changes modal', async function (assert) {
    await render(hbs`<UnsavedChangesModal @onSave={{this.save}} @onDiscard={{this.discard}} />`);
    assert.dom(GENERAL.modal.header('unsaved-changes')).hasText('Unsaved changes');
    assert
      .dom(GENERAL.modal.body('unsaved-changes'))
      .hasText(
        `You've made changes to the following: Description Secrets duration Would you like to apply them?`
      );
  });

  test('it shows unsaved changes modal with custom changedFields', async function (assert) {
    this.changedFields = ['Field 1', 'Field 2'];
    await render(
      hbs`<UnsavedChangesModal @onSave={{this.save}} @onDiscard={{this.discard}} @changedFields={{this.changedFields}}/>`
    );
    assert.dom(GENERAL.modal.header('unsaved-changes')).hasText('Unsaved changes');
    assert
      .dom(GENERAL.modal.body('unsaved-changes'))
      .hasText(`You've made changes to the following: Field 1 Field 2 Would you like to apply them?`);
  });
});
