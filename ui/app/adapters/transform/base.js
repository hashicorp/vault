import { assign } from '@ember/polyfills';
import { resolve, allSettled } from 'rsvp';
import ApplicationAdapter from '../application';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default ApplicationAdapter.extend({
  // TODO this adapter was copied over, much of this stuff may or may not need to be here.
  namespace: 'v1',

  pathForType(type) {
    return type.replace('transform/', '');
  },

  _url(modelType, backend, id) {
    let type = this.pathForType(modelType);
    let base = `/${this.namespace}/${encodePath(backend)}/${type}`;
    if (id) {
      return `${base}/${encodePath(id)}`;
    }
    // CBS TODO: if no id provided, should we assume it's a LIST?
    return base;
  },

  createOrUpdate(store, type, snapshot) {
    const serializer = store.serializerFor('transform'); // TODO replace transform with type.modelName
    const data = serializer.serialize(snapshot);
    const { id } = snapshot;
    let url = this._url(type.modelName, snapshot.record.get('backend'), id);

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
    return this.ajax(this._url(type.modelName, snapshot.record.get('backend'), id), 'DELETE');
  },

  urlForQuery(modelType, backend) {
    let base = this._url(modelType, backend);
    return base + '?list=true';
  },

  query(store, type, query) {
    return this.ajax(this.urlForQuery(type.modelName, query.backend), 'GET').then(result => {
      console.log(result);

      return result;
    });
  },

  queryRecord(store, type, query) {
    return this.ajax(this._url(type.modelName, query.backend, query.id), 'GET').then(result => {
      console.log(result);

      return result;
    });
  },
});
