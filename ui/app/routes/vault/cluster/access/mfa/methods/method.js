import Route from '@ember/routing/route';

export default class MfaMethodRoute extends Route {
  model({ id }) {
    return this.store.findRecord('mfa-method', id);
  }
  setupController(controller, model) {
    controller.set('model', model);
  }
}
