import { assign } from '@ember/polyfills';
import { allSettled } from 'rsvp';
import ApplicationAdapter from './application';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default ApplicationAdapter.extend({
  namespace: 'v1',

  createOrUpdate(store, type, snapshot) {
    const serializer = store.serializerFor(type.modelName);
    const data = serializer.serialize(snapshot);
    const { id } = snapshot;
    let url = this.urlForTransformations(snapshot.record.get('backend'), id);

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
    return this.ajax(this.urlForTransformations(snapshot.record.get('backend'), id), 'DELETE');
  },

  pathForType() {
    return 'transform';
  },

  urlForTransformations(backend, id) {
    let url = `${this.buildURL()}/${encodePath(backend)}/transformation`;
    if (id) {
      url = url + '/' + encodePath(id);
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
    const queryAjax = this.ajax(this.urlForTransformations(backend, id), 'GET', this.optionsForQuery(id));

    return allSettled([queryAjax]).then(results => {
      // query result 404d, so throw the adapterError
      if (!results[0].value) {
        throw results[0].reason;
      }
      let resp = {
        id,
        name: id,
        backend,
        data: {},
      };

      results.forEach(result => {
        if (result.value) {
          if (result.value.data.roles) {
            resp.data = assign({}, resp.data, { zero_address_roles: result.value.data.roles });
          } else {
            resp.data = assign({}, resp.data, result.value.data);
          }
        }
      });
      return resp;
    });
  },

  query(store, type, query) {
    return this.fetchByQuery(store, query);
  },

  queryRecord(store, type, query) {
    return this.fetchByQuery(store, query);
  },
});
