/* eslint-disable */
import { isEmpty } from '@ember/utils';
import ApplicationAdapter from './application';

export default ApplicationAdapter.extend({
  namespace: 'v1',
  _url(backend, id) {
    let url = `${this.buildURL()}/${backend}/metadata/`;
    if (!isEmpty(id)) {
      url = url + id;
    }
    return url;
  },

  // we override query here because the query object has a bunch of client-side
  // concerns and we only want to send "list" to the server
  query(store, type, query) {
    let { backend, id } = query;
    return this.ajax(this._url(backend, id), 'GET', { data: { list: true } }).then(resp => {
      resp.id = id;
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

  urlForUpdateRecord(store, type, snapshot) {
    let backend = snapshot.belongsTo('engine', { id: true });
    let { id } = snapshot;
    return this._url(backend, id);
  },

  urlForDeleteRecord(store, type, snapshot) {
    let backend = snapshot.belongsTo('engine', { id: true });
    let { id } = snapshot;
    return this._url(backend, id);
  },
});
