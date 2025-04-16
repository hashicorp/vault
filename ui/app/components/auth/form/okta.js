/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AuthBase from './base';

/**
 * @module Auth::Form::Okta
 * see Auth::Base
 * */

export default class AuthFormOkta extends AuthBase {
  loginFields = [{ name: 'username' }, { name: 'password' }];
}
