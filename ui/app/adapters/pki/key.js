import ApplicationAdapter from '../application';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default class PkiKeyAdapter extends ApplicationAdapter {
  namespace = 'v1';

  _urlForKey(backend, id) {
    let url = `${this.buildURL()}/${encodePath(backend)}/key`;
    if (id) {
      url = url + '/' + encodePath(id);
    }
    return url;
  }

  urlForQuery(backend, id) {
    let url = `${this.buildURL()}/${encodePath(backend)}/keys`;
    if (id) {
      url = url + '/' + encodePath(id);
    }
    return url;
  }

  query(store, type, query) {
    const { backend } = query;
    return this.ajax(this.urlForQuery(backend), 'GET', { data: { list: true } });
  }

  queryRecord(store, type, query) {
    const { backend, id } = query;
    return this.ajax(this._urlForKey(backend, id), 'GET');
  }
}
