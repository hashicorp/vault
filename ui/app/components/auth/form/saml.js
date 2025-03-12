/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AuthBase from './base';
import { action } from '@ember/object';

/**
 * @module Auth::Form::Saml
 * see Auth::Base
 */

export default class AuthFormSaml extends AuthBase {
  loginFields = ['role'];

  @action
  async login(event) {
    event.preventDefault();
    // base login flow
  }
}
