import Route from '@ember/routing/route';

export default class OidcScopesRoute extends Route {
  model() {
    return this.store.query('oidc/scope', {}).catch((err) => {
      if (err.httpStatus === 404) {
        return [];
      } else {
        throw err;
      }
    });
  }
}
