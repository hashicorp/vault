import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
// import getStorage from 'vault/lib/token-storage';

export default class VaultClusterPluginShowRoute extends Route {
  @service auth;

  beforeModel() {
    window.localStorage.setItem('example', 'foo-bar');
  }
  async model(params) {
    console.log('vault session', window.sessionStorage);
    console.log('vault local', window.localStorage);
    const response = await this.auth.ajax(`/v1/plugin/${params.plugin}`, 'GET', {});
    return response.data;
  }
}
