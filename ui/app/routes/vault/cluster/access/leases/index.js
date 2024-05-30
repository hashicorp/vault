/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class LeasesIndexRoute extends Route {
  @service router;

  beforeModel(transition) {
    if (this.modelFor('vault.cluster.access.leases').canList && transition.targetName === this.routeName) {
      return this.router.replaceWith('vault.cluster.access.leases.list-root');
    } else {
      return;
    }
  }
}
