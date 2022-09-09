import Route from '@ember/routing/route';

export default class OidcProvidersCreateRoute extends Route {
  model() {
    return this.store.createRecord('oidc/provider');
  }
}
