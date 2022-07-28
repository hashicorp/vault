import Route from '@ember/routing/route';

export default class OidcClientRoute extends Route {
  model({ name }) {
    return this.store.findRecord('oidc/client', name);
  }
}
