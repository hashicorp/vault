import ApplicationAdapter from './application';
import Ember from 'ember';

export default ApplicationAdapter.extend({
  namespace: 'v1',

  defaultSerializer: 'role',

  createOrUpdate(store, type, snapshot, requestType) {
    const serializer = store.serializerFor(type.modelName);
    const data = serializer.serialize(snapshot, requestType);
    const { id } = snapshot;
    let url = this.urlForRole(snapshot.record.get('backend'), id);

    return this.ajax(url, 'POST', { data });
  },

  createRecord() {
    return this.createOrUpdate(...arguments);
  },

  updateRecord() {
    return this.createOrUpdate(...arguments, 'update');
  },

  deleteRecord(store, type, snapshot) {
    const { id } = snapshot;
    return this.ajax(this.urlForRole(snapshot.record.get('backend'), id), 'DELETE');
  },

  pathForType() {
    return 'roles';
  },

  urlForRole(backend, id) {
    let url = `${this.buildURL()}/${backend}/roles`;
    if (id) {
      url = url + '/' + id;
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
    const { id, backend } = query;
    let zeroAddressAjax = Ember.RSVP.resolve();
    const queryAjax = this.ajax(this.urlForRole(backend, id), 'GET', this.optionsForQuery(id));
    if (!id) {
      zeroAddressAjax = this.findAllZeroAddress(store, query);
    }

    return Ember.RSVP.allSettled([queryAjax, zeroAddressAjax]).then(results => {
      // query result 404d, so throw the adapterError
      if (!results[0].value) {
        throw results[0].reason;
      }
      let resp = {
        id,
        name: id,
        backend,
      };

      results.forEach(result => {
        if (result.value) {
          if (result.value.data.roles) {
            resp = Ember.assign({}, resp, { zero_address_roles: result.value.data.roles });
          } else {
            resp = Ember.assign({}, resp, result.value);
          }
        }
      });
      return resp;
    });
  },

  findAllZeroAddress(store, query) {
    const { backend } = query;
    const url = `/v1/${backend}/config/zeroaddress`;
    return this.ajax(url, 'GET');
  },

  query(store, type, query) {
    return this.fetchByQuery(store, query);
  },

  queryRecord(store, type, query) {
    return this.fetchByQuery(store, query);
  },
});
