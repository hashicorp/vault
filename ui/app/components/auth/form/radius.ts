/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AuthBase from './base';

/**
 * @module Auth::Form::Radius
 * see Auth::Base
 */

export default class AuthFormRadius extends AuthBase {
  loginFields = [{ name: 'username' }, { name: 'password' }];
}
