/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AuthBase from './base';

/**
 * @module Auth::Form::Token
 * see Auth::Base
 * */

export default class AuthFormToken extends AuthBase {
  loginFields = [{ name: 'token' }];
}
