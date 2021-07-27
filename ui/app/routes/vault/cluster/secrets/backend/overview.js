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
  pathQuery(backend, endpoint) {
    return {
      id: `${backend}/${endpoint}/`,
    };
  },
  async fetchCapabilitiesRole(queryOptions) {
    return this.store.queryRecord('capabilities', this.pathQuery(queryOptions.backend, 'roles'));
  },
  async fetchCapabilitiesStaticRole(queryOptions) {
    return this.store.queryRecord('capabilities', this.pathQuery(queryOptions.backend, 'static-roles'));
  },
  async fetchCapabilitiesConnection(queryOptions) {
    return this.store.queryRecord('capabilities', this.pathQuery(queryOptions.backend, 'config'));
  },
  model() {
    let backend = this.enginePathParam();
    let queryOptions = { backend, id: '' };

    let connection = this.fetchConnection(queryOptions);
    let role = this.fetchAllRoles(queryOptions);
    let roleCapabilities = this.fetchCapabilitiesRole(queryOptions);
    let staticRoleCapabilities = this.fetchCapabilitiesStaticRole(queryOptions);
    let connectionCapabilities = this.fetchCapabilitiesConnection(queryOptions);

    return hash({
      backend,
      connections: connection,
      roles: role,
      engineType: 'database',
      id: backend,
      roleCapabilities,
      staticRoleCapabilities,
      connectionCapabilities,
    });
  },
  setupController(controller, model) {
    this._super(...arguments);
    let showEmptyState = model.connections === 404 && model.roles === 404;
    let noConnectionCapabilities =
      !model.connectionCapabilities.canList &&
      !model.connectionCapabilities.canCreate &&
      !model.connectionCapabilities.canUpdate;

    let emptyStateMessage = function() {
      if (noConnectionCapabilities) {
        return 'You cannot yet generate credentials.  Ask your administrator if you think you should have access.';
      } else {
        return 'You can connect and external database to Vault.  We recommend that you create a user for Vault rather than using the database root user.';
      }
    };
    controller.set('showEmptyState', showEmptyState);
    controller.set('emptyStateMessage', emptyStateMessage());
  },
});
