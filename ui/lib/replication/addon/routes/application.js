/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import { setProperties } from '@ember/object';
import Route from '@ember/routing/route';
import ClusterRoute from 'vault/mixins/cluster-route';

export default Route.extend(ClusterRoute, {
  version: service(),
  store: service(),
  auth: service(),
  router: service('app-router'),
  capabilities: service(),

  async fetchCapabilities() {
    const enablePath = (type, cluster) => `sys/replication/${type}/${cluster}/enable`;
    const perms = await this.capabilities.fetchMultiplePaths([
      enablePath('dr', 'primary'),
      enablePath('dr', 'primary'),
      enablePath('performance', 'secondary'),
      enablePath('performance', 'secondary'),
    ]);
    return {
      canEnablePrimaryDr: perms[enablePath('dr', 'primary')].canUpdate,
      canEnableSecondaryDr: perms[enablePath('dr', 'primary')].canUpdate,
      canEnablePrimaryPerformance: perms[enablePath('performance', 'secondary')].canUpdate,
      canEnableSecondaryPerformance: perms[enablePath('performance', 'secondary')].canUpdate,
    };
  },

  beforeModel() {
    if (this.auth.activeCluster.replicationRedacted) {
      // disallow replication access if endpoints are redacted
      return this.router.transitionTo('vault.cluster');
    }
    return this.version.fetchFeatures().then(() => {
      return this._super(...arguments);
    });
  },

  model() {
    return this.auth.activeCluster;
  },

  async afterModel(model) {
    const {
      canEnablePrimaryDr,
      canEnableSecondaryDr,
      canEnablePrimaryPerformance,
      canEnableSecondaryPerformance,
    } = await this.fetchCapabilities();

    setProperties(model, {
      canEnablePrimaryDr,
      canEnableSecondaryDr,
      canEnablePrimaryPerformance,
      canEnableSecondaryPerformance,
    });
    return model;
  },
  actions: {
    refresh() {
      this.refresh();
    },
  },
});
