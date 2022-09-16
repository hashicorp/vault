import Route from '@ember/routing/route';

export default class OidcClientsCreateRoute extends Route {
  model() {
    return this.store.createRecord('oidc/client');
  }
}
