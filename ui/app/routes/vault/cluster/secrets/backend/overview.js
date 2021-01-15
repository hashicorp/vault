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
    let secretEngine = this.store.peekRecord('secret-engine', backend);
    let type = secretEngine && secretEngine.get('engineType');

    let connection = this.store.query('database/connection', {});
    let role = this.store.query('database/role', {});
    let staticRole = this.store.query('database/static-role', {});

    return hash({
      connections: connection,
      roles: role,
      staticRoles: staticRole,
      engineType: 'database',
      id: type,
    });
  },
});
