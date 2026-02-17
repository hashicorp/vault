/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */
import Route from '@ember/routing/route';
import RouterService from '@ember/routing/router-service';
import { service } from '@ember/service';
import { computeNavBar, RouteName } from 'core/helpers/display-nav-item';

import type CurrentClusterService from 'vault/services/current-cluster';
import type ClusterModel from 'vault/models/cluster';

export default class ResilienceRecoveryRoute extends Route {
  @service declare readonly router: RouterService;
  @service declare readonly currentCluster: CurrentClusterService;

  beforeModel() {
    const cluster = this.currentCluster.cluster as ClusterModel | null;

    if (computeNavBar(this, RouteName.SECRETS_RECOVERY)) {
      this.router.replaceWith('vault.cluster.recovery.snapshots');
    } else if (computeNavBar(this, RouteName.SEAL)) {
      this.router.replaceWith('vault.cluster.settings.seal', cluster?.name);
    } else if (computeNavBar(this, RouteName.REPLICATION)) {
      this.router.replaceWith('vault.cluster.replication.index');
    }
  }
}
