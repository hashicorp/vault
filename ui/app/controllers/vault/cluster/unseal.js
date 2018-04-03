import Ember from 'ember';

export default Ember.Controller.extend({
  actions: {
    transitionToCluster() {
      return this.get('model').reload().then(() => {
        return this.transitionToRoute('vault.cluster', this.get('model.name'));
      });
    },
    isUnsealed(data) {
      return data.sealed === false;
    },
  },
});
