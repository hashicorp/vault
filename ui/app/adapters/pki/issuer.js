import ApplicationAdapter from '../application';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default class PkiIssuerAdapter extends ApplicationAdapter {
  namespace = 'v1';

  optionsForQuery(id) {
    const data = {};
    if (!id) {
      data['list'] = true;
    }
    return { data };
  }

  urlForQuery(backend, id) {
    let url = `${this.buildURL()}/${encodePath(backend)}/issuers`;
    if (id) {
      url = url + '/' + encodePath(id);
    }
    return url;
  }

  query(store, type, query) {
    const { backend, id } = query;
    return this.ajax(this.urlForQuery(backend, id), 'GET', this.optionsForQuery(id));
  }
}
