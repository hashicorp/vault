import ApplicationAdapater from '../application';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default ApplicationAdapater.extend({
  namespace: 'v1',
  modelName: 'transform/template',

  pathForType(type) {
    return type.replace('transform/', '');
  },

  _url(modelType, backend, id) {
    let type = this.pathForType(modelType);
    let base = `${this.buildURL()}/${encodePath(backend)}/${type}`;
    if (id) {
      return `${base}/${encodePath(id)}`;
    }
    return base + '?list=true';
  },

  query(store, type, query) {
    return this.ajax(this._url(this.modelName, query.backend), 'GET').then(result => {
      return result;
    });
  },

  // buildURL(modelName, id, snapshot, requestType, query) {
  //   return this._super(`${modelName}/`, id, snapshot, requestType, query);
  // },
});
