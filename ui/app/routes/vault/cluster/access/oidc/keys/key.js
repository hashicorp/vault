import Route from '@ember/routing/route';

export default class OidcKeyRoute extends Route {
  model({ name }) {
    return this.store.findRecord('oidc/key', name);
  }
}
