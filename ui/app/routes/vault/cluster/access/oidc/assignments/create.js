import Route from '@ember/routing/route';

export default class OidcAssignmentsCreateRoute extends Route {
  model() {
    return this.store.createRecord('oidc/assignment');
  }
}
