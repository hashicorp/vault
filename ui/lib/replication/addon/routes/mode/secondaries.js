/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { setProperties } from '@ember/object';
import { hash } from 'rsvp';
import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default Route.extend({
  router: service(),
  store: service(),
  model() {
    const replicationMode = this.paramsFor('mode').replication_mode;

    return hash({
      cluster: this.modelFor('mode'),
      canAddSecondary: this.store
        .findRecord('capabilities', `sys/replication/${replicationMode}/primary/secondary-token`)
        .then((c) => c.canUpdate),
      canRevokeSecondary: this.store
        .findRecord('capabilities', `sys/replication/${replicationMode}/primary/revoke-secondary`)
        .then((c) => c.canUpdate),
    }).then(({ cluster, canAddSecondary, canRevokeSecondary }) => {
      setProperties(cluster, {
        canRevokeSecondary,
        canAddSecondary,
      });
      return cluster;
    });
  },
  afterModel(model) {
    const replicationMode = this.paramsFor('mode').replication_mode;
    const modeModel = model[replicationMode];
    if (!modeModel.isPrimary || modeModel.replicationDisabled || modeModel.replicationUnsupported) {
      this.router.transitionTo('vault.cluster.replication.mode', replicationMode);
    }
  },
});
