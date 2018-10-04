import { isEmpty } from '@ember/utils';
import ApplicationAdapter from './application';

export default ApplicationAdapter.extend({
  namespace: 'v1',
  _url(backend, id) {
    let url = `${this.buildURL()}/${backend}/data/`;
    if (!isEmpty(id)) {
      url = url + id;
    }
    return url;
  },

  urlForFindRecord(id) {
    let [backend, path, version] = JSON.parse(id);
    return this._url(backend, path) + `?version=${version}`;
  },

  deleteRecord(store, type, snapshot) {
    // use adapterOptions to determine if it's delete or destroy for the version
    return this._super(...arguments);
  },
});
