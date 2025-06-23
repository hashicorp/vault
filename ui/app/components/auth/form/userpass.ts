/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AuthBase from './base';

import type { UsernameLoginResponse } from 'vault/vault/auth/methods';

/**
 * @module Auth::Form::Userpass
 * see Auth::Base
 * */

export default class AuthFormUserpass extends AuthBase {
  loginFields = [{ name: 'username' }, { name: 'password' }];

  async loginRequest(formData: { path: string; username: string; password: string }) {
    const { path, username, password } = formData;

    const { auth } = <UsernameLoginResponse>await this.api.auth.userpassLogin(username, path, {
      password,
    });

    // normalize auth data so stored token data has the same keys regardless of auth type
    return this.normalizeAuthResponse(auth, {
      displayName: auth.metadata?.username,
      path,
      tokenKey: 'clientToken',
      ttlKey: 'leaseDuration',
    });
  }
}
