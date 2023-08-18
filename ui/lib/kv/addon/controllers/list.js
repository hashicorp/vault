/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Controller from '@ember/controller';
import { tracked } from '@glimmer/tracking';

export default class KvListController extends Controller {
  queryParams = ['pageFilter', 'currentPage'];
  // ARG TODO does this need to be tracked?
  @tracked currentPage = 1;
}
