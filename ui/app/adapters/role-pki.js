import ApplicationAdapter from './application';
import Ember from 'ember';

export default ApplicationAdapter.extend({
  namespace: 'v1',

  createOrUpdate(store, type, snapshot, requestType) {
    const serializer = store.serializerFor(type.modelName);
    const data = serializer.serialize(snapshot, requestType);
    const { id } = snapshot;
    let url = this.urlForRole(snapshot.record.get('backend'), id);

    return this.ajax(url, 'POST', { data });
  },

  createRecord() {
    return this.createOrUpdate(...arguments);
  },

  updateRecord() {
    return this.createOrUpdate(...arguments, 'update');
  },

  deleteRecord(store, type, snapshot) {
    const { id } = snapshot;
    return this.ajax(this.urlForRole(snapshot.record.get('backend'), id), 'DELETE');
  },

  pathForType() {
    return 'roles';
  },

  urlForRole(backend, id) {
    let url = `${this.buildURL()}/${backend}/roles`;
    if (id) {
      url = url + '/' + id;
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
    const { id, backend } = query;
    return this.ajax(this.urlForRole(backend, id), 'GET', this.optionsForQuery(id)).then(resp => {
      const data = {
        id,
        name: id,
        backend,
      };

      return Ember.assign({}, resp, data);
    });
  },

  query(store, type, query) {
    return this.fetchByQuery(store, query);
  },

  queryRecord(store, type, query) {
    return this.fetchByQuery(store, query);
  },
});
