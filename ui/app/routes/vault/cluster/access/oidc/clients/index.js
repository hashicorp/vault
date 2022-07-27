import Route from '@ember/routing/route';

export default class OidcClientsRoute extends Route {
  model() {
    return this.store.query('oidc/client', {}).catch((err) => {
      if (err.httpStatus === 404) {
        return [];
      } else {
        throw err;
      }
    });
  }
}
