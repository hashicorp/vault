import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
export default class PkiRoleGenerateRoute extends Route {
  @service store;
  @service secretMountPath;
  @service pathHelp;

  beforeModel() {
    // Must call this promise before the model hook otherwise
    // the model doesn't hydrate from OpenAPI correctly.
    return this.pathHelp.getNewModel('pki/certificate/generate', this.secretMountPath.currentPath);
  }

  async model() {
    const { role } = this.paramsFor('roles/role');
    return this.store.createRecord('pki/certificate/generate', {
      role,
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    const { role } = this.paramsFor('roles/role');
    controller.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview' },
      { label: 'roles', route: 'roles.index' },
      { label: role, route: 'roles.role.details' },
      { label: 'generate certificate' },
    ];
    // This is updated on successful generate in the controller
    controller.hasSubmitted = false;
  }
}
