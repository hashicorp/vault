/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { A } from '@ember/array';
import { hasWhitespace, isNonString } from 'vault/utils/validators';

const WHITESPACE_WARNING =
  "This key contains whitespace. If this is desired, you'll need to encode it with %20 in API requests.";

const NON_STRING_WARNING =
  'This value will be saved as a string. If you need to save a non-string value, please use the JSON editor.';

class Kv {
  @tracked key;
  @tracked value;
  @tracked state;
  @tracked isInvalidKey;

  constructor({ key, value = undefined, state = 'disabled' }) {
    this.key = key;
    this.value = value;
    this.state = state;
  }

  get keyHasWarning() {
    return hasWhitespace(this.key) ? WHITESPACE_WARNING : '';
  }

  get valueHasWarning() {
    if (this.newValue === null) return '';
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

    if (state === 'deleted') {
      this.value = null;
    }
  }
}

export default class KvPatchEditor extends Component {
  @tracked patchData; // key value pairs in form
  @tracked showSubkeys = false;

  // tracked variables for new (initially empty) row of inputs
  // once a user clicks "Add" a Kv class is instantiated for that row
  // and it is added to the patchData array
  @tracked isInvalidKey = '';
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
    return hasWhitespace(this.newKey) ? WHITESPACE_WARNING : '';
  }

  get newValueWarning() {
    if (this.newValue === null) return '';
    return isNonString(this.newValue) ? NON_STRING_WARNING : '';
  }

  generateData(key, value, state) {
    return new Kv({ key, value, state });
  }

  resetNewRow() {
    this.newKey = undefined;
    this.newValue = undefined;
  }

  @action
  updateKey(KV, event) {
    const key = event.target.value;
    const isInvalid = this.validateKey(key);

    if (KV) {
      KV.isInvalidKey = isInvalid;
      if (isInvalid) return; // don't set if invalid
      KV.key = key;
    } else {
      this.isInvalidKey = isInvalid;
      if (isInvalid) return; // don't set if invalid
      this.newKey = key;
    }
  }

  @action
  updateValue(event) {
    this.newValue = event.target.value;
  }

  @action
  addRow() {
    if (!this.newKey || this.isInvalidKey) return;
    const KV = this.generateData(this.newKey, this.newValue, 'enabled');
    this.patchData.pushObject(KV);
    // reset tracked values after adding them to patchData
    this.resetNewRow();
  }

  @action
  undo(KV) {
    if (this.isOriginalSubkey(KV.key)) {
      // reset state to disabled and value to undefined
      KV.reset();
    } else {
      // remove row all together
      this.patchData.removeObject(KV);
    }
  }

  @action
  submit(event) {
    event.preventDefault();
    // patchData will not include the last row if a user has not clicked "add"
    // manually check for data and add it to this.patchData
    if (this.newKey && this.newValue) {
      this.addRow();
    }

    const data = this.patchData.reduce((obj, KV) => {
      // only included edited inputs
      const { state } = KV;
      if (state === 'enabled' || state === 'deleted') {
        obj[KV.key] = KV.value;
      }
      return obj;
    }, {});

    this.args.onSubmit(data);
  }

  validateKey(key) {
    return this.patchData.any((KV) => KV.key === key)
      ? `"${key}" key already exists. Update the value of the existing key or rename this one.`
      : '';
  }
}
