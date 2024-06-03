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
      canEnablePrimary: this.store
        .findRecord('capabilities', 'sys/replication/primary/enable')
        .then((c) => c.canUpdate),
      canEnableSecondary: this.store
        .findRecord('capabilities', 'sys/replication/secondary/enable')
        .then((c) => c.canUpdate),
    }).then(({ canEnablePrimary, canEnableSecondary }) => {
      setProperties(model, {
        canEnablePrimary,
        canEnableSecondary,
      });
      return model;
    });
  },
  actions: {
    refresh() {
      this.refresh();
    },
  },
});
