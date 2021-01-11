import ApplicationAdapter from '../application';

export default ApplicationAdapter.extend({
  namespace: 'v1',
  // ARG unsure if anything pathForType is doing
  pathForType() {
    return 'connection';
  },
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
    return this.ajax(this.urlFor('database'), 'GET', this.optionsForQuery()).then(resp => {
      // resp.id = id;
      resp.backend = 'database';
      return resp;
    });
  },
  query() {
    return this.fetchByQuery();
  },
});
