import { alias } from '@ember/object/computed';
import { inject as service } from '@ember/service';
import Controller from '@ember/controller';
import { task, timeout } from 'ember-concurrency';

export default Controller.extend({
  rm: service('replication-mode'),
  replicationMode: alias('rm.mode'),
  waitForNewClusterToInit: task(function* (replicationMode) {
    // waiting for the newly enabled cluster to init
    // this ensures we don't hit a capabilities-self error, called in the model of the mode/index route
    yield timeout(1000);
    return this.transitionToRoute('mode', replicationMode);
  }),
  actions: {
    onEnable(replicationMode, mode) {
      if (replicationMode == 'dr' && mode === 'secondary') {
        return this.transitionToRoute('vault.cluster');
      } else if (replicationMode === 'dr') {
        return this.transitionToRoute('mode', replicationMode);
      } else {
        this.waitForNewClusterToInit.perform(replicationMode);
      }
    },
    onDisable() {
      return this.transitionToRoute('index');
    },
  },
});
