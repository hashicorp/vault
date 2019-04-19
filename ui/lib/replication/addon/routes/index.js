import { inject as service } from '@ember/service';
import Route from '@ember/routing/route';

export default Route.extend({
  replicationMode: service(),
  beforeModel() {
    this.get('replicationMode').setMode(null);
  },
  model() {
    return this.modelFor('vault.cluster.replication');
  },
});
