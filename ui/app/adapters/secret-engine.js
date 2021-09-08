import { assign } from '@ember/polyfills';
import ApplicationAdapter from './application';
import { encodePath } from 'vault/utils/path-encoding-helpers';
import { splitObject } from 'vault/helpers/split-object';

export default ApplicationAdapter.extend({
  url(path) {
    const url = `${this.buildURL()}/mounts`;
    return path ? url + '/' + encodePath(path) : url;
  },

  urlForConfig(path) {
    return `/v1/${path}/config`;
  },

  internalURL(path) {
    let url = `/${this.urlPrefix()}/internal/ui/mounts`;
    if (path) {
      url = `${url}/${encodePath(path)}`;
    }
    return url;
  },

  pathForType() {
    return 'mounts';
  },

  async query(store, type, query) {
    let mountModel, configModel;
    try {
      mountModel = await this.ajax(this.internalURL(query.path), 'GET');
      // if kv2 then add the config data to the mountModel
      if (mountModel.data.type === 'kv' && mountModel.data.options.version === '2') {
        configModel = await this.ajax(this.urlForConfig(query.path), 'GET');
        let data = mountModel.data;
        data = { ...mountModel.data, ...configModel.data };
        mountModel.data = data;
      }
    } catch (error) {
      throw new Error(error);
    }
    return mountModel;
  },

  createRecord(store, type, snapshot) {
    const serializer = store.serializerFor(type.modelName);
    let data = serializer.serialize(snapshot);
    const path = snapshot.attr('path');
    // for kv2 we make two network requests
    // ARG TODO make sure to have it equal 2 and it might be a string. Need to double check.
    if (data.type === 'kv' && data.options.version !== 1) {
      console.log('here');
      // data has both data for sys mount and the config, we need to separate them
      let splitObjects = splitObject(data, ['max_versions', 'delete_version_after', 'cas_required']);
      let configData;
      [configData, data] = splitObjects;

      if (!data.id) {
        data.id = path;
      }
      // first create the engine
      return this.ajax(this.url(path), 'POST', { data })
        .then(() => {
          // second create the config request
          return this.ajax(this.urlForConfig(path), 'POST', { data: configData }).catch(() => {
            // error here means you do not have update capabilities to config endpoint. If that's the case we show a flash message in the component and continue with the transition.
            // the error is handled by mount-backend-form component because the save method returns an error and we filter through the error to determine if we should proceed or not.
            throw new Error('configIssue');
          });
        })
        .then(() => {
          return {
            data: assign({}, data, { path: path + '/', id: path }),
          };
        })
        .catch(e => {
          // if not mount error and already returned then you need to assign the id
          if (e.message === 'configIssue') {
            // still proceed with mount sys call and do not throw mountIssue error. The flash message warning is handled before saving
            return {
              data: assign({}, data, { path: path + '/', id: path }),
            };
          }
          // if they do not hit a config error then it's a sys/mount error and they should not proceed.
          if (e.httpStatus === 400) {
            throw new Error('samePath');
          }
          throw new Error('mountIssue');
        });
    } else {
      return this.ajax(this.url(path), 'POST', { data }).then(() => {
        // ember data doesn't like 204s if it's not a DELETE
        return {
          data: assign({}, data, { path: path + '/', id: path }),
        };
      });
    }
  },

  findRecord(store, type, path, snapshot) {
    if (snapshot.attr('type') === 'ssh') {
      return this.ajax(`/v1/${encodePath(path)}/config/ca`, 'GET');
    }
    return;
  },

  queryRecord(store, type, query) {
    if (query.type === 'aws') {
      return this.ajax(`/v1/${encodePath(query.backend)}/config/lease`, 'GET').then(resp => {
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
      const path = encodePath(snapshot.id);
      return this.ajax(`/v1/${path}/${apiPath}`, options.isDelete ? 'DELETE' : 'POST', { data });
    }
  },

  saveAWSRoot(store, type, snapshot) {
    let { data } = snapshot.adapterOptions;
    const path = encodePath(snapshot.id);
    return this.ajax(`/v1/${path}/config/root`, 'POST', { data });
  },

  saveAWSLease(store, type, snapshot) {
    let { data } = snapshot.adapterOptions;
    const path = encodePath(snapshot.id);
    return this.ajax(`/v1/${path}/config/lease`, 'POST', { data });
  },

  saveZeroAddressConfig(store, type, snapshot) {
    const path = encodePath(snapshot.id);
    const roles = store
      .peekAll('role-ssh')
      .filterBy('zeroAddress')
      .mapBy('id')
      .join(',');
    const url = `/v1/${path}/config/zeroaddress`;
    const data = { roles };
    if (roles === '') {
      return this.ajax(url, 'DELETE');
    }
    return this.ajax(url, 'POST', { data });
  },
});
