import ApplicationAdapter from '../application';

export default class OidcKeyAdapter extends ApplicationAdapter {
  namespace = 'v1';

  pathForType() {
    // backend name prepended in buildURL method
    return 'identity/oidc/key';
  }

  query(store, type, query) {
    const url = this.urlForQuery(query, type.modelName);
    return this.ajax(url, 'GET', { data: { list: true } });
  }
}
