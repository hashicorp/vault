/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AuthBase from './base';
import { action } from '@ember/object';

/**
 * @module Auth::Form::Ldap
 * see Auth::Base
 */

export default class AuthFormLdap extends AuthBase {
  loginFields = ['username', 'password'];

  @action
  async login(event) {
    event.preventDefault();
    // base login flow
  }
}
