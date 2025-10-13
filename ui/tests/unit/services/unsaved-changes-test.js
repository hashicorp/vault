/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Service | unsaved-changes', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    this.unsavedChanges = this.owner.lookup('service:unsaved-changes');
  });

  test('it should set properties', async function (assert) {
    const initialState = { title: 'Title 1', description: 'Old description' };
    const currentState = { title: 'New Title 1', description: 'New description' };
    this.unsavedChanges.setupProperties(initialState, currentState);
    assert.deepEqual(this.unsavedChanges.initialState, initialState);
    assert.deepEqual(this.unsavedChanges.currentState, currentState);
  });

  test('it should update changedFields when getDiff is called', async function (assert) {
    const initialState = { title: 'Title 1', description: 'Old description' };
    const currentState = { title: 'Title 1', description: 'New description' };
    this.unsavedChanges.setupProperties(initialState, currentState);
    assert.deepEqual(this.unsavedChanges.changedFields, []);
    this.unsavedChanges.getDiff();
    assert.deepEqual(this.unsavedChanges.changedFields, ['description']);
  });
});
