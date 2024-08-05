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
  isExistingKey = (key) => Object.keys(this.args.kvData).includes(key);

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
  }

  @action
  handleCancel(key) {
    if (this.isExistingKey(key)) {
      this.state.set(key, 'disabled');
    } else {
      this.patchData.deleteKey(key);
    }
  }

  @action
  onBlur(event) {
    const { name, value } = event.target;
    this[name] = value;
  }

  @action
  addData(key, value) {
    this.patchData.set(key, value);
    this.newKey = '';
    this.newValue = '';
    this.state.set(key, 'enabled');
  }

  @action
  submit(event) {
    // if a user doesn't click add make sure we include the final row of key/values
    // if there are matching keys, show the validation warning and disable the add
    // and prevent it overriding the existing key/value onBlur
    event.preventDefault();
    const patchData = new FormData(event.target);
    const data = Object.fromEntries(patchData.entries());
    this.args.onSubmit(data);
  }
}
