import Ember from 'ember';

export default Ember.Controller.extend({
  wizard: Ember.inject.service(),
  actions: {
    onMountSuccess: function(type) {
      let transition = this.transitionToRoute('vault.cluster.access.methods');
      return transition.followRedirects().then(() => {
        this.get('wizard').transitionFeatureMachine(this.get('wizard.featureState'), 'CONTINUE', type);
      });
    },
    onConfigError: function(modelId) {
      return this.transitionToRoute('vault.cluster.settings.auth.configure', modelId);
    },
  },
});
