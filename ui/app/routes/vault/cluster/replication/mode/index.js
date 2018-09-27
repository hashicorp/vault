import { inject as service } from '@ember/service';
import Route from '@ember/routing/route';

export default Route.extend({
  replicationMode: service(),
  beforeModel() {
    const replicationMode = this.paramsFor('vault.cluster.replication.mode').replication_mode;
    this.get('replicationMode').setMode(replicationMode);
  },
  model() {
    return this.modelFor('vault.cluster.replication.mode');
  },
});
