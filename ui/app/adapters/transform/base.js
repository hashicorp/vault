// import { assign } from '@ember/polyfills';
// import { resolve, allSettled } from 'rsvp';
import ApplicationAdapter from '../application';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default ApplicationAdapter.extend({
  namespace: 'v1',

  pathForType(type) {
    return type.replace('transform/', '');
  },

  url(backend, modelType, id) {
    let type = this.pathForType(modelType);
    let url = `/${this.namespace}/${encodePath(backend)}/${encodePath(type)}`;
    if (id) {
      return `${url}/${encodePath(id)}`;
    }
    return url;
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
    // console.log({ type });
    return this.fetchByQuery(query);
  },

  queryRecord(store, type, query) {
    // console.log({ type });
    return this.ajax(this._url(type.modelName, query.backend, query.id), 'GET').then(result => {
      // console.log(result);

      return result;
    });
  },

  // buildUrl(modelName, id, snapshot, requestType, query, returns) {},
});
