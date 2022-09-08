import Route from '@ember/routing/route';

export default class OidcClientProvidersRoute extends Route {
  model() {
    const model = this.modelFor('vault.cluster.access.oidc.clients.client');
    return this.store
      .query('oidc/provider', {
        allowed_client_id: model.clientId,
      })
      .catch((err) => {
        if (err.httpStatus === 404) {
          return [];
        } else {
          throw err;
        }
      });
  }
}
