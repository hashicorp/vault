import Route from '@ember/routing/route';

export default class MfaLoginEnforcementCreateRoute extends Route {
  setupController(controller) {
    super.setupController(...arguments);
    // if route was refreshed after type select recreate method model
    const { type } = controller;
    if (type) {
      controller.set('method', this.store.createRecord('mfa-method', { type }));
    }
  }
}
