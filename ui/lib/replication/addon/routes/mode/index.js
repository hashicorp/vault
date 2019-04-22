import { inject as service } from '@ember/service';
import Route from '@ember/routing/route';

export default Route.extend({
  replicationMode: service(),
  model() {
    const replicationMode = this.paramsFor('mode').replication_mode;
    this.get('replicationMode').setMode(replicationMode);
    return this.modelFor('mode');
  },
});
