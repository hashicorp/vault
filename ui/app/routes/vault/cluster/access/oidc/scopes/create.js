import Route from '@ember/routing/route';

export default class OidcScopesCreateRoute extends Route {
  model() {
    return this.store.createRecord('oidc/scope');
  }
}
