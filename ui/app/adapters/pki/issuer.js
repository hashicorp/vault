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
    const baseUrl = `${this.buildURL()}/${encodePath(backend)}`;
    if (id) {
      return `${baseUrl}/issuer/${encodePath(id)}`;
    } else {
      return `${baseUrl}/issuers`;
    }
  }

  createRecord(store, type, snapshot) {
    const { record, adapterOptions } = snapshot;
    let url = this.urlForQuery(record.backend);
    if (adapterOptions.import) {
      url = `${url}/import/bundle`;
    } else {
      // TODO WIP generate root or intermediate CSR actions from issuers index page
      // certType = 'root' || 'intermediate',   // record.type is internal or exported
      // url = ` ${url}/generate/${certType}/${record.type}`;
      throw new Error('createRecord method in adapters/pki/issuer.js is incomplete.');
    }
    return this.ajax(url, 'POST', { data: this.serialize(snapshot) }).then((resp) => {
      return resp;
    });
  }

  updateRecord(store, type, snapshot) {
    const { backend, issuerId } = snapshot.record;
    const data = this.serialize(snapshot);
    const url = this.urlForQuery(backend, issuerId);
    return this.ajax(url, 'POST', { data });
  }

  query(store, type, query) {
    return this.ajax(this.urlForQuery(query.backend), 'GET', this.optionsForQuery());
  }

  queryRecord(store, type, query) {
    const { backend, id } = query;
    return this.ajax(this.urlForQuery(backend, id), 'GET', this.optionsForQuery(id));
  }
}
