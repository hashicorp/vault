import ApplicationAdapter from '../application';
import { encodePath } from 'vault/utils/path-encoding-helpers';
export default class PkiKeyAdapter extends ApplicationAdapter {
  namespace = 'v1';

  createRecord(store, type, snapshot) {
    const { record, adapterOptions } = snapshot;
    let url = this.getUrl(record.backend);
    if (adapterOptions.import) {
      url = url + '/import';
    } else {
      url = url + '/generate/' + record.type;
    }
    return this.ajax(url, 'POST', { data: this.serialize(snapshot) }).then((resp) => {
      return resp;
    });
  }

  getUrl(backend, id) {
    const url = `${this.buildURL()}/${encodePath(backend)}`;
    if (id) {
      return url + '/key/' + encodePath(id);
    }
    return url + '/keys';
  }

  query(store, type, query) {
    const { backend } = query;
    return this.ajax(this.getUrl(backend), 'GET', { data: { list: true } });
  }

  queryRecord(store, type, query) {
    const { backend, id } = query;
    return this.ajax(this.getUrl(backend, id), 'GET');
  }

  deleteRecord(store, type, snapshot) {
    const { id, record } = snapshot;
    return this.ajax(this.getUrl(record.backend, id), 'DELETE');
  }
}
