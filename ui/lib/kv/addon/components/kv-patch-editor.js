/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { A } from '@ember/array';
import { containsWhiteSpace } from 'vault/utils/validators';

// TODO validations
/*
  - show warning for matching keys, disable add
  - show warning that numbers will be submitted as strings
  - show warning for whitespace 
  */

const WHITESPACE_WARNING =
  "This key contains whitespace. If this is desired, you'll need to encode it with %20 in APi requests.";

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
    const isValid = containsWhiteSpace(this.key); // returns false (invalid) if contains whitespace
    return isValid ? '' : WHITESPACE_WARNING;
  }

  reset() {
    this.value = undefined;
    this.state = 'disabled';
  }

  @action
  onBlur(input, evt) {
    this[input] = evt.target.value;
  }

  @action
  onClick(state) {
    this.state = state;

    if (state === 'deleted') {
      this.value = null;
    }
  }
}

export default class KvPatchEditor extends Component {
  @tracked patchData; // key value pairs in form

  // tracked variables for new (initially empty) row of inputs
  // once a user clicks "Add" a Kv class is instantiated for that row
  // and it is added to the patchData array
  @tracked isInvalidKey = '';
  @tracked newKey;
  @tracked newValue;

  isOriginalSubkey = (key) => this.args.subkeyArray.includes(key);

  constructor() {
    super(...arguments);
    const kvData = this.args.subkeyArray.map((key) => this.generateData(key));
    this.patchData = A(kvData);
    this.resetNewRow();
  }

  get keyHasWarning() {
    const isValid = containsWhiteSpace(this.newKey); // returns false (invalid) if contains whitespace
    return isValid ? '' : WHITESPACE_WARNING;
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
  handleNewRow(event) {
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
