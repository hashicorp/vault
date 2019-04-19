import { on } from '@ember/object/evented';
import { inject as service } from '@ember/service';
import Route from '@ember/routing/route';

const SUPPORTED_REPLICATION_MODES = ['dr', 'performance'];

export default Route.extend({
  replicationMode: service(),

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

  setReplicationMode: on('activate', 'enter', function() {
    const replicationMode = this.paramsFor(this.routeName).replication_mode;
    this.get('replicationMode').setMode(replicationMode);
  }),
});
