/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import Route from '@ember/routing/route';
import { service } from '@ember/service';
import SecretsEngineResource from 'vault/resources/secrets/engine';

export default class VaultClusterDashboardRoute extends Route {
  @service store;
  @service namespace;
  @service version;
  @service api;
  @service('checklist-state') checklistState;

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
      this.api.sys.internalUiListEnabledVisibleMounts().catch(() => ({})),
    ];
    const [vaultConfiguration, { secret }] = await Promise.all(requests);
    const secretsEngines = this.api
      .responseObjectToArray(secret, 'path')
      .map((engine) => new SecretsEngineResource(engine));

    const isChecklistNamespaceScope = this.namespace.inRootNamespace && !hasChroot;

    if (
      this.version.isEnterprise &&
      isChecklistNamespaceScope &&
      this.checklistState.hasChecklistEntryAccess
    ) {
      // Fetch persisted checklist state; isAvailable is set to false on 403/404/5xx so
      // the component degrades gracefully without blocking the rest of the page.
      try {
        await this.checklistState.fetchState.perform();
      } catch {
        // Any unexpected error from the checklist fetch must not block page render.
      }
    }

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
