/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { hash } from 'rsvp';
// eslint-disable-next-line ember/no-mixins
import ClusterRoute from 'vault/mixins/cluster-route';
import { action } from '@ember/object';

export default class VaultClusterDashboardRoute extends Route.extend(ClusterRoute) {
  @service store;
  @service namespace;
  @service version;

  async getVaultConfiguration() {
    try {
      if (!this.namespace.inRootNamespace) return null;

      const adapter = this.store.adapterFor('application');
      const configState = await adapter.ajax('/v1/sys/config/state/sanitized', 'GET');
      return configState.data;
    } catch (e) {
      return null;
    }
  }

  model() {
    const clusterModel = this.modelFor('vault.cluster');
    const hasChroot = clusterModel?.hasChrootNamespace;
    const replication =
      hasChroot || clusterModel.replicationRedacted
        ? null
        : {
            dr: clusterModel.dr,
            performance: clusterModel.performance,
          };
    return hash({
      replication,
      secretsEngines: this.store.query('secret-engine', {}),
      license: this.store.queryRecord('license', {}).catch(() => null),
      isRootNamespace: this.namespace.inRootNamespace && !hasChroot,
      version: this.version,
      vaultConfiguration: hasChroot ? null : this.getVaultConfiguration(),
    });
  }

  @action
  refreshRoute() {
    this.refresh();
  }
}
