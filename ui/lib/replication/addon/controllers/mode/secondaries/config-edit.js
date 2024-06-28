/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { alias } from '@ember/object/computed';
import { service } from '@ember/service';
import Controller from '@ember/controller';

export default Controller.extend({
  flashMessages: service(),
  router: service(),
  rm: service('replication-mode'),
  replicationMode: alias('rm.mode'),
  actions: {
    resetConfig(config) {
      if (config.isNew) {
        config.setProperties({
          mode: null,
          paths: [],
        });
      } else {
        config.rollbackAttributes();
      }
    },

    saveConfig(config) {
      // if the mode is null, we want no filtering, so we should delete any existing config
      const isDelete = config.mode === null;
      const flash = this.flashMessages;
      const id = config.id;
      const redirectArgs = isDelete
        ? ['vault.cluster.replication.mode.secondaries', this.replicationMode]
        : ['vault.cluster.replication.mode.secondaries.config-show', id];
      const modelMethod = isDelete ? config.destroyRecord : config.save;

      modelMethod
        .call(config)
        .then(() => {
          this.router
            .transitionTo(...redirectArgs)
            .followRedirects()
            .then(() => {
              flash.success(
                `The performance mount filter config for the secondary ${id} was successfully ${
                  isDelete ? 'deleted' : 'saved'
                }.`
              );
            });
        })
        .catch((e) => {
          const errString = e.errors.join('.');
          flash.error(
            `There was an error ${isDelete ? 'deleting' : 'saving'} the config for ${id}: ${errString}`
          );
        });
    },
  },
});
