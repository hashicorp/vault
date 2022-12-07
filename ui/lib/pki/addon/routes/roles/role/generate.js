import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
export default class PkiRoleGenerateRoute extends Route {
  @service store;
  @service secretMountPath;

  async model() {
    const { role } = this.paramsFor('roles/role');
    // const adapter = this.store.adapterFor('pki/role');
    return {
      role,
      backend: this.secretMountPath.currentPath,
    };
    // return adapter.generateCertificate(this.secretMountPath.currentPath, role);
  }
  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    const { role } = this.paramsFor('roles/role');
    const backend = this.secretMountPath.currentPath || 'pki';
    controller.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: backend, route: 'overview' },
      { label: 'roles', route: 'roles.index' },
      { label: role, route: 'roles.role.details' },
      { label: 'generate certificate' },
    ];
    // This is updated on successful generate in the controller
    controller.hasSubmitted = false;
  }
}
