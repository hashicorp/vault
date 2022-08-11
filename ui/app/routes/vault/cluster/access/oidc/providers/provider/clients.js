import Route from '@ember/routing/route';

export default class OidcProviderClientsRoute extends Route {
  async model() {
    const providerModel = this.modelFor('vault.cluster.access.oidc.providers.provider');
    const { allowedClientIds } = providerModel;
    return await this.store
      .query('oidc/client', {})
      .then((clientRecords) => {
        // return all clientRecords if glob is in array of IDs
        if (allowedClientIds.includes('*')) {
          return clientRecords;
        }
        return clientRecords.filter((client) => allowedClientIds.includes(client.id));
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
