/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AuthBase from './base';
import { service } from '@ember/service';
import { action } from '@ember/object';

/**
 * @module Auth::Form::Token
 *
 * */

export default class AuthFormToken extends AuthBase {
  @service auth;

  showFields = ['token'];
  type = 'token';

  // extrapolated data from auth service and SUPPORTED_AUTH_BACKENDS
  // depending on ember data affects on auth service, use this data here instead
  url = '/v1/auth/token/lookup-self';
  displayNamePath = 'display_name';
  tokenPath = 'id';

  @action
  async login(event) {
    event.preventDefault();
    const data = {};
    this.showFields.forEach((f) => {
      data[f] = this.state[f];
    });
    const authResponse = await this.auth.authenticate({
      clusterId: this.args.cluster.id,
      backend: this.type,
      data,
      selectedAuth: this.type,
    });

    // responsible for redirect after auth data is persisted
    // if auth service authSuccess method is called here, then we'd do that before calling parent onSuccess
    this.args.onSuccess(authResponse, this.type);
  }
}
