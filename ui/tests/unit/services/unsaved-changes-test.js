/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Service | unsaved-changes', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    this.unsavedChanges = this.owner.lookup('service:unsaved-changes');
  });

  test('it should get changedFields', async function (assert) {
    const initialState = { title: 'Title 1', description: 'Old description' };
    const currentState = { title: 'Title 1', description: 'New description' };
    this.unsavedChanges.initialState = initialState;
    this.unsavedChanges.currentState = currentState;

    assert.deepEqual(this.unsavedChanges.changedFields, ['description']);
  });
  test('it should get hasChanges', async function (assert) {
    const initialState = { title: 'Title 1', description: 'Old description' };
    const currentState = { title: 'Title 1', description: 'New description' };
    this.unsavedChanges.initialState = initialState;
    this.unsavedChanges.currentState = currentState;

    assert.true(this.unsavedChanges.hasChanges, 'shows that there are unsaved changes');
    currentState.description = 'Old description';
    this.unsavedChanges.initialState = initialState;
    this.unsavedChanges.currentState = currentState;
    assert.false(this.unsavedChanges.hasChanges, 'shows that there are no unsaved changes');
  });
});
