import Route from '@ember/routing/route';
import { action } from '@ember/object';

export default class ConfigRoute extends Route {
  model() {
    return this.store.queryRecord('clients/config', {});
  }
  @action
  async loading(transition) {
    // eslint-disable-next-line ember/no-controller-access-in-routes
    let controller = this.controllerFor('vault.cluster.clients.config');
    controller.set('currentlyLoading', true);
    transition.promise.finally(function () {
      controller.set('currentlyLoading', false);
    });
  }
}
