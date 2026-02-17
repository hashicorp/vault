/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Helper from '@ember/component/helper';
import { service } from '@ember/service';
import { getOwner, setOwner } from '@ember/owner';

import type CurrentClusterService from 'vault/services/current-cluster';
import type VersionService from 'vault/services/version';
import type NamespaceService from 'vault/services/namespace';
import type ClusterModel from 'vault/models/cluster';
import type PermissionsService from 'vault/services/permissions';
import type FlagsService from 'vault/services/flags';

export enum RouteName {
  SECRETS_SYNC = 'secrets-sync',
  SECRETS_RECOVERY = 'secrets-recovery',
  SEAL = 'seal',
  REPLICATION = 'replication',
  VAULT_USAGE = 'vault-usage',
  LICENSE = 'license',
}

export enum NavSection {
  RESILIENCE_AND_RECOVERY = 'resilience-and-recovery',
  REPORTING = 'reporting',
  CLIENT_COUNT = 'client-count',
}

export default class NavBar extends Helper {
  @service declare readonly currentCluster: CurrentClusterService;
  @service declare readonly version: VersionService;
  @service declare readonly namespace: NamespaceService;
  @service declare readonly permissions: PermissionsService;
  @service declare readonly flags: FlagsService;

  compute([navItem]: string[]) {
    const { SECRETS_RECOVERY, SEAL, REPLICATION, VAULT_USAGE, LICENSE, SECRETS_SYNC } = RouteName;
    const { RESILIENCE_AND_RECOVERY, REPORTING, CLIENT_COUNT } = NavSection;

    switch (navItem) {
      // secrets sync nav items
      case SECRETS_SYNC:
        return this.supportsSecretsSync;
      // client count nav items
      case CLIENT_COUNT:
        return this.supportsClientCount;
      // reporting nav items
      case VAULT_USAGE:
        return this.canAccessVaultUsageDashboard;
      case LICENSE:
        return this.supportsLicense;
      case REPORTING:
        return this.canAccessVaultUsageDashboard || this.supportsLicense;
      // resilience and recovery nav items
      case SECRETS_RECOVERY:
        return this.supportsSnapshots;
      case SEAL:
        return this.canSeal;
      case REPLICATION:
        return this.supportsReplication;
      case RESILIENCE_AND_RECOVERY:
        return this.supportsSnapshots || this.canSeal || this.supportsReplication;
      default:
        return true;
    }
  }

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

  get supportsReplication() {
    return (
      this.version.isEnterprise &&
      this.isRootNamespace &&
      !this.cluster?.replicationRedacted &&
      this.permissions.hasNavPermission('status', 'replication')
    );
  }

  get canSeal() {
    return (
      this.isRootNamespace &&
      this.permissions.hasNavPermission('status', 'seal') &&
      !this.cluster?.dr?.isSecondary
    );
  }

  get supportsSnapshots() {
    return (this.cluster && !this.cluster?.dr?.isSecondary) || !this.flags.isHvdManaged;
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

  get supportsLicense() {
    return (
      this.version.isEnterprise &&
      this?.version?.features &&
      this.isRootNamespace &&
      this.permissions.hasNavPermission('status', 'license') &&
      !this.cluster?.dr?.isSecondary
    );
  }

  get supportsClientCount() {
    return (
      this.permissions.hasNavPermission('clients', 'activity') &&
      !this.cluster?.dr?.isSecondary &&
      !this.hasChrootNamespace &&
      !this.version.hasPKIOnly
    );
  }

  get supportsSecretsSync() {
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
}

export function computeNavBar(context: object, navItem: string): boolean {
  const navBar = new NavBar();
  const owner = getOwner(context);

  if (owner) {
    setOwner(navBar, owner);
    return navBar.compute([navItem]);
  }
  return true; // when in doubt
}
