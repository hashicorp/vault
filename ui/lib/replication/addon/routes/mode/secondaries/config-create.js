/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { hash } from 'rsvp';
import { service } from '@ember/service';
import Base from '../../replication-base';

export default Base.extend({
  flashMessages: service(),
  router: service(),

  modelPath: 'model.config',

  findOrCreate(id) {
    const flash = this.flashMessages;
    return this.store
      .findRecord('path-filter-config', id)
      .then(() => {
        // if we find a record, transition to the edit view
        return this.router
          .transitionTo('vault.cluster.replication.mode.secondaries.config-edit', id)
          .followRedirects()
          .then(() => {
            flash.info(
              `${id} already had a path filter config, so we loaded the config edit screen for you.`
            );
          });
      })
      .catch((e) => {
        if (e.httpStatus === 404) {
          return this.store.createRecord('path-filter-config', {
            id,
            mode: null,
            paths: [],
          });
        } else {
          throw e;
        }
      });
  },

  redirect(model) {
    const cluster = model.cluster;
    const replicationMode = this.replicationMode;
    if (
      !this.version.hasPerfReplication ||
      replicationMode !== 'performance' ||
      !cluster[replicationMode].isPrimary ||
      !cluster.canAddSecondary
    ) {
      return this.router.transitionTo('vault.cluster.replication.mode', replicationMode);
    }
  },

  model(params) {
    return hash({
      cluster: this.modelFor('mode'),
      config: this.findOrCreate(params.secondary_id),
    });
  },
});
