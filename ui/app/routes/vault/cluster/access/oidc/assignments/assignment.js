import Route from '@ember/routing/route';

export default class OidcAssignmentRoute extends Route {
  model({ name }) {
    return this.store.findRecord('oidc/assignment', name);
  }
}
