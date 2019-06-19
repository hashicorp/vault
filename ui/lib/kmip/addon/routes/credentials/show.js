import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default Route.extend({
  store: service(),
  secretMountPath: service(),
  credParams() {
    let { role_name: role, scope_name: scope } = this.paramsFor('credentials');
    return {
      role,
      scope,
    };
  },
  model(params) {
    let { role, scope } = this.credParams();
    return this.store.queryRecord('kmip/credential', {
      role,
      scope,
      backend: this.secretMountPath.currentPath,
      id: params.serial,
    });
  },

  setupController(controller) {
    let { role, scope } = this.credParams();
    this._super(...arguments);
    controller.setProperties({ role, scope });
  },
});
