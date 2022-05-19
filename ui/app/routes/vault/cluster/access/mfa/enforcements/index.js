import Route from '@ember/routing/route';

export default class MfaEnforcementsRoute extends Route {
  model() {
    return this.store.query('mfa-login-enforcement', {}).catch((err) => {
      if (err.httpStatus === 404) {
        return [];
      } else {
        throw err;
      }
    });
  }
  setupController(controller, model) {
    controller.set('model', model);
  }
}
