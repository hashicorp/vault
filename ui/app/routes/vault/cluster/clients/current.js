import Route from '@ember/routing/route';
import RSVP from 'rsvp';
import { inject as service } from '@ember/service';

export default class CurrentRoute extends Route {
  @service store;

  async model() {
    const parentModel = this.modelFor('vault.cluster.clients');

    return RSVP.hash({
      config: parentModel.config,
      monthly: await this.store.queryRecord('clients/monthly', {}),
      versionHistory: parentModel.versionHistory,
    });
  }
}
