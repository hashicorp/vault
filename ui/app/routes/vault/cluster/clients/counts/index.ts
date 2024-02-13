/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

import type RouterService from '@ember/routing/router-service';

export default class ClientsCountsOverviewRoute extends Route {
  @service declare readonly router: RouterService;

  redirect() {
    this.router.transitionTo('vault.cluster.clients.counts.overview');
  }
}
