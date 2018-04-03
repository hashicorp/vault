import Ember from 'ember';
import ClusterRoute from 'vault/mixins/cluster-route';

export default Ember.Route.extend(ClusterRoute, {
  model() {
    return this.store.query('secret-engine', {});
  },
});
