/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { inject as service } from '@ember/service';
import { setProperties } from '@ember/object';
import { hash } from 'rsvp';
import Route from '@ember/routing/route';
import ClusterRoute from 'vault/mixins/cluster-route';

export default Route.extend(ClusterRoute, {
  version: service(),
  store: service(),
  auth: service(),

  beforeModel() {
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
        .then((c) => c.get('canUpdate')),
      canEnableSecondary: this.store
        .findRecord('capabilities', 'sys/replication/secondary/enable')
        .then((c) => c.get('canUpdate')),
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
