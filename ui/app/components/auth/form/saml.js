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
      helperText: 'Leave blank to sign in with the default role, if one is configured.',
    },
  ];
}
