import Ember from 'ember';

const SUPPORTED_REPLICATION_MODES = ['dr', 'performance'];

export default Ember.Route.extend({
  replicationMode: Ember.inject.service(),

  beforeModel() {
    const replicationMode = this.paramsFor(this.routeName).replication_mode;
    if (!SUPPORTED_REPLICATION_MODES.includes(replicationMode)) {
      return this.transitionTo('vault.cluster.replication');
    } else {
      return this._super(...arguments);
    }
  },

  model() {
    return this.modelFor('vault.cluster.replication');
  },

  setReplicationMode: Ember.on('activate', 'enter', function() {
    const replicationMode = this.paramsFor(this.routeName).replication_mode;
    this.get('replicationMode').setMode(replicationMode);
  }),
});
