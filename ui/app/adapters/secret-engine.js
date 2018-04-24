import Ember from 'ember';
import ApplicationAdapter from './application';

export default ApplicationAdapter.extend({
  url(path) {
    const url = `${this.buildURL()}/mounts`;
    return path ? url + '/' + path : url;
  },

  pathForType() {
    return 'mounts';
  },

  query(store, type, query) {
    let url = `/${this.urlPrefix()}/internal/ui/mounts`;
    if (query.path) {
      url = `${url}/${query.path}`;
    }
    return this.ajax(url, 'GET');
  },

  createRecord(store, type, snapshot) {
    const serializer = store.serializerFor(type.modelName);
    const data = serializer.serialize(snapshot);
    const path = snapshot.attr('path');

    return this.ajax(this.url(path), 'POST', { data }).then(() => {
      // ember data doesn't like 204s if it's not a DELETE
      return {
        data: Ember.assign({}, data, { path: path + '/', id: path }),
      };
    });
  },

  findRecord(store, type, path, snapshot) {
    if (snapshot.attr('type') === 'ssh') {
      return this.ajax(`/v1/${path}/config/ca`, 'GET');
    }
    return;
  },

  queryRecord(store, type, query) {
    if (query.type === 'aws') {
      return this.ajax(`/v1/${query.backend}/config/lease`, 'GET').then(resp => {
        resp.path = query.backend + '/';
        return resp;
      });
    }
    return;
  },

  updateRecord(store, type, snapshot) {
    const { apiPath, options, adapterMethod } = snapshot.adapterOptions;
    if (adapterMethod) {
      return this[adapterMethod](...arguments);
    }
    if (apiPath) {
      const serializer = store.serializerFor(type.modelName);
      const data = serializer.serialize(snapshot);
      const path = snapshot.id;
      return this.ajax(`/v1/${path}/${apiPath}`, options.isDelete ? 'DELETE' : 'POST', { data });
    }
  },

  saveAWSRoot(store, type, snapshot) {
    let { data } = snapshot.adapterOptions;
    const path = snapshot.id;
    return this.ajax(`/v1/${path}/config/root`, 'POST', { data });
  },

  saveAWSLease(store, type, snapshot) {
    let { data } = snapshot.adapterOptions;
    const path = snapshot.id;
    return this.ajax(`/v1/${path}/config/lease`, 'POST', { data });
  },

  saveZeroAddressConfig(store, type, snapshot) {
    const path = snapshot.id;
    const roles = store.peekAll('role-ssh').filterBy('zeroAddress').mapBy('id').join(',');
    const url = `/v1/${path}/config/zeroaddress`;
    const data = { roles };
    if (roles === '') {
      return this.ajax(url, 'DELETE');
    }
    return this.ajax(url, 'POST', { data });
  },
});
