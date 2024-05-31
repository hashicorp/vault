/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';

export default class VaultController extends Controller {
  queryParams = [
    {
      wrappedToken: 'wrapped_token',
      redirectTo: 'redirect_to',
    },
  ];
  wrappedToken = '';
  redirectTo = '';
}
