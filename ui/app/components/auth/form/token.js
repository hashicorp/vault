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
  loginFields = ['token'];

  // idea: future improvement
  // instead of the auth service "authenticate" calling await adapter.authenticate(options);
  // which is prepared for any/all methods
  // specific method data (like from SUPPORTED_AUTH_BACKENDS) could live here
  // and remove the need for the authenticate method in the cluster.js adapter

  // url = '/v1/auth/token/lookup-self';
  // displayNamePath = 'display_name';
  // tokenPath = 'id';
}
