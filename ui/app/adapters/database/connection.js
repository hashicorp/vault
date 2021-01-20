import ApplicationAdapter from '../application';

export default ApplicationAdapter.extend({
  namespace: 'v1',

  urlFor(backend, id) {
    let url = `${this.buildURL()}/${backend}/config`;
    if (id) {
      url = `${this.buildURL()}/${backend}/config/${id}`;
    }
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
      // resp.id = id;
      resp.backend = backend;
      if (id) {
        resp.id = id;
      }
      return resp;
    });
  },
  query(store, type, query) {
    // ARG TODO pass in id
    return this.fetchByQuery(store, query);
  },

  queryRecord(store, type, query) {
    // ARG TODO unsure if using??
    return this.fetchByQuery(store, query);
  },
});
