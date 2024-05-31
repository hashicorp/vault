/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { hash } from 'rsvp';
import Base from '../../replication-base';

export default Base.extend({
  modelPath: 'model.config',

  model(params) {
    return hash({
      cluster: this.modelFor('mode.secondaries'),
      config: this.store.findRecord('path-filter-config', params.secondary_id),
    });
  },

  redirect(model) {
    const cluster = model.cluster;
    const replicationMode = this.paramsFor('mode').replication_mode;
    if (
      !this.version.hasPerfReplication ||
      replicationMode !== 'performance' ||
      !cluster[replicationMode].isPrimary ||
      !cluster.canAddSecondary
    ) {
      return this.router.transitionTo('vault.cluster.replication.mode', replicationMode);
    }
  },
});
