/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';

export default class AgentsRegistryController extends Controller {
  queryParams = ['page', 'pageSize'];
  page = 1;
  pageSize = 10;
}
