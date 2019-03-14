import { isEmpty } from '@ember/utils';
import ApplicationAdapter from './application';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default ApplicationAdapter.extend({
  namespace: 'v1',

  createOrUpdate(store, type, snapshot) {
    const serializer = store.serializerFor(type.modelName);
    const data = serializer.serialize(snapshot);
    const { id } = snapshot;

    return this.ajax(this.urlForSecret(snapshot.attr('backend'), id), 'POST', { data });
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

  urlForSecret(backend, id) {
    let url = `${this.buildURL()}/${encodePath(backend)}/`;
    if (!isEmpty(id)) {
      url = url + encodePath(id);
    }

    return url;
  },

  pathForType() {
    return 'mounts';
  },

  optionsForQuery(id, action, wrapTTL) {
    let data = {};
    if (action === 'query') {
      data.list = true;
    }
    if (wrapTTL) {
      return { data, wrapTTL };
    }
    return { data };
  },

  fetchByQuery(query, action) {
    const { id, backend, wrapTTL } = query;
    return this.ajax(this.urlForSecret(backend, id), 'GET', this.optionsForQuery(id, action, wrapTTL)).then(
      resp => {
        if (wrapTTL) {
          return resp;
        }
        resp.id = id;
        resp.backend = backend;
        return resp;
      }
    );
  },

  query(store, type, query) {
    return this.fetchByQuery(query, 'query');
  },

  queryRecord(store, type, query) {
    return this.fetchByQuery(query, 'queryRecord');
  },
});
