import Ember from 'ember';

export default Ember.Controller.extend({
  actions: {
    onMountSuccess: function() {
      return this.transitionToRoute('vault.cluster.access.methods');
    },
    onConfigError: function(modelId) {
      return this.transitionToRoute('vault.cluster.settings.auth.configure', modelId);
    },
  },
});
