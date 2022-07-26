import Route from '@ember/routing/route';

export default class OidcProvidersRoute extends Route {
  model() {
    return this.store.query('oidc/provider', {}).catch((err) => {
      if (err.httpStatus === 404) {
        return [];
      } else {
        throw err;
      }
    });
  }
}
