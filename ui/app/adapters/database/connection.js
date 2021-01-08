import ApplicationAdapter from '../application';

export default ApplicationAdapter.extend({
  namespace: 'v1',
  pathForType() {
    return 'connection';
  },

  urlFor(backend, id) {
    // example : id mydb
    let url = `${this.buildURL()}/${backend}/config`;
    if (id) {
      url = `${this.buildURL()}/${backend}/config/${id}`;
    }
    console.log('URLLLLL', url);
    return url;
  },

  optionsForQuery(id) {
    let data = {};
    if (!id) {
      data['list'] = true;
    }
    return { data };
  },

  fetchByQuery(store, query) {
    const { backend, id } = query;
    return this.ajax(this.urlFor(backend, id), 'GET', this.optionsForQuery(id)).then(resp => {
      // MONDAY: why is AJax returning something different?
      const data = {
        backend,
      };
      console.log('🇲🇻 RESPONSONDFLJl', resp);
      return resp;
      // return assign({}, resp, data);
    });
  },

  query(store, type, query) {
    return this.fetchByQuery(store, query);
  },

  queryRecord(store, type, query) {
    return this.fetchByQuery(store, query);
  },
});
