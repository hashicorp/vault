import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
export default class MfaRoute extends Route {
  @service router;

  model() {
    return this.store.findAll('mfa-method').then((data) => {
      return data;
    });
  }
  afterModel(model) {
    if (model.get('length') === 0) {
      this.router.transitionTo('vault.cluster.access.mfa.configure');
    }
  }
  setupController(controller, model) {
    controller.set('model', model);
  }
}
