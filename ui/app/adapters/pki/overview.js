import ApplicationAdapter from '../application';
// import { encodePath } from 'vault/utils/path-encoding-helpers';
// import ControlGroupError from '../../lib/control-group-error';

export default class KeymgmtKeyAdapter extends ApplicationAdapter {
  namespace = 'v1';

  pathForType() {
    // backend name prepended in buildURL method
    return 'pki';
  }

  // buildURL(modelName, id, snapshot, requestType, query) {
  //   let url = super.buildURL(...arguments);
  // if (snapshot) {
  //   url = url.replace('pki', `${snapshot.attr('backend')}/certs`);
  // } else if (query) {
  //   url = url.replace('pki', `${query.backend}/certs`);
  // }
  //   return url;
  // }

  // optionsForQuery(id) {
  //   let data = {};
  //   if (!id) {
  //     data['list'] = true;
  //   }
  //   return { data };
  // }

  // urlFor(backend, id) {
  //   let url = `${this.buildURL()}/${backend}/certs`;
  //   if (id) {
  //     url = `${this.buildURL()}/${backend}/cert/${id}`;
  //   }
  //   return url;
  // }

  // fetchByQuery(store, query) {
  //   const { backend, id } = query;
  //   return this.ajax(this.urlFor(backend, id), 'GET', this.optionsForQuery(id)).then((resp) => {
  //     const data = {
  //       backend,
  //     };
  //     if (id) {
  //       data.serial_number = id;
  //       data.id = id;
  //       data.id_for_nav = `cert/${id}`;
  //     }
  //     return assign({}, resp, data);
  //   });
  // }

  // async queryRecord(store, type, query) {
  //   return this.fetchByQuery(store, query);
  // }

  // async query(store, type, query) {
  //   return this.fetchByQuery(store, query);
  // }
}
