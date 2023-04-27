/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Controller from '@ember/controller';

export default class MfaMethodsListController extends Controller {
  queryParams = ['page'];

  page = 1;
}
