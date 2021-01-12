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
  fetchByQuery() {
    // ARG todo pass in id later.
    const backendName = 'database'; // TODO: grab this dynamically
    return this.ajax(this.urlFor(backendName), 'GET', this.optionsForQuery()).then(resp => {
      // resp.id = id;
      resp.backend = backendName;
      return resp;
    });
  },
  query() {
    return this.fetchByQuery();
  },
});
