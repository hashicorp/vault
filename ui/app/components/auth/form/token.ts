/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AuthBase from './base';

import type { TokenLoginApiResponse } from 'vault/vault/auth/methods';

/**
 * @module Auth::Form::Token
 * see Auth::Base
 * */

export default class AuthFormToken extends AuthBase {
  loginFields = [{ name: 'token' }];

  async loginRequest(formData: { token: string }) {
    const { token } = formData;

    const { data } = (await this.api.auth.tokenLookUpSelf(
      this.api.buildHeaders({ token })
    )) as TokenLoginApiResponse;

    // normalize auth data so stored token data has the same keys regardless of auth type
    return this.normalizeAuthResponse(data, {
      authMountPath: '',
      token: data.id,
      ttl: data.ttl,
    });
  }
}
