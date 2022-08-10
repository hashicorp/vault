import Route from '@ember/routing/route';

export default class OidcProviderRoute extends Route {
  model({ name }) {
    return this.store.findRecord('oidc/provider', name);
  }
}
