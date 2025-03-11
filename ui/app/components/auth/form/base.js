/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';

/**
 * @module Auth::Base
 *
 * @example
 *
 * @param {string} param -
 */

export default class AuthBase extends Component {
  @service auth;

  maybeMask = (field) => {
    if (field === 'token' || field === 'password') {
      return 'password';
    }
    return 'text';
  };

  @action
  async login(event) {
    event.preventDefault();
    // base login flow
  }

  onError(error) {
    if (!this.auth.mfaError) {
      const errorMessage = `Authentication failed: ${this.auth.handleError(error)}`;
      this.args.onError(errorMessage);
    }
  }
}
