import { alias } from '@ember/object/computed';
import { inject as service } from '@ember/service';
import Controller from '@ember/controller';

export default Controller.extend({
  rm: service('replication-mode'),
  replicationMode: alias('rm.mode'),
  actions: {
    onEnable(mode) {
      return this.transitionToRoute('mode', mode);
    },
    onDisable() {
      return this.transitionToRoute('index');
    },
  },
});
