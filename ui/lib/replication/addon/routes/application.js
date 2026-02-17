/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Route from '@ember/routing/route';

export default class ApplicationRoute extends Route {
  @service version;
  @service store;
  @service auth;
  @service('app-router') router;
  @service capabilities;

  async fetchCapabilities() {
    const enablePath = (type, cluster) => `sys/replication/${type}/${cluster}/enable`;
    const perms = await this.capabilities.fetch([
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
  }

  beforeModel() {
    if (this.auth.activeCluster.replicationRedacted) {
      // disallow replication access if endpoints are redacted
      return this.router.transitionTo('vault.cluster');
    }
    return this.version.fetchFeatures();
  }

  model() {
    return this.auth.activeCluster;
  }

  async afterModel(model) {
    const {
      canEnablePrimaryDr,
      canEnableSecondaryDr,
      canEnablePrimaryPerformance,
      canEnableSecondaryPerformance,
    } = await this.fetchCapabilities();

    model.canEnablePrimaryDr = canEnablePrimaryDr;
    model.canEnableSecondaryDr = canEnableSecondaryDr;
    model.canEnablePrimaryPerformance = canEnablePrimaryPerformance;
    model.canEnableSecondaryPerformance = canEnableSecondaryPerformance;
    return model;
  }
}
