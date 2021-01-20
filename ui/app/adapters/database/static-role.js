import ApplicationAdapter from '../application';

export default ApplicationAdapter.extend({
  namespace: 'v1',

  urlFor(backend, id) {
    let url = `${this.buildURL()}/${backend}/static-roles`;
    if (id) {
      url = `${this.buildURL()}/${backend}/static-roles/${id}`;
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
    // ARG TODO pass in id later
    return this.ajax(this.urlFor(backend), 'GET', this.optionsForQuery()).then(resp => {
      // resp.id = id;
      resp.backend = backend;
      return resp;
    });
  },
  query(store, type, query) {
    return this.fetchByQuery(store, query);
  },
});
