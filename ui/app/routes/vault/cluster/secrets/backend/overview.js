import Route from '@ember/routing/route';
import ClusterRoute from 'vault/mixins/cluster-route';
import { hash } from 'rsvp';

// ARG TODO query permissions!

export default Route.extend(ClusterRoute, {
  model() {
    let connection = this.store.query('database/connection', {});
    let role = this.store.query('database/role', {});
    // FUTURE find records for the models
    return hash({
      connections: connection,
      roles: role,
    });
  },
});
