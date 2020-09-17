import { inject as service } from '@ember/service';
import Controller from '@ember/controller';

export default Controller.extend({
  wizard: service(),

  actions: {
    transitionToCluster(resp) {
      return this.model.reload().then(() => {
        this.wizard.transitionTutorialMachine(this.get('wizard.currentState'), 'CONTINUE', resp);
        return this.transitionToRoute('vault.cluster', this.get('model.name'));
      });
    },

    setUnsealState(resp) {
      this.wizard.set('componentState', resp);
    },

    isUnsealed(data) {
      return data.sealed === false;
    },
  },
});
