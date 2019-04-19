import { alias } from '@ember/object/computed';
import { inject as service } from '@ember/service';
import Controller from '@ember/controller';

const CONFIG_DEFAULTS = {
  mode: 'whitelist',
  paths: [],
};

export default Controller.extend({
  flashMessages: service(),
  rm: service('replication-mode'),
  replicationMode: alias('rm.mode'),
  actions: {
    resetConfig(config) {
      if (config.get('isNew')) {
        config.setProperties(CONFIG_DEFAULTS);
      } else {
        config.rollbackAttributes();
      }
    },

    saveConfig(config, isDelete) {
      const flash = this.get('flashMessages');
      const id = config.id;
      const redirectArgs = isDelete
        ? [
            'vault.cluster.replication.mode.secondaries',
            this.model.cluster.get('name'),
            this.get('replicationMode'),
          ]
        : ['vault.cluster.replication.mode.secondaries.config-show', id];
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
        .catch(e => {
          const errString = e.errors.join('.');
          flash.error(
            `There was an error ${isDelete ? 'deleting' : 'saving'} the config for ${id}: ${errString}`
          );
        });
    },
  },
});
