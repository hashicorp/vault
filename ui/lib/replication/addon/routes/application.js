/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import { setProperties } from '@ember/object';
import { hash } from 'rsvp';
import Route from '@ember/routing/route';
import ClusterRoute from 'vault/mixins/cluster-route';

export default Route.extend(ClusterRoute, {
  version: service(),
  store: service(),
  auth: service(),
  router: service(),

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

  afterModel(model) {
    return hash({
      canEnablePrimaryPerformance: this.store
        .findRecord('capabilities', 'sys/replication/performance/primary/enable')
        .then((c) => c.canUpdate),
      canEnableSecondaryPerformance: this.store
        .findRecord('capabilities', 'sys/replication/performance/secondary/enable')
        .then((c) => c.canUpdate),
      canEnablePrimaryDr: this.store
        .findRecord('capabilities', 'sys/replication/dr/primary/enable')
        .then((c) => c.canUpdate),
      canEnableSecondaryDr: this.store
        .findRecord('capabilities', 'sys/replication/dr/secondary/enable')
        .then((c) => c.canUpdate),
    }).then(
      ({
        canEnablePrimaryPerformance,
        canEnableSecondaryPerformance,
        canEnablePrimaryDr,
        canEnableSecondaryDr,
      }) => {
        setProperties(model, {
          canEnablePrimaryPerformance,
          canEnableSecondaryPerformance,
          canEnablePrimaryDr,
          canEnableSecondaryDr,
        });
        return model;
      }
    );
  },
  actions: {
    refresh() {
      this.refresh();
    },
  },
});
