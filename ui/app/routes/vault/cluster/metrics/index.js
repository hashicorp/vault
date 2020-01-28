import Route from '@ember/routing/route';
import ClusterRoute from 'vault/mixins/cluster-route';
import { hash } from 'rsvp';

export default Route.extend(ClusterRoute, {
  model() {
    // return this.store.queryRecord('metrics', {});
    return hash({
      tokens: this.store.queryRecord('tokens', {}), // ARG: calling model token, model and adapter are related.
      https: this.store.queryRecord('http-requests', {}),
    });
  },
});
