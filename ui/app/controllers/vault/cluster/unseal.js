import Ember from 'ember';

export default Ember.Controller.extend({
  wizard: Ember.inject.service(),

  actions: {
    transitionToCluster(resp) {
      return this.get('model').reload().then(() => {
        this.get('wizard').transitionTutorialMachine(this.get('wizard.currentState'), 'CONTINUE', resp);
        return this.transitionToRoute('vault.cluster', this.get('model.name'));
      });
    },

    setUnsealState(resp) {
      this.get('wizard').set('componentState', resp);
    },

    isUnsealed(data) {
      return data.sealed === false;
    },
  },
});
