/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AuthBase from './base';

/**
 * @module Auth::Form::Userpass
 * see Auth::Base
 * */

export default class AuthFormUserpass extends AuthBase {
  loginFields = [{ name: 'username' }, { name: 'password' }];
}
