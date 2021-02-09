import ApplicationAdapter from '../application';

export default ApplicationAdapter.extend({
  namespace: 'v1',

  urlFor(backend, id, type = 'REST') {
    if (type === 'ROTATE') {
      return `${this.buildURL()}/${backend}/rotate-root/${id}`;
    } else if (type === 'RESET') {
      return `${this.buildURL()}/${backend}/reset/${id}`;
    }
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
    return this.fetchByQuery(store, query);
  },

  queryRecord(store, type, query) {
    return this.fetchByQuery(store, query);
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

  updateRecord() {
    return this.createRecord(...arguments);
  },

  deleteRecord(store, type, snapshot) {
    const id = snapshot.id;
    return this.ajax(this.urlFor('database', id), 'DELETE');
  },

  rotateRootCredentials(backend, id) {
    return this.ajax(this.urlFor('database', id, 'ROTATE'), 'POST');
  },

  resetConnection(backend, id) {
    return this.ajax(this.urlFor(backend, id, 'RESET'), 'POST');
  },
});
