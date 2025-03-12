/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AuthBase from './base';
import { action } from '@ember/object';

// TODO separate these?? they use the same endpoint, so I think it makes sense to keep as one...tbd
/**
 * @module Auth::Form::OidcJwt
 * see Auth::Base
 */

export default class AuthFormOidcJwt extends AuthBase {
  loginFields = ['role'];

  @action
  async login(event) {
    event.preventDefault();
    // base login flow
  }
}
