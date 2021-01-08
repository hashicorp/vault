import Route from '@ember/routing/route';
import ClusterRoute from 'vault/mixins/cluster-route';
import { hash } from 'rsvp';

// TODO query permissions!

export default Route.extend(ClusterRoute, {
  model() {
    // let { backend } = this.paramsFor('vault.cluster.secrets.backend');

    // FUTURE find records for the models

    return hash({
      connections: [],
      roles: [],
    });
  },
});
