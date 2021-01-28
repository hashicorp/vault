import ApplicationAdapter from '../application';

export default ApplicationAdapter.extend({
  namespace: 'v1',

  urlFor(backend, id) {
    let url = `${this.buildURL()}/${backend}/roles`;
    if (id) {
      url = `${this.buildURL()}/${backend}/roles/${id}`;
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
    // ARG todo pass in id later.
    const { backend, id } = query;
    return this.ajax(this.urlFor(backend, id), 'GET', this.optionsForQuery(id)).then(resp => {
      // resp.id = id;
      resp.backend = backend;
      return resp;
    });
  },
  query(store, type, query) {
    return this.fetchByQuery(store, query);
  },
});
