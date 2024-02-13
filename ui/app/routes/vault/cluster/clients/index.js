/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class ClientsIndexRoute extends Route {
  @service router;

  redirect() {
    this.router.transitionTo('vault.cluster.clients.counts.overview');
  }
}
