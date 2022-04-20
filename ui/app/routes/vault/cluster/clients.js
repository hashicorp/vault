import Route from '@ember/routing/route';
import RSVP from 'rsvp';
import { action } from '@ember/object';
import getStorage from 'vault/lib/token-storage';

const INPUTTED_START_DATE = 'vault:ui-inputted-start-date';
export default class ClientsRoute extends Route {
  async getVersionHistory() {
    try {
      let arrayOfModels = [];
      let response = await this.store.findAll('clients/version-history'); // returns a class with nested models
      response.forEach((model) => {
        arrayOfModels.push({
          id: model.id,
          previousVersion: model.previousVersion,
          timestampInstalled: model.timestampInstalled,
        });
      });
      return arrayOfModels;
    } catch (e) {
      console.debug(e);
      return [];
    }
  }

  async model() {
    let config = await this.store.queryRecord('clients/config', {}).catch((e) => {
      console.debug(e);
      // swallowing error so activity can show if no config permissions
      return {};
    });

    return RSVP.hash({
      config,
      versionHistory: this.getVersionHistory(),
    });
  }

  @action
  async loading(transition) {
    // eslint-disable-next-line ember/no-controller-access-in-routes
    let controller = this.controllerFor(this.routeName);
    controller.set('currentlyLoading', true);
    transition.promise.finally(function () {
      controller.set('currentlyLoading', false);
    });
  }

  @action
  deactivate() {
    // when navigating away from parent route, delete manually inputted license start date
    getStorage().removeItem(INPUTTED_START_DATE);
  }
}
