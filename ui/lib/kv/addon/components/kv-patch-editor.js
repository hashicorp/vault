/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { A } from '@ember/array';

// TODO validations
/*
  - show warning for matching keys, disable add
  - show warning that numbers will be submitted as strings
  - show warning for whitespace 
  */

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

  get isPatchable() {
    // only include edited inputs in payload
    return this.state === 'enabled' || this.state === 'deleted';
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

  // tracked variables for new row of inputs after user clicks "add"
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

  get isAddAllowed() {
    return !this.newKey || this.isInvalidKey ? false : true;
  }

  generateData(key, value, state) {
    return new Kv({ key, value, state });
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
    if (!this.isAddAllowed) return;
    const data = this.generateData(this.newKey, this.newValue, 'enabled');
    this.patchData.pushObject(data);
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
    // manually check for for data and add it the payload
    if (this.newKey && this.newValue) {
      this.addRow();
    }

    // collect only relevant inputs
    const data = this.patchData.reduce((obj, KV) => {
      if (KV.isPatchable) {
        obj[KV.key] = KV.value;
      }
      return obj;
    }, {});

    this.args.onSubmit(data);
  }
}
