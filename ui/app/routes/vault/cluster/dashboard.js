/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
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
      const adapter = this.store.adapterFor('application');
      const configState = await adapter.ajax('/v1/sys/config/state/sanitized', 'GET');
      return configState.data;
    } catch (e) {
      return null;
    }
  }

  model() {
    const clusterModel = this.modelFor('vault.cluster');
    const replication = {
      dr: clusterModel.dr,
      performance: clusterModel.performance,
    };

    return hash({
      replication,
      secretsEngines: this.store.query('secret-engine', {}),
      license: this.store.queryRecord('license', {}).catch(() => null),
      isRootNamespace: this.namespace.inRootNamespace,
      version: this.version,
      vaultConfiguration: this.getVaultConfiguration(),
    });
  }

  @action
  refreshRoute() {
    this.refresh();
  }
}
