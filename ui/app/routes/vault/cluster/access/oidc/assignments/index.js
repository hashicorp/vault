import Route from '@ember/routing/route';

export default class OidcAssignmentsRoute extends Route {
  model() {
    return this.store.query('oidc/assignment', {}).catch((err) => {
      if (err.httpStatus === 404) {
        return [];
      } else {
        throw err;
      }
    });
  }
}
