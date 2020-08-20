import ApplicationAdapter from '../application';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default ApplicationAdapter.extend({
  namespace: 'v1',

  pathForType(type) {
    return type.replace('transform/', '');
  },

  createOrUpdate(store, type, snapshot) {
    const { modelName } = type;
    const serializer = store.serializerFor(modelName);
    const data = serializer.serialize(snapshot);
    const { id } = snapshot;
    let url = this.url(snapshot.record.get('backend'), modelName, id);

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
    return this.ajax(this.url(snapshot.record.get('backend'), type.modelName, id), 'DELETE');
  },

  url(backend, modelType, id) {
    let type = this.pathForType(modelType);
    let url = `/${this.namespace}/${encodePath(backend)}/${encodePath(type)}`;
    if (id) {
      return `${url}/${encodePath(id)}`;
    }
    return url + '?list=true';
  },

  fetchByQuery(query) {
    const { backend, modelName, id } = query;
    return this.ajax(this.url(backend, modelName, id), 'GET').then(resp => {
      if (id) {
        return {
          id,
          ...resp,
        };
      }
      return resp;
    });
  },

  query(store, type, query) {
    console.log({ type });
    return this.fetchByQuery(query);
  },

  queryRecord(store, type, query) {
    console.log({ type });
    return this.ajax(this.url(query.backend, type.modelName, query.id), 'GET').then(result => {
      // console.log(result);

      return result;
    });
  },

  // buildUrl(modelName, id, snapshot, requestType, query, returns) {},
});
