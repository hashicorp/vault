import Ember from 'ember';

export default Ember.Route.extend({
  beforeModel() {
    return this.replaceWith('vault.cluster.settings.auth.enable');
  },
});
