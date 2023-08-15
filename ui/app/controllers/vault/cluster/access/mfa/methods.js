/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';

export default class MfaMethodsListController extends Controller {
  queryParams = ['page'];

  page = 1;
}
