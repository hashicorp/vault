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

    const { data } = <TokenLoginApiResponse>(
      await this.api.auth.tokenLookUpSelf(this.api.buildHeaders({ token }))
    );
    // normalize auth data so stored token data has the same keys regardless of auth type
    return this.normalizeAuthResponse(data, {
      displayName: data.displayName,
      path: '',
      tokenKey: 'id',
      ttlKey: 'ttl',
    });
  }
}
