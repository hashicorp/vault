import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default Route.extend({
  store: service(),
  model(params) {
    let model = { id: 'serial-beep-boop', get: true, activate: true };
    return model;
    return this.store.findRecord('kmip/role');
  },

  setupController(controller) {
    this._super(...arguments);
    let { scope_name: scope, role_name: role } = this.paramsFor('role');
    controller.setProperties({ role, scope });
  },
});
