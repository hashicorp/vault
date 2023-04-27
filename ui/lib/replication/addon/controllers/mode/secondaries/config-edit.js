/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { alias } from '@ember/object/computed';
import { inject as service } from '@ember/service';
import Controller from '@ember/controller';

export default Controller.extend({
  flashMessages: service(),
  rm: service('replication-mode'),
  replicationMode: alias('rm.mode'),
  actions: {
    resetConfig(config) {
      if (config.get('isNew')) {
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
        ? ['mode.secondaries', this.replicationMode]
        : ['mode.secondaries.config-show', id];
      const modelMethod = isDelete ? config.destroyRecord : config.save;

      modelMethod
        .call(config)
        .then(() => {
          this.transitionToRoute(...redirectArgs)
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
