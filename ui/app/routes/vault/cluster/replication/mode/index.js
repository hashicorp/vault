import Ember from 'ember';

export default Ember.Route.extend({
  replicationMode: Ember.inject.service(),
  beforeModel() {
    const replicationMode = this.paramsFor('vault.cluster.replication.mode').replication_mode;
    this.get('replicationMode').setMode(replicationMode);
  },
  model() {
    return this.modelFor('vault.cluster.replication.mode');
  },
});
