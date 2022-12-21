import ApplicationAdapter from '../application';
import { encodePath } from 'vault/utils/path-encoding-helpers';
export default class PkiKeyAdapter extends ApplicationAdapter {
  namespace = 'v1';

  createRecord(store, type, snapshot) {
    const { record } = snapshot;
    const url = this.getUrl(record.backend) + '/generate/' + record.type;
    return this.ajax(url, 'POST', { data: this.serialize(snapshot) }).then((resp) => {
      return resp;
    });
  }

  updateRecord(store, type, snapshot) {
    const { record } = snapshot;
    const { key_name } = this.serialize(snapshot);
    const url = this.getUrl(record.backend, record.id);
    return this.ajax(url, 'POST', { data: { key_name } });
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
