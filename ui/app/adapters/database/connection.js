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
  fetchByQuery(id) {
    // ARG todo pass in id later.
    const backendName = 'database'; // TODO: grab this dynamically
    return this.ajax(this.urlFor(backendName, id), 'GET', this.optionsForQuery(id)).then(resp => {
      // resp.id = id;
      resp.backend = backendName;
      if (id) {
        resp.id = id;
      }
      return resp;
    });
  },
  query() {
    return this.fetchByQuery();
  },

  queryRecord(store, type, query) {
    const { id, backend } = query;
    console.log('querying record for ', backend, id);
    return this.fetchByQuery(id);
  },
});
