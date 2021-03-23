import Route from '@ember/routing/route';
import fetch from 'fetch';

export default class VaultClusterSecretsCustomRoute extends Route {
  async model(params) {
    // params.backend
    const { backend } = params;
    const resp = await fetch(`/v1/${backend}/?help=1`, {
      headers: {
        Accept: 'application/json',
        'Content-Type': 'application/json',
        'X-Vault-Token': 'root',
      },
    });
    const json = await resp.json();
    const tabs = json.openapi?.paths || {};
    let info = {};
    Object.keys(tabs).forEach(key => {
      // only add items that are gettable
      if (tabs[key]?.get) {
        const openApi = tabs[key];
        const tabInfo = {
          key: key,
          description: openApi.description,
          gettable: !!openApi.get,
          creatable: !!openApi['x-vault-createSupported'],
          urlPath: key,
          ...openApi,
        };
        const getParams = tabs[key].get.parameters?.find(p => p.name === 'list');
        if (getParams) {
          tabInfo.listable = true;
        }
        // tabs[key].key = key;
        // todo clean tab name and add to info
        info[encodeURIComponent(key)] = tabInfo;
      }
    });
    return {
      backend,
      info,
      paramText: '',
    };
  }
}
