/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { CLUSTER } from 'vault/lib/route-paths';

export default class InitRoute extends Route {
  @service router;

  beforeModel() {
    const cluster = this.modelFor(CLUSTER);
    // if it doesn't need init, nav to cluster route
    if (!cluster.needsInit) {
      return this.router.replaceWith(CLUSTER);
    }
  }
}
