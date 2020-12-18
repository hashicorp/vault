/* eslint-disable */
import { isEmpty } from '@ember/utils';
import ApplicationAdapter from './application';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default ApplicationAdapter.extend({
  namespace: 'v1',
  _url(backend, id) {
    let url = `${this.buildURL()}/${encodePath(backend)}/metadata/`;
    if (!isEmpty(id)) {
      url = url + encodePath(id);
    }
    return url;
  },

  // we override query here because the query object has a bunch of client-side
  // concerns and we only want to send "list" to the server
  query(store, type, query) {
    let { backend, id } = query;
    return this.ajax(this._url(backend, id), 'GET', { data: { list: true } }).then(resp => {
      resp.id = id;
      resp.backend = backend;
      return resp;
    });
  },

  urlForQueryRecord(query) {
    let { id, backend } = query;
    return this._url(backend, id);
  },

  queryRecord(store, type, query) {
    let { backend, id } = query;
    return this.ajax(this._url(backend, id), 'GET').then(resp => {
      resp.id = id;
      resp.backend = backend;
      return resp;
    });
  },

  detailURL(snapshot) {
    let backend = snapshot.belongsTo('engine', { id: true }) || snapshot.attr('engineId');
    let { id } = snapshot;
    return this._url(backend, id);
  },

  urlForUpdateRecord(store, type, snapshot) {
    return this.detailURL(snapshot);
  },
  urlForCreateRecord(modelName, snapshot) {
    return this.detailURL(snapshot);
  },
  urlForDeleteRecord(store, type, snapshot) {
    return this.detailURL(snapshot);
  },
});
