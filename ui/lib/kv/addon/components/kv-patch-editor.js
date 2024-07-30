/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { assert } from '@ember/debug';

class StateManager {
  @tracked _state;
  possibleStates = ['disabled', 'readonly', 'deleted'];

  constructor(keys) {
    // initially disable all inputs
    this._state = Object.keys(keys).reduce((obj, key) => {
      obj[key] = 'disabled';
      return obj;
    }, {});
  }

  set(key, status) {
    assert(
      `state must be one of the following: ${this.possibleStates.join(
        ', '
      )}, you attempted to set "${status}"`,
      this.possibleStates.includes(status)
    );
    const newState = { ...this._state };
    newState[key] = status;
    this._state = newState;
  }

  get(key) {
    return this._state[key];
  }
}

export default class KvPatchEditor extends Component {
  @tracked state;
  @tracked patchData = {};

  getState = (key) => this.state.get(key);

  constructor() {
    super(...arguments);
    this.args.value;
    this.state = new StateManager(this.args.value);
  }

  @action
  setState(key, status) {
    this.state.set(key, status);
  }

  @action
  submit(event) {
    event.preventDefault();
    const formData = new FormData(event.target);
    const data = Object.fromEntries(formData.entries());
    this.args.onSubmit(data);
  }
}
