import Route from '@ember/routing/route';

export default class OidcKeysCreateRoute extends Route {
  model() {
    return this.store.createRecord('oidc/key');
  }
}
