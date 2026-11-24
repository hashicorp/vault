/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

export default class PoliciesAclIndexController extends Controller {
  // page refreshes the model (re-fetches from API) on change.
  // pageSize, sortBy, sortOrder are client-side only — no model reload needed.
  queryParams = ['page', 'pageSize', 'sortBy', 'sortOrder'];
  @tracked page = 1;
  @tracked pageSize = 10;
  @tracked sortBy: string | undefined = undefined;
  @tracked sortOrder: 'asc' | 'desc' | undefined = undefined;

  @action
  updateSort(sortBy: string, sortOrder: 'asc' | 'desc') {
    this.sortBy = sortBy;
    this.sortOrder = sortOrder;
  }
}
