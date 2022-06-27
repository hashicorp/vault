import Route from '@ember/routing/route';

export default class OidcAssignmentRoute extends Route {
  model(params) {
    return this.store.findRecord('oidc/assignment', params.name);
  }
}
