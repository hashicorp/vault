import ApplicationAdapter from '../application';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default ApplicationAdapter.extend({
  namespace: 'v1',
  pathForType(type) {
    return type.replace('kmip/', '');
  },

  _url(modelType, meta = {}, id) {
    let { backend, scope, role } = meta;
    let type = this.pathForType(modelType);
    let base;
    switch (type) {
      case 'scope':
        base = `${encodePath(backend)}/scope`;
        break;
      case 'role':
        base = `${encodePath(backend)}/scope/${encodePath(scope)}/role`;
        break;
      case 'credential':
        base = `${encodePath(backend)}/scope/${encodePath(scope)}/role/${encodePath(role)}/credential`;
        break;
    }
    if (id && type === 'credential') {
      return `/v1/${base}/lookup?serial_number=${encodePath(id)}`;
    }

    if (id) {
      return `/v1/${base}/${encodePath(id)}`;
    }
    return `/v1/${base}`;
  },

  urlForQuery(query, modelType) {
    let base = this._url(modelType, query);
    return base + '?list=true';
  },

  query(store, type, query) {
    return this.ajax(this.urlForQuery(query, type.modelName), 'GET').then(resp => {
      // remove pagination query items here
      const { size, page, responsePath, pageFilter, ...modelAttrs } = query;
      resp._requestQuery = modelAttrs;
      return resp;
    });
  },

  queryRecord(store, type, query) {
    let id = query.id;
    delete query.id;
    return this.ajax(this._url(type.modelName, query, id), 'GET').then(resp => {
      resp.id = id;
      resp = { ...resp, ...query };
      return resp;
    });
  },
  buildURL(modelName, id, snapshot, requestType, query) {
    if (requestType === 'createRecord') {
      return this._super(...arguments);
    }
    return this._super(`${modelName}`, id, snapshot, requestType, query);
  },
});
