import { inject as service } from '@ember/service';
import Route from '@ember/routing/route';

const SUPPORTED_REPLICATION_MODES = ['dr', 'performance'];

export default Route.extend({
  replicationMode: service(),
  store: service(),
  model(params) {
    const replicationMode = params.replication_mode;

    if (!SUPPORTED_REPLICATION_MODES.includes(replicationMode)) {
      return this.transitionTo('application');
    } else {
      this.replicationMode.setMode(replicationMode);
      return this.modelFor('application');
    }
  },
});
