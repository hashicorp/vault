/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';

export default class MfaEnforcementListController extends Controller {
  queryParams = ['page'];
  page = 1;
}
