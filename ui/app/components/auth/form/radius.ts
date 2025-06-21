/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AuthBase from './base';

import type { UsernameLoginResponse } from 'vault/vault/auth/methods';

/**
 * @module Auth::Form::Radius
 * see Auth::Base
 */

export default class AuthFormRadius extends AuthBase {
  loginFields = [{ name: 'username' }, { name: 'password' }];

  async loginRequest(formData: { path: string; username: string; password: string }) {
    const { path, username, password } = formData;

    const { auth } = <UsernameLoginResponse>await this.api.auth.radiusLogin(path, {
      username,
      password,
    });

    return this.normalizeAuthResponse(auth, {
      displayName: auth.metadata?.username,
      path,
      tokenKey: 'clientToken',
      ttlKey: 'leaseDuration',
    });
  }
}
