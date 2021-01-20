import Route from '@ember/routing/route';
import { hash } from 'rsvp';

// ARG TODO query permissions!

export default Route.extend({
  type: '',
  enginePathParam() {
    let { backend } = this.paramsFor('vault.cluster.secrets.backend');
    return backend;
  },
  model() {
    let backend = this.enginePathParam();
    let queryOptions = { backend, id: '' };
    let secretEngine = this.store.peekRecord('secret-engine', backend);
    let type = secretEngine && secretEngine.get('engineType');
    let connection = this.store.query('database/connection', queryOptions);
    let role = this.store.query('database/role', queryOptions);
    let staticRole = this.store.query('database/static-role', queryOptions);

    return hash({
      connections: connection,
      roles: role,
      staticRoles: staticRole,
      engineType: 'database',
      id: type,
    });
  },
});
