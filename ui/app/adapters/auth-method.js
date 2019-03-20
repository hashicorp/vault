import { assign } from '@ember/polyfills';
import { get, set } from '@ember/object';
import ApplicationAdapter from './application';
import DS from 'ember-data';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default ApplicationAdapter.extend({
  url(path) {
    const url = `${this.buildURL()}/auth`;
    return path ? url + '/' + encodePath(path) : url;
  },

  // used in updateRecord on the model#tune action
  pathForType() {
    return 'mounts/auth';
  },

  findAll(store, type, sinceToken, snapshotRecordArray) {
    let isUnauthenticated = get(snapshotRecordArray || {}, 'adapterOptions.unauthenticated');
    if (isUnauthenticated) {
      let url = `/${this.urlPrefix()}/internal/ui/mounts`;
      return this.ajax(url, 'GET', {
        unauthenticated: true,
      })
        .then(result => {
          return {
            data: result.data.auth,
          };
        })
        .catch(() => {
          return {
            data: {},
          };
        });
    }
    return this.ajax(this.url(), 'GET').catch(e => {
      if (e instanceof DS.AdapterError) {
        set(e, 'policyPath', 'sys/auth');
      }
      throw e;
    });
  },

  createRecord(store, type, snapshot) {
    const serializer = store.serializerFor(type.modelName);
    const data = serializer.serialize(snapshot);
    const path = snapshot.attr('path');

    return this.ajax(this.url(path), 'POST', { data }).then(() => {
      // ember data doesn't like 204s if it's not a DELETE
      return {
        data: assign({}, data, { path: path + '/', id: path }),
      };
    });
  },

  urlForDeleteRecord(id, modelName, snapshot) {
    return this.url(snapshot.id);
  },

  exchangeOIDC(path, state, code) {
    return this.ajax(`/v1/auth/${encodePath(path)}/oidc/callback`, 'GET', { data: { state, code } });
  },
});
