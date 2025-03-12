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
 * @param {string} authType - chosen login method type
 * @param {object} cluster - The cluster model which contains information such as cluster id, name and boolean for if the cluster is in standby
 * @param {function} onError - callback if there is a login error
 * @param {function} onSuccess - calls onAuthResponse in auth/page redirects if successful
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

  // if we move auth service authSuccess method here (or to each auth method component)
  // then call that before calling parent this.args.onSuccess
  onSuccess(authResponse) {
    //  responsible for redirect after auth data is persisted
    this.args.onSuccess(authResponse, this.args.authType);
  }

  onError(error) {
    if (!this.auth.mfaError) {
      const errorMessage = `Authentication failed: ${this.auth.handleError(error)}`;
      this.args.onError(errorMessage);
    }
  }
}
