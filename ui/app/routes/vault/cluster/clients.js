import Route from '@ember/routing/route';
import RSVP from 'rsvp';

export default class ClientsRoute extends Route {
  async getVersionHistory() {
    try {
      let arrayOfModels = [];
      let response = await this.store.findAll('clients/version-history'); // returns a class with nested models
      response.forEach((model) => {
        arrayOfModels.push({
          id: model.id,
          perviousVersion: model.previousVersion,
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
      // monthly: await this.store.queryRecord('clients/monthly', {}),
      versionHistory: this.getVersionHistory(),
    });
  }
}
