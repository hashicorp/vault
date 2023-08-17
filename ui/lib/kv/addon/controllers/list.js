/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Controller from '@ember/controller';
import { tracked } from '@glimmer/tracking';

export default class KvListController extends Controller {
  queryParams = ['pageFilter', 'currentPage'];

  @tracked currentPage = 1;
}
