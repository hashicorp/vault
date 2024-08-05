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

  constructor(kvObject) {
    // initially disable all inputs
    this._state = Object.keys(kvObject).reduce((obj, key) => {
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

  constructor(kvObject) {
    this._kvData = kvObject;
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
  @tracked formData;
  @tracked newKey = '';
  @tracked newValue = '';

  getState = (key) => this.state.get(key);
  isExistingKey = (key) => Object.keys(this.args.kvObject).includes(key);

  constructor() {
    super(...arguments);
    this.state = new StateManager(this.args.kvObject);
    this.formData = new KvData(this.args.kvObject);
  }

  get inputData() {
    return this.formData._kvData;
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
      this.formData.deleteKey(key);
    }
  }

  @action
  onBlur(event) {
    const { name, value } = event.target;
    this[name] = value;
  }

  @action
  addData(key, value) {
    this.formData.set(key, value);
    this.newKey = '';
    this.newValue = '';
    this.state.set(key, 'enabled');
  }

  @action
  submit(event) {
    event.preventDefault();
    const formData = new FormData(event.target);
    const data = Object.fromEntries(formData.entries());
    this.args.onSubmit(data);
  }
}
