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

    const { auth } = (await this.api.auth.githubLogin(path, {
      token,
    })) as GithubLoginApiResponse;

    const { org, username } = auth?.metadata || {};
    const displayName = org && username ? `${org}/${username}` : username || org || '';

    return this.normalizeAuthResponse(auth, {
      authMountPath: path,
      displayName,
      token: auth.client_token,
      ttl: auth.lease_duration,
    });
  }
}
