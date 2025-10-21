/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { create } from 'jsondiffpatch';

// this service tracks the unsaved changes modal state.
export default class UnsavedChangesService extends Service {
  @tracked changedFields: Array<string> = [];
  @tracked showModal = false;

  @tracked initialState: Record<string, unknown> | undefined;
  @tracked currentState: Record<string, unknown> | undefined;

  setupProperties(
    initialState: Record<string, unknown> | undefined,
    currentState: Record<string, unknown> | undefined
  ) {
    this.initialState = initialState;
    this.currentState = currentState;
  }

  getDiff() {
    const diffpatcher = create({});
    const delta = diffpatcher.diff(this.initialState, this.currentState);

    const changedFields = delta ? Object.keys(delta) : [];

    this.changedFields = changedFields;

    return changedFields;
  }

  get hasChanges() {
    return this.changedFields.length > 0;
  }
}
