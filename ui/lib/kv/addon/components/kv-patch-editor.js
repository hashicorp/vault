/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

class StateManager {
  @tracked _state;
  possibleStates = ['enabled', 'disabled', 'deleted'];

  constructor(kvData) {
    // initially disable all inputs
    this._state = Object.keys(kvData).reduce((obj, key) => {
      obj[key] = 'disabled';
      return obj;
    }, {});
  }

  set(key, status) {
    const newState = { ...this._state };
    newState[key] = status;
    this._state = newState;
  }

  get(key) {
    return this._state[key];
  }
}

class KvData {
  @tracked _kvData;

  constructor(kvData) {
    this._kvData = kvData;
  }

  set(key, value) {
    const newObject = { ...this._kvData, [key]: value };
    // trigger an update
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
  @tracked state;
  @tracked patchData;
  @tracked newKey = '';
  @tracked newValue = '';

  getState = (key) => this.state.get(key);
  getValue = (key) => this.patchData.get(key);

  isOriginalKey = (key) => Object.keys(this.args.kvData).includes(key);

  constructor() {
    super(...arguments);
    this.state = new StateManager(this.args.kvData);
    this.patchData = new KvData(this.args.kvData);
  }

  get inputData() {
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
      this.patchData.set(key, '');
    } else {
      this.patchData.deleteKey(key);
    }
  }

  @action
  onBlurNew(event) {
    const { name, value } = event.target;
    this[name] = value;
  }

  @action
  onBlur(key, type, event) {
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
    // check the last row for data, add it the payload if there are values
    if (this.newKey && this.newValue) {
      this.patchData.set(this.newKey, this.newValue);
    }

    const data = this.patchData._kvData;
    for (const [key, value] of Object.entries(this.patchData._kvData)) {
      if (value === '') {
        delete data[key];
      }
    }
    // if a user doesn't click add make sure we include the final row of key/values
    // if there are matching keys, show the validation warning and disable the add
    // and prevent it overriding the existing key/value onBlur
    // console.log(data, 'tracked patch data!!');
    this.args.onSubmit(this.patchData);
  }
}
