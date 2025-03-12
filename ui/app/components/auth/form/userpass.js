/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AuthBase from './base';
import { action } from '@ember/object';
import { service } from '@ember/service';

/**
 * @module Auth::Form::Userpass
 *
 * */

export default class AuthFormUserpass extends AuthBase {
  @service auth;

  loginFields = ['username', 'password'];

  @action
  async login(event) {
    event.preventDefault();
    const formData = new FormData(event.target);
    const data = {};

    for (const key of formData.keys()) {
      data[key] = formData.get(key);
    }

    try {
      const authResponse = await this.auth.authenticate({
        clusterId: this.args.cluster.id,
        backend: this.args.authType,
        data,
        selectedAuth: this.args.authType,
      });

      // responsible for redirect after auth data is persisted
      // if auth service authSuccess method is called here, then we'd do that before calling parent onSuccess
      this.onSuccess(authResponse);
    } catch (error) {
      this.onError(error);
    }
  }
}
