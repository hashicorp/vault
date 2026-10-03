/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';

import type CurrentClusterService from 'vault/services/current-cluster';
import type VersionService from 'vault/services/version';
import type NamespaceService from 'vault/services/namespace';
import type ClusterModel from 'vault/models/cluster';
import type PermissionsService from 'vault/services/permissions';

interface Args {
  isEngine?: boolean;
}

export default class SidebarNavReportingComponent extends Component<Args> {
  @service declare readonly currentCluster: CurrentClusterService;
  @service declare readonly version: VersionService;
  @service declare readonly namespace: NamespaceService;
  @service declare readonly permissions: PermissionsService;

  get cluster() {
    return this.currentCluster.cluster as ClusterModel | null;
  }

  get hasChrootNamespace() {
    return this.cluster?.hasChrootNamespace;
  }

  get isRootNamespace() {
    // should only return true if we're in the true root namespace
    return this.namespace.inRootNamespace && !this.hasChrootNamespace;
  }

  get canAccessVaultUsageDashboard() {
    /*
    A user can access Vault Usage if they satisfy the following conditions:
      1) They have access to sys/v1/utilization-report endpoint
      2) They are either
        a) enterprise cluster and root namespace
        b) hvd cluster and /admin namespace
    */

    const hasPermission = this.permissions.hasNavPermission('monitoring');
    const isEnterprise = this.version.isEnterprise;
    const isCorrectNamespace = this.isRootNamespace || this.namespace.inHvdAdminNamespace;

    return hasPermission && isEnterprise && isCorrectNamespace;
  }
}
