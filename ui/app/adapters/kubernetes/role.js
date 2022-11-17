import NamedPathAdapter from 'vault/adapters/named-path';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default class KubernetesRoleAdapter extends NamedPathAdapter {
  getURL(backend, name) {
    const base = `${this.buildURL()}/${encodePath(backend)}/roles`;
    return name ? `${base}/${name}` : base;
  }
  urlForQuery({ backend }) {
    return this.getURL(backend);
  }
  urlForUpdateRecord(name, modelName, snapshot) {
    return this.getURL(snapshot.attr('backend'), name);
  }
  urlForDeleteRecord(name, modelName, snapshot) {
    return this.getURL(snapshot.attr('backend'), name);
  }

  queryRecord(store, type, query) {
    const { backend, name } = query;
    return this.ajax(this.getURL(backend, name), 'GET').then((resp) => {
      resp.backend = backend;
      resp.name = name;
      return resp;
    });
  }
}
