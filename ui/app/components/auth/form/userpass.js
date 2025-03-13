/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AuthBase from './base';

/**
 * @module Auth::Form::Userpass
 *
 * */

export default class AuthFormUserpass extends AuthBase {
  loginFields = ['username', 'password'];
}
