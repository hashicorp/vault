import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class VaultClusterPluginsRoute extends Route {
  @service auth;
  async model() {
    // TODO:
    // get token from localstorage,
    // wrap it (endpoint),
    // pass to component so it can provide it in the url
    const response = await this.auth.ajax('/v1/plugins', 'GET', {});
    return {
      ...response.data,
      wrappedToken: 's.qMktRJnvcosL9T8EPApDqsHL',
    };
  }
}
