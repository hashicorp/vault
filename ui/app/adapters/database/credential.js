import ApplicationAdapter from '../application';

export default ApplicationAdapter.extend({
  namespace: 'v1',

  fetchByQuery(store, query) {
    const { backend, secret } = query;
    return this.ajax(`${this.buildURL()}/${backend}/creds/${secret}`, 'GET').then(resp => {
      return resp;
    });
  },
  queryRecord(store, type, query) {
    return this.fetchByQuery(store, query);
  },
});
