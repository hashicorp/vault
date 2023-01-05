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

  hasConfig() {
    const endpoint = `${this.win.origin}/v1/${this.secretMountPath.currentPath}/issuers?list=true`;
    return this.auth
      .ajax(endpoint, 'GET', {})
      .then(() => true)
      .catch(() => false);
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

  fetchEngine() {
    return this.store
      .query('secret-engine', {
        path: this.secretMountPath.currentPath,
      })
      .then((model) => {
        if (model) {
          return model.get('firstObject');
        }
      });
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

    controller.canViewRoles = resolvedModel.roles.length;
    controller.canViewIssuers = resolvedModel.issuers.length;

    controller.roles = resolvedModel.roles.map((role) => {
      return { name: role.id, id: role.id };
    });
    controller.issuers = resolvedModel.issuers.map((issuer) => {
      return { name: issuer.id, id: issuer.id };
    });
    controller.certificates = resolvedModel.certificates.map((certificate) => {
      return { name: certificate.id, id: certificate.id };
    });
    controller.breadcrumbs = [{ label: 'secrets', route: 'secrets', linkExternal: true }, { label: backend }];
  }
}
