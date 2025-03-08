/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import AuthTemplate from './base';

/**
 * @module Auth::Form::Userpass
 * */

export default class AuthFormUserpass extends AuthTemplate {
  showFields = ['username', 'password'];

  @action
  authenticate(event) {
    // do something
  }
}
