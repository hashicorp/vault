/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
// eslint-disable-next-line ember/no-mixins
import ClusterRoute from 'vault/mixins/cluster-route';
import { action } from '@ember/object';
import SecretsEngineResource from 'vault/resources/secrets/engine';

export default class VaultClusterDashboardRoute extends Route.extend(ClusterRoute) {
  @service store;
  @service namespace;
  @service version;
  @service api;

  async getVaultConfiguration(hasChroot) {
    try {
      if (!this.namespace.inRootNamespace || hasChroot) {
        return null;
      }
      const { data } = await this.api.sys.readSanitizedConfigurationState();
      return data;
    } catch (e) {
      return null;
    }
  }

  async model() {
    const clusterModel = this.modelFor('vault.cluster');
    const adapter = this.store.adapterFor('application');
    const hasChroot = clusterModel?.hasChrootNamespace;
    const replication =
      hasChroot || clusterModel.replicationRedacted
        ? null
        : {
            dr: clusterModel.dr,
            performance: clusterModel.performance,
          };
    const requests = [
      this.getVaultConfiguration(hasChroot),
      adapter.ajax('/v1/sys/internal/ui/mounts', 'GET').catch(() => ({})),
    ];
    const [vaultConfiguration, { data }] = await Promise.all(requests);
    const secret = data.secret;
    const secretsEngines = this.api
      .responseObjectToArray(secret, 'path')
      .map((engine) => new SecretsEngineResource(engine));

    return {
      replication,
      secretsEngines,
      isRootNamespace: this.namespace.inRootNamespace && !hasChroot,
      version: this.version,
      vaultConfiguration,
    };
  }

  @action
  refreshRoute() {
    this.refresh();
  }
}
