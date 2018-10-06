/* eslint-disable */
import { isEmpty } from '@ember/utils';
import ApplicationAdapter from './application';

export default ApplicationAdapter.extend({
  namespace: 'v1',
  _url(backend, id, infix = 'data') {
    let url = `${this.buildURL()}/${backend}/${infix}/`;
    if (!isEmpty(id)) {
      url = url + id;
    }
    return url;
  },

  urlForFindRecord(id) {
    let [backend, path, version] = JSON.parse(id);
    return this._url(backend, path) + `?version=${version}`;
  },

  urlForCreateRecord(modelName, snapshot) {
    let backend = snapshot.belongsTo('secret').belongsTo('engine').id;
    let path = snapshot.attr('path');
    return this._url(backend, path);
  },

  createRecord(store, modelName, snapshot) {
    let backend = snapshot.belongsTo('secret').belongsTo('engine').id;
    let path = snapshot.attr('path');
    return this._super(...arguments).then(resp => {
      resp.id = JSON.stringify([backend, path, resp.version]);
      return resp;
    });
  },

  urlForUpdateRecord(id) {
    let [backend, path] = JSON.parse(id);
    return this._url(backend, path);
  },

  deleteRecord(store, type, snapshot) {
    // use adapterOptions to determine if it's delete or destroy for the version
    // deleteType should be 'delete', 'destroy', 'undelete'
    let infix = snapshot.adapterOptions.deleteType;
    let [backend, path, version] = JSON.parse(snapshot.id);

    return this.ajax(this._url(backend, path, infix), 'POST', { data: { versions: [version] } });
  },

  handleResponse(/*status, headers, payload, requestData*/) {
    // the body of the 404 will have some relevant information
    return this._super(...arguments);
  },
});
