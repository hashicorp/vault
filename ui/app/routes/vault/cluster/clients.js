import Route from '@ember/routing/route';
import RSVP from 'rsvp';
import { action } from '@ember/object';
import getStorage from 'vault/lib/token-storage';
import { inject as service } from '@ember/service';
const INPUTTED_START_DATE = 'vault:ui-inputted-start-date';

export default class ClientsRoute extends Route {
  @service store;
  async getVersionHistory() {
    try {
      const arrayOfModels = [];
      const response = await this.store.findAll('clients/version-history'); // returns a class with nested models
      response.forEach((model) => {
        arrayOfModels.push({
          id: model.id,
          previousVersion: model.previousVersion,
          timestampInstalled: model.timestampInstalled,
        });
      });
      return arrayOfModels;
    } catch (e) {
      console.debug(e); // eslint-disable-line
      return [];
    }
  }

  async model() {
    const config = await this.store.queryRecord('clients/config', {}).catch((e) => {
      console.debug(e); // eslint-disable-line
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
    const controller = this.controllerFor(this.routeName);
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
