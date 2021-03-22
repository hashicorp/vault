import Route from '@ember/routing/route';
import fetch from 'fetch';

export default class VaultClusterSecretsCustomRoute extends Route {
  async model(params) {
    // params.backend
    const resp = await fetch(`/v1/secret/?help=1`, {
      headers: {
        Accept: 'application/json',
        'Content-Type': 'application/json',
        'X-Vault-Token': 'root',
      },
    });
    const json = await resp.json();
    const tabs = json.openapi?.paths || {};
    let tabsArray = [];
    Object.keys(tabs).forEach(key => {
      if (tabs[key]?.get) {
        tabsArray.push(key);
      }
    });
    return {
      tabs: tabsArray,
    };
  }
}
