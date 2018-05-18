import Ember from 'ember';

export default Ember.Route.extend({
  model(params) {
    return this.modelFor('vault.cluster.secrets.backend');
  },

});
