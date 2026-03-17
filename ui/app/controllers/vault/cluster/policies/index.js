/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */
import Controller from '@ember/controller';
import { action } from '@ember/object';

export default class ManagePoliciesController extends Controller {
  queryParams = ['page', 'pageFilter'];
  page = 1;
  pageFilter = null;

  @action
  refreshRoute() {
    this.send('reload');
  }
}
