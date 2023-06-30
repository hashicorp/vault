import ApplicationAdapter from './application';

export default class SecretListAdapter extends ApplicationAdapter {
  namespace = 'v1';

  urlForQuery(query, modelName) {
    const { backend } = query;
    return this.buildURL() + `/${backend}/keys`;
  }
  queryRecord(store, type, query) {
    const url = this.urlForQuery(query, type.modelName);
    const queryParams = { list: true };
    return this.ajax(url, 'GET', { data: queryParams }).then((resp) => {
      return {
        id: `${query.modelType}-list-${query.backend}`,
        backend: query.backend,
        ...resp,
      };
    });
  }
}
