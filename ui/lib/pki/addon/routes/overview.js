import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { hash } from 'rsvp';

export default class PkiOverviewRoute extends Route {
  @service secretMountPath;
  @service auth;
  @service store;

  hasConfig() {
    // When the engine is configured, it creates a default issuer.
    // If the issuers list is empty, we know it hasn't been configured
    return (
      this.store
        .query('pki/issuer', { backend: this.secretMountPath.currentPath })
        .then(() => true)
        // this endpoint is unauthenticated, so we're not worried about permissions errors
        .catch(() => false)
    );
  }

  async fetchAllRoles() {
    try {
      return await this.store.query('pki/role', { backend: this.secretMountPath.currentPath });
    } catch (e) {
      return e.httpStatus;
    }
  }

  async fetchAllIssuers() {
    try {
      return await this.store.query('pki/issuer', { backend: this.secretMountPath.currentPath });
    } catch (e) {
      return e.httpStatus;
    }
  }

  async model() {
    return hash({
      hasConfig: this.hasConfig(),
      engine: this.modelFor('application'),
      roles: this.fetchAllRoles(),
      issuers: this.fetchAllIssuers(),
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
  }
}
