import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class OidcAssignmentsRoute extends Route {
  @service store;
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
