import Route from '@ember/routing/route';

export default class OidcProviderClientsRoute extends Route {
  async model() {
    const model = this.modelFor('vault.cluster.access.oidc.providers.provider');
    if (model.allowedClientIds.includes('*')) {
      return this.store.query('oidc/client', {}).catch((err) => {
        if (err.httpStatus === 404) {
          return [];
        } else {
          throw err;
        }
      });
    }
    return await model.allowedClientIds.map((client) => {
      return this.store.findRecord('oidc/client', client);
    });
  }
}
