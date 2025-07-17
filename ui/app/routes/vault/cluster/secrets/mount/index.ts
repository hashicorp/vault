import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { PluginsCatalogListPluginsWithTypeListEnum } from '@hashicorp/vault-client-typescript';

import type ApiService from 'vault/services/api';

// sample response from plugin/catalog/secret
// {
//   "data": {
//       "keys": [
//           "aws",
//           "azure",
//           "custom-auth-plugin",
//           "gcp",
//           "ldap"
//       ]
//   }
// }

export default class VaultClusterSecretsMountIndexRoute extends Route {
  @service declare readonly api: ApiService;

  async model() {
    const { keys } = await this.api.sys.pluginsCatalogListPluginsWithType(
      'secret',
      PluginsCatalogListPluginsWithTypeListEnum.TRUE
    );
    return keys;
  }
}
