/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AuthBase from './base';

/**
 * @module Auth::Form::Saml
 * see Auth::Base
 */

export default class AuthFormSaml extends AuthBase {
  loginFields = [
    {
      name: 'role',
      helperText: 'Vault will use the default role to sign in if this field is left blank.',
    },
  ];
}
