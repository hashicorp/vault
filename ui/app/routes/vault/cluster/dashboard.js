import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { hash } from 'rsvp';

export default class VaultClusterDashboardRoute extends Route {
  @service store;

  model() {
    return hash({
      secretEngines: this.store.query('secret-engine', {}),
    });
  }
}
