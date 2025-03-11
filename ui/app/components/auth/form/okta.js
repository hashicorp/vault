/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AuthBase from './base';
import { action } from '@ember/object';

/**
 * @module Auth::Form::Okta
 *
 * */

export default class AuthFormOkta extends AuthBase {
  showFields = ['username', 'password'];
  type = 'okta';

  @action
  async login(event) {
    event.preventDefault();
    // do something
  }
}
