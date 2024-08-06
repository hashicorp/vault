/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

// TODO validations
/*
  - show warning for matching keys, disable add
  - show warning that numbers will be submitted as strings
  - show warning for whitespace 
  */
class InputStateManager {
  @tracked _state;
  possibleStates = ['enabled', 'disabled', 'deleted'];

  constructor(emptySubkeys) {
    // initially set all inputs to disabled
    this._state = Object.keys(emptySubkeys).reduce((obj, key) => {
      obj[key] = 'disabled';
      return obj;
    }, {});
  }

  set(key, status) {
    const newState = { ...this._state };
    newState[key] = status;
    // trigger update
    this._state = newState;
  }

  get(key) {
    return this._state[key];
  }
}

class KvData {
  @tracked _kvData;

  constructor(kvData) {
    this._kvData = { ...kvData };
  }

  set(key, value) {
    const newObject = this._kvData;
    newObject[key] = value;
    // trigger update
    this._kvData = newObject;
  }

  get(key) {
    return this._kvData[key];
  }

  deleteKey(key) {
    const newObject = { ...this._kvData };
    delete newObject[key];
    // trigger an update
    this._kvData = newObject;
  }
}

export default class KvPatchEditor extends Component {
  @tracked state; // disabled, enabled or deleted input states
  @tracked patchData; // key value pairs in form

  // tracked variables for new row of inputs after user clicks "add"
  @tracked newKey = '';
  @tracked newValue = '';

  getState = (key) => this.state.get(key);
  getValue = (key) => this.patchData.get(key);

  isOriginalKey = (key) => Object.keys(this.args.emptySubkeys).includes(key);

  constructor() {
    super(...arguments);
    this.state = new InputStateManager(this.args.emptySubkeys);
    this.patchData = new KvData(this.args.emptySubkeys);
  }

  get formData() {
    return this.patchData._kvData;
  }

  @action
  setState(key, status) {
    this.state.set(key, status);

    if (status === 'deleted') {
      this.patchData.set(key, null);
    }
  }

  @action
  handleCancel(key) {
    if (this.isOriginalKey(key)) {
      this.setState(key, 'disabled');
      // reset value to empty string
      this.patchData.set(key, '');
    } else {
      // remove row all together
      this.patchData.deleteKey(key);
    }
  }

  @action
  onBlurNew(event) {
    const { name, value } = event.target;
    this[name] = value;
  }

  @action
  onBlurExisting(key, type, event) {
    if (type === 'key') {
      // store value of original key
      const value = this.patchData.get(key);
      const newKey = event.target.value;
      // delete old key
      this.patchData.deleteKey(key);
      // add new one
      this.patchData.set(newKey, value);
    }
    if (type === 'value') {
      const value = event.target.value;
      this.patchData.set(key, value);
    }
  }

  @action
  handleAddClick() {
    this.patchData.set(this.newKey, this.newValue);
    this.setState(this.newKey, 'enabled');
    // reset tracked values after adding them to patchData
    this.newKey = '';
    this.newValue = '';
  }

  @action
  submit(event) {
    event.preventDefault();
    // patchData will not include the last row if a user has not clicked "add"
    // manually check for for data and add it the payload
    if (this.newKey && this.newValue) {
      this.handleAddClick();
    }

    const data = this.formData;
    // remove any empty strings from the payload
    for (const [key, value] of Object.entries(data)) {
      if (value === '') {
        delete data[key];
      }
    }
    this.args.onSubmit(data);
  }
}
