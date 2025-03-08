/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { allSupportedAuthBackends, supportedAuthBackends } from 'vault/helpers/supported-auth-backends';

/**
 * @module Auth::Template
 *
 * @example
 *
 * @param {string} authMethodQueryParam - auth method type to login with, updated by selecting an auth method from the dropdown
 */

class AuthState {
  @tracked type = '';
  @tracked mount = '';
  // fields
  @tracked token = '';
  @tracked username = '';
  @tracked password = '';
  @tracked role = '';
  @tracked jwt = '';

  resetFields() {
    this.token = '';
    this.userame = '';
    this.password = '';
    this.role = '';
    this.jwt = '';
  }

  constructor(type) {
    this.type = type;
  }
}

export default class AuthBase extends Component {
  @service version;

  @tracked state;

  constructor() {
    super(...arguments);
    this.state = new AuthState(this.args.authType);
  }

  @action
  handleError() {
    // do something
  }

  @action
  handleInput(evt) {
    // For changing values in this backing class, not on form
    const { name, value } = evt.target;
    this[name] = value;

    if (this.args.onUpdate) {
      // Do parent side effects like update query params
      this.args.onUpdate(name, value);
    }
  }

  maybeMask = (field) => {
    if (field === 'token' || field === 'password') {
      return 'password';
    }
    return 'text';
  };
}
