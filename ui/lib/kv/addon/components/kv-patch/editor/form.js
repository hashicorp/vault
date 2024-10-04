/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { A } from '@ember/array';
import {
  hasWhitespace,
  isNonString,
  WHITESPACE_WARNING,
  NON_STRING_WARNING,
} from 'vault/utils/model-helpers/validators';

/**
 * @module KvPatch::Editor::Form
 * @description
 * This component renders one of two ways to patch a KV v2 secret (the other is using the JSON editor).
 * Each top-level subkey returned by the API endpoint renders in a disabled column with an empty (also disabled) value input beside it.
 * Initially, an edit or delete button is left of the value input. Clicking "Delete" marks a key for deletion (it does not remove the row).
 * Clicking "Edit" enables the value input (the key input for retrieved subkeys is never editable). Users can then input a new value for that key.
 * If either button is clicked it is replaced by a "Cancel" button. Canceling empties the value input and returns it to a 'disabled' state
 *
 * Additionally, there is one empty row at the bottom for adding new key/value pairs.
 * Clicking "Add" adds the new key/value pair to the internally tracked state (an array) and creates a new empty row.
 * Newly added keys are editable and therefore never disabled.
 * A newly added pair can be undone by clicking "Remove" which deletes the row and removes it from the tracked array.
 *
 * Clicking the "Reveal subkeys in JSON" toggle displays the full, nested subkey structure returned by the API.
 *
 * @example
 * <KvPatch::Editor::Form @subkeys={{@subkeys}} @onSubmit={{perform this.save}} @onCancel={{this.onCancel}} @isSaving={{this.save.isRunning}} />
 *
 * @param {boolean} isSaving - if true, disables the save and cancel buttons. useful if the onSubmit callback is a concurrency task
 * @param {function} onCancel - called when form is canceled
 * @param {function} onSubmit - called when form is saved, called with with the key value object containing patch data
 * @param {object} subkeys - leaf keys of a kv v2 secret, all values (unless a nested object with more keys) return null. https://developer.hashicorp.com/vault/api-docs/secret/kv/kv-v2#read-secret-subkeys
 * @param {string} submitError - error message string from parent if submit failed
 */

export class KeyValueState {
  @tracked key;
  @tracked value;
  @tracked state; // 'enabled', 'disabled' or 'deleted'
  @tracked keyError;

  constructor({ key, value = undefined, state = 'disabled' }) {
    this.key = key;
    this.value = value;
    this.state = state;
  }

  get keyWarning() {
    return hasWhitespace(this.key) ? WHITESPACE_WARNING('this key') : '';
  }

  get valueWarning() {
    if (this.value === null) return '';
    return isNonString(this.value) ? NON_STRING_WARNING : '';
  }

  reset() {
    this.value = undefined;
    this.state = 'disabled';
  }

  @action
  updateValue(event) {
    this.value = event.target.value;
  }

  @action
  updateState(state) {
    this.state = state;
  }
}

export default class KvPatchEditorForm extends Component {
  @tracked patchData; // key value pairs in form
  @tracked showSubkeys = false;
  @tracked validationError;

  // tracked variables for new (initially empty) row of inputs.
  // once a user clicks "Add" a KeyValueState class is instantiated for that row
  @tracked newKey;
  @tracked newValue;

  isOriginalSubkey = (key) => Object.keys(this.args.subkeys).includes(key);

  constructor() {
    super(...arguments);
    const kvData = Object.keys(this.args.subkeys).map((key) => this.generateData(key));
    this.patchData = A(kvData);
    this.resetNewRow();
  }

  get newKeyWarning() {
    return hasWhitespace(this.newKey) ? WHITESPACE_WARNING('this key') : '';
  }

  get newValueWarning() {
    if (this.newValue === null) return '';
    return isNonString(this.newValue) ? NON_STRING_WARNING : '';
  }

  get newKeyError() {
    return this.validateKey(this.newKey);
  }

  generateData(key, value, state) {
    return new KeyValueState({ key, value, state });
  }

  resetNewRow() {
    this.newKey = undefined;
    this.newValue = undefined;
  }

  validateKey(key) {
    return this.patchData.any((KV) => KV.key === key)
      ? `"${key}" key already exists. Update the value of the existing key or rename this one.`
      : '';
  }

  @action
  updateKey(KV, event) {
    // KV is KeyValueState class
    const key = event.target.value;
    // if a user refocuses an input that already has a key
    // validateKey miscalculates and thinks it's a duplicate
    if (KV.key === key) return; // so we return if values match
    const isInvalid = this.validateKey(key);
    KV.keyError = isInvalid;
    if (isInvalid) return;
    // only set if valid, otherwise key matches original
    // subkey and input state updates to readonly
    KV.key = key;
  }

  @action
  updateNewKey(event) {
    const key = event.target.value;
    this.newKey = key;
  }

  @action
  updateNewValue(event) {
    this.newValue = event.target.value;
  }

  @action
  addRow() {
    if (!this.newKey || this.newKeyError) return;
    const KV = this.generateData(this.newKey, this.newValue, 'enabled');
    this.patchData.pushObject(KV);
    // reset tracked values after adding them to patchData
    this.resetNewRow();
  }

  @action
  undoKey(KV) {
    if (this.isOriginalSubkey(KV.key)) {
      // reset state to 'disabled' and value to undefined
      KV.reset();
    } else {
      // remove row all together
      this.patchData.removeObject(KV);
    }
  }

  @action
  submit(event) {
    event.preventDefault();
    if (this.newKeyError || this.patchData.any((KV) => KV.keyError)) {
      this.validationError = 'This form contains validations errors, please resolve those before submitting.';
      return;
    }

    // patchData will not include the last row if a user has not clicked "Add"
    // manually check for data and add it to this.patchData
    if (this.newKey) {
      this.addRow();
    }

    const data = this.patchData.reduce((obj, KV) => {
      // only include edited inputs
      const { state } = KV;
      if (state === 'enabled' || state === 'deleted') {
        const value = state === 'deleted' ? null : KV.value;
        obj[KV.key] = value;
      }
      return obj;
    }, {});

    this.args.onSubmit(data);
  }
}
