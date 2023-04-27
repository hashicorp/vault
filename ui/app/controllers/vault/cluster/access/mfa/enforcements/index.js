/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Controller from '@ember/controller';

export default class MfaEnforcementListController extends Controller {
  queryParams = ['page'];
  page = 1;
}
