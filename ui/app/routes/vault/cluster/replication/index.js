import Ember from 'ember';

export default Ember.Route.extend({
  replicationMode: Ember.inject.service(),
  beforeModel() {
    this.get('replicationMode').setMode(null);
  },
  model() {
    return this.modelFor('vault.cluster.replication');
  },
});
