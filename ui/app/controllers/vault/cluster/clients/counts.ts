/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';

export default class ClientsCountsController extends Controller {
  queryParams = ['start_time', 'end_time', 'ns', 'authMount'];
}
