import Route from '@ember/routing/route';

export default class MfaMethodRoute extends Route {
  model(params) {
    return this.store.findRecord('mfa-method', params['method_id']);
  }
  setupController(controller, model) {
    controller.set('model', model);
  }
}
