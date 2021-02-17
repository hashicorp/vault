import ApplicationAdapter from '../application';

export default ApplicationAdapter.extend({
  namespace: 'v1',

  fetchByQuery(store, query) {
    const { backend, roleType, secret } = query;
    let creds = roleType === 'static' ? 'static-creds' : 'creds';
    return this.ajax(
      `${this.buildURL()}/${encodeURIComponent(backend)}/${creds}/${encodeURIComponent(secret)}`,
      'GET'
    );
  },
  queryRecord(store, type, query) {
    return this.fetchByQuery(store, query);
  },
});
