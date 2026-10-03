/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class ClusterIndexRoute extends Route {
  @service router;
  beforeModel() {
    return this.router.transitionTo('vault.cluster.dashboard');
  }
}
