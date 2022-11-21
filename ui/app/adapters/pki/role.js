import ApplicationAdapter from '../application';
import { assign } from '@ember/polyfills';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default class PkiRoleAdapter extends ApplicationAdapter {
  namespace = 'v1';

  _urlForRole(backend, id) {
    let url = `${this.buildURL()}/${encodePath(backend)}/roles`;
    if (id) {
      url = url + '/' + encodePath(id);
    }
    return url;
  }

  _optionsForQuery(id) {
    const data = {};
    if (!id) {
      data['list'] = true;
    }
    return { data };
  }

  createRecord(store, type, snapshot) {
    const name = snapshot.attr('name');
    const url = this._urlForRole(snapshot.record.backend, name);

    return this.ajax(url, 'POST', { data: this.serialize(snapshot) }).then(() => {
      return {
        id: name,
        name,
        backend: snapshot.record.backend,
      };
    });
  }

  fetchByQuery(store, query) {
    const { id, backend } = query;

    return this.ajax(this._urlForRole(backend, id), 'GET', this._optionsForQuery(id)).then((resp) => {
      const data = {
        id,
        name: id,
        backend,
      };

      return assign({}, resp, data);
    });
  }

  query(store, type, query) {
    return this.fetchByQuery(store, query);
  }

  queryRecord(store, type, query) {
    return this.fetchByQuery(store, query);
  }
}
