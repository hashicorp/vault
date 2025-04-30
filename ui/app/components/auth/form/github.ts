/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AuthBase from './base';

/**
 * @module Auth::Form::Github
 * see Auth::Base
 */

export default class AuthFormGithub extends AuthBase {
  loginFields = [{ name: 'token', label: 'Github token' }];
}
