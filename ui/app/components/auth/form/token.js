/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import AuthBase from './base';

/**
 * @module Auth::Form::Token
 *
 * */

export default class AuthFormToken extends AuthBase {
  showFields = ['token'];

  @action
  authenticate(event) {
    // do something
  }
}
