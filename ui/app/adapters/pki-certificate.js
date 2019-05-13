import { assign } from '@ember/polyfills';
import Adapter from './pki';

export default Adapter.extend({
  url(role) {
    return `/v1/${role.backend}/issue/${role.name}`;
  },

  urlFor(backend, id) {
    let url = `${this.buildURL()}/${backend}/certs`;
    if (id) {
      url = `${this.buildURL()}/${backend}/cert/${id}`;
    }
    return url;
  },
  optionsForQuery(id) {
    let data = {};
    if (!id) {
      data['list'] = true;
    }
    return { data };
  },

  fetchByQuery(store, query) {
    const { backend, id } = query;
    return this.ajax(this.urlFor(backend, id), 'GET', this.optionsForQuery(id)).then(resp => {
      const data = {
        backend,
      };
      if (id) {
        data.serial_number = id;
        data.id = id;
        data.id_for_nav = `cert/${id}`;
      }
      return assign({}, resp, data);
    });
  },

  query(store, type, query) {
    return this.fetchByQuery(store, query);
  },

  queryRecord(store, type, query) {
    return this.fetchByQuery(store, query);
  },

  updateRecord(store, type, snapshot) {
    if (snapshot.adapterOptions.method !== 'revoke') {
      return;
    }
    const id = snapshot.id;
    const backend = snapshot.record.get('backend');
    const data = {
      serial_number: id,
    };
    return this.ajax(`${this.buildURL()}/${backend}/revoke`, 'POST', { data }).then(resp => {
      const data = {
        id,
        serial_number: id,
        backend,
      };
      return assign({}, resp, data);
    });
  },
});
