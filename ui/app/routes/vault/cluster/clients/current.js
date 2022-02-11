import Route from '@ember/routing/route';
import RSVP from 'rsvp';
import { action } from '@ember/object';

export default class CurrentRoute extends Route {
  async model() {
    let parentModel = this.modelFor('vault.cluster.clients');

    return RSVP.hash({
      config: parentModel.config,
      monthly: await this.store.queryRecord('clients/monthly', {}),
      versionHistory: parentModel.versionHistory,
    });
  }

  @action
  async loading(transition) {
    // eslint-disable-next-line ember/no-controller-access-in-routes
    let controller = this.controllerFor(this.routeName);
    if (controller) {
      // must use set here or it does not work see docs
      controller.set('currentlyLoading', true);
      transition.promise.finally(function () {
        controller.set('currentlyLoading', false);
      });
    }
  }
}
