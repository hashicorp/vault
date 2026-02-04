/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { NavSection } from 'core/helpers/display-nav-item';

export default class SidebarNavClusterComponent extends Component {
  @service currentCluster;
  @service flags;
  @service version;
  @service namespace;
  @service permissions;

  navSection = {
    resilienceAndRecovery: NavSection.RESILIENCE_AND_RECOVERY,
    reporting: NavSection.REPORTING,
    clientCount: NavSection.CLIENT_COUNT,
  };

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
