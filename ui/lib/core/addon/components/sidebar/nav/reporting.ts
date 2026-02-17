/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { RouteName } from 'core/helpers/display-nav-item';

import type CurrentClusterService from 'vault/services/current-cluster';
import type ClusterModel from 'vault/models/cluster';

interface Args {
  isEngine?: boolean;
}

export default class SidebarNavReportingComponent extends Component<Args> {
  @service declare readonly currentCluster: CurrentClusterService;

  routeName = {
    vaultUsage: RouteName.VAULT_USAGE,
    license: RouteName.LICENSE,
  };

  get cluster() {
    return this.currentCluster.cluster as ClusterModel | null;
  }
}
