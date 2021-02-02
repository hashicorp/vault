import Route from '@ember/routing/route';
import { hash } from 'rsvp';

export default Route.extend({
  type: '',
  enginePathParam() {
    let { backend } = this.paramsFor('vault.cluster.secrets.backend');
    return backend;
  },
  async fetchConnection(queryOptions) {
    try {
      return await this.store.query('database/connection', queryOptions);
    } catch (e) {
      return e.httpStatus;
    }
  },
  async fetchAllRoles(queryOptions) {
    try {
      return await this.store.query('database/role', queryOptions);
    } catch (e) {
      return e.httpStatus;
    }
  },
  model() {
    let backend = this.enginePathParam();
    let queryOptions = { backend, id: '' };
    let secretEngine = this.store.peekRecord('secret-engine', backend);
    let type = secretEngine && secretEngine.get('engineType');

    let connection = this.fetchConnection(queryOptions);
    let role = this.fetchAllRoles(queryOptions);

    return hash({
      backend,
      connections: connection,
      roles: role,
      engineType: 'database',
      id: type,
    });
  },
});
