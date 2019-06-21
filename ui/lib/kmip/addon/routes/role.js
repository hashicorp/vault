import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default Route.extend({
  store: service(),
  secretMountPath: service(),
  pathHelp: service(),

  beforeModel() {
    return this.pathHelp.getNewModel('kmip/role', this.secretMountPath.currentPath);
  },

  model(params) {
    return this.store.queryRecord('kmip/role', {
      backend: this.secretMountPath.currentPath,
      scope: params.scope_name,
      id: params.role_name,
    });
  },

  setupController(controller) {
    this._super(...arguments);
    let { scope_name: scope, role_name: role } = this.paramsFor('role');
    controller.setProperties({ role, scope });
  },
});
