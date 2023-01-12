import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { hash } from 'rsvp';

export default class PkiOverviewRoute extends Route {
  @service secretMountPath;
  @service auth;
  @service store;

  get win() {
    return this.window || window;
  }

  async hasConfig() {
    // When the engine is configured, it creates a default issuer.
    // If the issuers list is empty, we know it hasn't been configured
    const endpoint = `${this.win.origin}/v1/${this.secretMountPath.currentPath}/issuers?list=true`;

    try {
      await this.auth.ajax(endpoint, 'GET', {});
      return true;
    } catch (e) {
      return false;
    }
  }

  async fetchEngine() {
    const model = await this.store.query('secret-engine', {
      path: this.secretMountPath.currentPath,
    });
    return model.get('firstObject');
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

  async fetchAllCertificates() {
    try {
      return await this.store.query('pki/certificate', { backend: this.secretMountPath.currentPath });
    } catch (e) {
      return e.httpStatus;
    }
  }

  async model() {
    return hash({
      hasConfig: this.hasConfig(),
      engine: this.fetchEngine(),
      roles: this.fetchAllRoles(),
      issuers: this.fetchAllIssuers(),
      certificates: this.fetchAllCertificates(),
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    const backend = this.secretMountPath.currentPath || 'pki';

    controller.backend = backend;
    controller.rolesCanList = !(resolvedModel.roles === 403);
    controller.issuersCanList = !(resolvedModel.issuers === 403);
    controller.breadcrumbs = [{ label: 'secrets', route: 'secrets', linkExternal: true }, { label: backend }];
  }
}
