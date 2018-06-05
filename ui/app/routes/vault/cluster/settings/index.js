import Ember from 'ember';

export default Ember.Route.extend({
  beforeModel: function(transition) {
    if (transition.targetName === this.routeName) {
      transition.abort();
      this.replaceWith('vault.cluster.settings.mount-secret-backend');
    }
  },
});
