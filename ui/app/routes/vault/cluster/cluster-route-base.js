// this is the base route for
// all of the CLUSTER_ROUTES that are states before you can use vault
//
import Ember from 'ember';
import ClusterRoute from 'vault/mixins/cluster-route';

export default Ember.Route.extend(ClusterRoute, {
  model() {
    return this.modelFor('vault.cluster');
  },

  resetController(controller) {
    controller.reset && controller.reset();
  },
});
