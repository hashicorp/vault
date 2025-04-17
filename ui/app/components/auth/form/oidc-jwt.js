/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AuthBase from './base';

/**
 * @module Auth::Form::OidcJwt
 * see Auth::Base
 *
 * OIDC can be configured at 'jwt' or 'oidc', see https://developer.hashicorp.com/vault/docs/auth/jwt
 * we use the same template because displaying the JWT token input depends on the error message returned when fetching
 * the role
 */

export default class AuthFormOidcJwt extends AuthBase {
  loginFields = [
    {
      name: 'role',
      helperText: 'Vault will use the default role to sign in if this field is left blank.',
    },
  ];
}
