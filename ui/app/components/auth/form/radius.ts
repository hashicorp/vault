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

    const { auth } = (await this.api.auth.radiusLoginWithUsername(username, path, {
      password,
    })) as UsernameLoginResponse;

    return this.normalizeAuthResponse(auth, {
      authMountPath: path,
      displayName: auth?.metadata?.username,
      token: auth.client_token,
      ttl: auth.lease_duration,
    });
  }
}
