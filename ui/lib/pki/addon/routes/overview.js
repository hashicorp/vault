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

  async fetchAllRolesCapabilities() {
    const query = { id: `${this.secretMountPath.currentPath}/roles` };
    return await this.store.queryRecord('capabilities', query);
  }

  async fetchAllIssuersCapabilities() {
    const query = { id: `${this.secretMountPath.currentPath}/issuers` };
    return await this.store.queryRecord('capabilities', query);
  }
  async fetchAllCertificatesCapabilities() {
    const query = { id: `${this.secretMountPath.currentPath}/certificates` };
    return await this.store.queryRecord('capabilities', query);
  }

  async model() {
    return hash({
      hasConfig: this.hasConfig(),
      engine: this.fetchEngine(),
      roles: this.fetchAllRoles(),
      issuers: this.fetchAllIssuers(),
      certificates: this.fetchAllCertificates(),
      rolesCapabilities: this.fetchAllRolesCapabilities(),
      issuersCapabilities: this.fetchAllIssuersCapabilities(),
      certificateCapabilities: this.fetchAllCertificatesCapabilities(),
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    const backend = this.secretMountPath.currentPath || 'pki';

    const { rolesCapabilities, roles, certificates, hasConfig } = resolvedModel;

    if (rolesCapabilities.canList && roles.length) {
      controller.roleOptions = roles.map((role) => {
        return { name: role.id, id: role.id };
      });
    }

    if (hasConfig && certificates.length) {
      controller.certificateOptions = certificates.map((certificate) => {
        return { name: certificate.id, id: certificate.id };
      });
    }

    controller.breadcrumbs = [{ label: 'secrets', route: 'secrets', linkExternal: true }, { label: backend }];
  }
}
