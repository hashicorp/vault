import Route from '@ember/routing/route';
import ClusterRoute from 'vault/mixins/cluster-route';
import { hash } from 'rsvp';

// TODO query permissions!

export default Route.extend(ClusterRoute, {
  model() {
    // let { backend } = this.paramsFor('vault.cluster.secrets.backend');

    let connection = this.store.query('database/connection', {});
    console.log(connection.get('length'), 'LENGHT');
    // FUTURE find records for the models
    return hash({
      connections: connection,
      roles: [],
    });
  },
});
