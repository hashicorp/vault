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

  createRecord(store, type, snapshot) {
    const serializer = store.serializerFor(type.modelName);
    const data = serializer.serialize(snapshot);
    const id = snapshot.attr('name');
    const backend = snapshot.attr('backend');

    return this.ajax(this.urlFor(backend, id), 'POST', { data }).then(() => {
      // ember data doesn't like 204s if it's not a DELETE
      return {
        data: {
          id,
          ...data,
        },
      };
    });
  },

  deleteRecord(store, type, snapshot) {
    const id = snapshot.id;
    return this.ajax(this.urlFor('database', id), 'DELETE');
  },
});
