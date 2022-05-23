import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class MfaMethodsRoute extends Route {
  @service router;

  model() {
    return this.store.query('mfa-method', {}).catch((err) => {
      if (err.httpStatus === 404) {
        return [];
      } else {
        throw err;
      }
    });
  }

  afterModel(model) {
    if (model.length === 0) {
      this.router.transitionTo('vault.cluster.access.mfa');
    }
  }
  setupController(controller, model) {
    controller.set('model', model);
  }
}
