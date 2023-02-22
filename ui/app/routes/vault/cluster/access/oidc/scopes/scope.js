import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class OidcScopeRoute extends Route {
  @service store;

  model({ name }) {
    return this.store.findRecord('oidc/scope', name);
  }
}
