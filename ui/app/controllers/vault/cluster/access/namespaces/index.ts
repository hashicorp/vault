/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import { service } from '@ember/service';
import Controller from '@ember/controller';
import { tracked } from '@glimmer/tracking';

import type RouterService from '@ember/routing/router-service';

export default class NamespacesIndexController extends Controller {
  @service declare readonly router: RouterService;

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

  @action
  refreshRoute() {
    this.router.refresh('vault.cluster.access.namespaces.index');
  }
}
