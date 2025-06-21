/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AuthBase from './base';

import type { GithubLoginApiResponse } from 'vault/vault/auth/methods';

/**
 * @module Auth::Form::Github
 * see Auth::Base
 */

export default class AuthFormGithub extends AuthBase {
  loginFields = [{ name: 'token', label: 'Github token' }];

  async loginRequest(formData: { path: string; token: string }) {
    const { path, token } = formData;

    const { auth } = <GithubLoginApiResponse>await this.api.auth.githubLogin(path, {
      token,
    });

    const { metadata } = auth;
    // normalize auth data so stored token data has the same keys regardless of auth type

    return this.normalizeAuthResponse(auth, {
      displayName: `${metadata?.org}/${metadata.username}`,
      path,
      tokenKey: 'clientToken',
      ttlKey: 'leaseDuration',
    });
  }
}
