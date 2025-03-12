/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AuthBase from './base';
import { action } from '@ember/object';

/**
 * @module Auth::Form::Github
 * see Auth::Base
 */

export default class AuthFormGithub extends AuthBase {
  loginFields = ['token'];

  @action
  async login(event) {
    event.preventDefault();
    // base login flow
  }
}
