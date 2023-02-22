import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class OidcProviderRoute extends Route {
  @service store;

  model({ name }) {
    return this.store.findRecord('oidc/provider', name);
  }
}
