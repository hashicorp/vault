import Route from '@ember/routing/route';

export default class OidcProviderClientsRoute extends Route {
  async model() {
    const { allowedClientIds } = this.modelFor('vault.cluster.access.oidc.providers.provider');
    return await this.store.query('oidc/client', { paramKey: 'client_id', filterFor: allowedClientIds });
  }
}
