import PkiRoleAdapter from './pki-role';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default class PkiRoleEngineAdapter extends PkiRoleAdapter {
  _urlForRole(backend, id) {
    let url = `${this.buildURL()}/${encodePath(backend)}/roles`;
    if (id) {
      url = url + '/' + encodePath(id);
    }
    return url;
  }

  createRecord(store, type, snapshot) {
    let name = snapshot.attr('name');
    let url = this._urlForRole(snapshot.record.backend, name);
    return this.ajax(url, 'POST', { data: this.serialize(snapshot) }).then(() => {
      return {
        id: name,
        name,
        backend: snapshot.record.backend,
      };
    });
  }
}
