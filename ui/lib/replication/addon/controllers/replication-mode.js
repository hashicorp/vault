import { alias } from '@ember/object/computed';
import { inject as service } from '@ember/service';
import Controller from '@ember/controller';

export default Controller.extend({
  rm: service('replication-mode'),
  replicationMode: alias('rm.mode'),
  actions: {
    onEnable(replicationMode, mode) {
      if (replicationMode == 'dr' && mode === 'secondary') {
        this.router.transitionTo('vault.cluster');
      }
      return this.transitionToRoute('mode', replicationMode);
    },
    onDisable() {
      return this.transitionToRoute('index');
    },
  },
});
