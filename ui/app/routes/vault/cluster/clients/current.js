import Route from '@ember/routing/route';
import RSVP from 'rsvp';

export default class CurrentRoute extends Route {
  async model() {
    let parentModel = this.modelFor('vault.cluster.clients');

    return RSVP.hash({
      config: parentModel.config,
      monthly: await this.store.queryRecord('clients/monthly', {}),
      versionHistory: parentModel.versionHistory,
    });
  }
}
