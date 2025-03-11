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

  showFields = ['username', 'password'];
  type = 'userpass';

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
        backend: this.type,
        data,
        selectedAuth: this.type,
      });

      // responsible for redirect after auth data is persisted
      // if auth service authSuccess method is called here, then we'd do that before calling parent onSuccess
      this.args.onSuccess(authResponse, this.type);
    } catch (error) {
      this.onError(error);
    }
  }
}
