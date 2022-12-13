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

  query(store, type, query) {
    const { backend } = query;
    return this.ajax(this.getURL(backend), 'GET', { data: { list: true } }).then((resp) => {
      return resp.data.keys.map((name) => ({ name, backend }));
    });
  }
  queryRecord(store, type, query) {
    const { backend, name } = query;
    return this.ajax(this.getURL(backend, name), 'GET').then((resp) => {
      resp.data.backend = backend;
      resp.data.name = name;
      return resp.data;
    });
  }
  generateCredentials(backend, data) {
    const generateCredentialsUrl = `${this.buildURL()}/${encodePath(backend)}/creds/${data.role}`;

    return this.ajax(generateCredentialsUrl, 'POST', { data }).then((response) => {
      const { lease_id, lease_duration, data } = response;

      return {
        lease_id,
        lease_duration,
        ...data,
      };
    });
  }
}
