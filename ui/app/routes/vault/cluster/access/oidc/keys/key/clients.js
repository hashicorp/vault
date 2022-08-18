import Route from '@ember/routing/route';

export default class OidcKeyClientsRoute extends Route {
  async model() {
    const { allowedClientIds } = this.modelFor('vault.cluster.access.oidc.keys.key');
    return await this.store.query('oidc/client', { paramKey: 'client_id', filterFor: allowedClientIds });
  }
}
