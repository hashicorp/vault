import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class KmipRoleRoute extends Route {
  @service store;
  @service secretMountPath;
  @service pathHelp;

  async beforeModel() {
    return await this.pathHelp.getNewModel('kmip/role', this.secretMountPath.currentPath);
  }

  async model(params) {
    return await this.store.queryRecord('kmip/role', {
      backend: this.secretMountPath.currentPath,
      scope: params.scope_name,
      id: params.role_name,
    });
  }

  setupController(controller, model) {
    let { scope_name: scope, role_name: role } = this.paramsFor('role');
    controller.setProperties({ role, scope });
    controller.set('model', model);
  }
}
