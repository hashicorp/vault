import { inject as service } from '@ember/service';
import Route from '@ember/routing/route';

const SUPPORTED_REPLICATION_MODES = ['dr', 'performance'];

export default Route.extend({
  replicationMode: service(),
  store: service(),
  beforeModel() {
    const replicationMode = this.paramsFor(this.routeName).replication_mode;
    if (!SUPPORTED_REPLICATION_MODES.includes(replicationMode)) {
      return this.transitionTo('index');
    }
  },
  model() {
    const replicationMode = this.paramsFor(this.routeName).replication_mode;
    this.replicationMode.setMode(replicationMode);
    return this.modelFor('application');
  },
});
