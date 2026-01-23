/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';

export default class SidebarNavClusterComponent extends Component {
  @service currentCluster;
  @service flags;
  @service version;
  @service auth;
  @service namespace;
  @service permissions;

  get cluster() {
    return this.currentCluster.cluster;
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

  get showSecretsSync() {
    // always show for HVD managed clusters
    if (this.flags.isHvdManaged) return true;

    if (this.flags.secretsSyncIsActivated) {
      // activating the feature requires different permissions than using the feature.
      // we want to show the link to allow activation regardless of permissions to sys/sync
      // and only check permissions if the feature has been activated
      return this.permissions.hasNavPermission('sync');
    }

    // otherwise we show the link depending on whether or not the feature exists
    return this.version.hasSecretsSync;
  }

  get accessRoute() {
    if (this.permissions.hasPermission('policies')) {
      return 'vault.cluster.policies';
    }

    if (this.permissions.hasPermission('access')) {
      return 'vault.cluster.access';
    }

    return null;
  }

  get accessRouteModels() {
    if (this.permissions.hasPermission('policies')) {
      return this.routeParamsFor('policies').models;
    }
    if (this.permissions.hasPermission('access')) {
      return this.routeParamsFor('access').models;
    }

    return null;
  }

  routeParamsFor(routeName) {
    return this.permissions.navPathParams(routeName);
  }
}
