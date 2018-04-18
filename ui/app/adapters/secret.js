import Ember from 'ember';
import ApplicationAdapter from './application';
const { computed } = Ember;

export default ApplicationAdapter.extend({
  namespace: 'v1',

  headers: computed(function() {
    return {
      'X-Vault-Kv-Client': 'v2',
    };
  }),

  createOrUpdate(store, type, snapshot) {
    const serializer = store.serializerFor(type.modelName);
    const data = serializer.serialize(snapshot);
    const { id } = snapshot;

    return this.ajax(this.urlForSecret(snapshot.attr('backend'), id), 'POST', {
      data: { data },
    });
  },

  createRecord() {
    return this.createOrUpdate(...arguments);
  },

  updateRecord() {
    return this.createOrUpdate(...arguments);
  },

  deleteRecord(store, type, snapshot) {
    const { id } = snapshot;
    return this.ajax(this.urlForSecret(snapshot.attr('backend'), id), 'DELETE');
  },

  urlForSecret(backend, id, infix = 'data') {
    let url = `${this.buildURL()}/${backend}/${infix}/`;
    if (!Ember.isEmpty(id)) {
      url = url + id;
    }

    return url;
  },

  optionsForQuery(id, action) {
    let data = {};
    if (action === 'query') {
      data['list'] = true;
    }

    return { data };
  },

  urlForQuery(query) {
    let { id, backend } = query;
    return this.urlForSecret(backend, id, 'metadata');
  },

  urlForQueryRecord(query) {
    let { id, backend } = query;
    return this.urlForSecret(backend, id);
  },

  query(store, type, query) {
    return this.ajax(
      this.urlForQuery(query, type.modelName),
      'GET',
      this.optionsForQuery(query.id, 'query')
    ).then(resp => {
      resp.id = query.id;
      return resp;
    });
  },

  queryRecord(store, type, query) {
    return this.ajax(
      this.urlForQueryRecord(query, type.modelName),
      'GET',
      this.optionsForQuery(query.id, 'queryRecord')
    ).then(resp => {
      resp.id = query.id;
      return resp;
    });
  },
});
