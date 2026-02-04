/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { RouteName } from 'core/helpers/display-nav-item';

import type CurrentClusterService from 'vault/services/current-cluster';
import type VersionService from 'vault/services/version';
import type ClusterModel from 'vault/models/cluster';

interface Args {
  isEngine?: boolean;
}

export default class SidebarNavResilienceAndRecoveryComponent extends Component<Args> {
  @service declare readonly currentCluster: CurrentClusterService;
  @service declare readonly version: VersionService;

  routeName = {
    secretsRecovery: RouteName.SECRETS_RECOVERY,
    seal: RouteName.SEAL,
    replication: RouteName.REPLICATION,
  };

  get cluster() {
    return this.currentCluster.cluster as ClusterModel | null;
  }
}
