import ApplicationAdapter from '../application';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default class PkiKeyAdapter extends ApplicationAdapter {
  namespace = 'v1';

  optionsForQuery(id) {
    const data = {};
    if (!id) {
      data['list'] = true;
    }
    return { data };
  }

  _urlForKey(backend, id) {
    let url = `${this.buildURL()}/${encodePath(backend)}/key`;
    if (id) {
      url = url + '/' + encodePath(id);
    }
    return url;
  }

  urlForKeys(backend, method) {
    let url = `${this.buildURL()}/${encodePath(backend)}/keys`;
    // methods can be 'generate' or 'import'
    if (method) {
      url = url + '/' + method;
    }
    return url;
  }

  query(store, type, query) {
    const { backend } = query;
    return this.ajax(this.urlForKeys(backend), 'GET', { list: true });
  }

  queryRecord(store, type, query) {
    const { backend, id } = query;
    return this.ajax(this._urlForKey(backend, id), 'GET');
  }
}
