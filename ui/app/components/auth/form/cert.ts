/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import AuthBase from './base';

import type { CertLoginApiResponse } from 'vault/vault/auth/methods';

/**
 * @module Auth::Form::Cert
 * see Auth::Base
 * */

export default class AuthFormCert extends AuthBase {
  loginFields = [
    {
      name: 'name',
      label: 'Role name',
      helperText: 'Leave blank to match any certificate role.',
    },
  ];

  async loginRequest(formData: { name: string; path: string }) {
    const { path, name } = formData;

    const { auth } = (await this.api.auth.certLogin(path, {
      name,
    })) as CertLoginApiResponse;

    const { cert_name, common_name } = auth?.metadata || {};
    const displayName =
      cert_name && common_name ? `${cert_name}/${common_name}` : cert_name || common_name || '';

    return this.normalizeAuthResponse(auth, {
      authMountPath: path,
      displayName,
      token: auth.client_token,
      ttl: auth.lease_duration,
    });
  }
}
