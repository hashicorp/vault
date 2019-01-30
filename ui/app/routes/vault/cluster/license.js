import Route from '@ember/routing/route';
import ClusterRoute from 'vault/mixins/cluster-route';
import { inject as service } from '@ember/service';

export default Route.extend(ClusterRoute, {
  version: service(),
  beforeModel() {
    if (this.version.isOSS) {
      this.transitionTo('vault.cluster');
    }
  },

  model() {
    return this.store.queryRecord('license', {});
  },

  actions: {
    doRefresh() {
      this.refresh();
    },
  },
});
