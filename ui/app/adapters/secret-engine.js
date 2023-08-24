/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

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
      // TODO kv engine cleanup - this logic can be removed when KV exists in separate ember engine
      // if kv2 then add the config data to the mountModel
      // version comes in as a string
      if (mountModel?.data?.type === 'kv' && mountModel?.data?.options?.version === '2') {
        configModel = await this.ajax(this.urlForConfig(query.path), 'GET');
        mountModel.data = { ...mountModel.data, ...configModel.data };
      }
    } catch (error) {
      // no path means this was an error on listing
      if (!query.path || !mountModel) {
        throw error;
      }
      // control groups will throw a 403 permission denied error. If this happens return the mountModel
      // error is handled on routing
    }
    return mountModel;
  },

  async createRecord(store, type, snapshot) {
    const serializer = store.serializerFor(type.modelName);
    let data = serializer.serialize(snapshot);
    const path = snapshot.attr('path');
    // for kv2 we make two network requests
    data.config.id = path; // config relationship needs an id so use path for now
    if (data.type === 'kv' && data.options.version === 2) {
      // data has both data for sys mount and the config, we need to separate them
      const splitObjects = splitObject(data, ['max_versions', 'delete_version_after', 'cas_required']);
      let configData;
      [configData, data] = splitObjects;

      if (!data.id) {
        data.id = path;
      }
      // first create the engine
      await this.ajax(this.url(path), 'POST', { data });

      // second post to config
      try {
        await this.ajax(this.urlForConfig(path), 'POST', { data: configData });
      } catch (e) {
        // error here means you do not have update capabilities to config endpoint. If that's the case we show a flash message in the component and continue with the transition.
        // the error is handled by mount-backend-form component which checks capabilities before hitting the save to the adapter.
        // we do not handle the error here because we want the secret-engine to mount successfully and to continue the flow.
      }
      return {
        data: assign({}, data, { path: path + '/', id: path }),
      };
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
      return this.ajax(`/v1/${encodePath(query.backend)}/config/lease`, 'GET').then((resp) => {
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
    const { data } = snapshot.adapterOptions;
    const path = encodePath(snapshot.id);
    return this.ajax(`/v1/${path}/config/root`, 'POST', { data });
  },

  saveAWSLease(store, type, snapshot) {
    const { data } = snapshot.adapterOptions;
    const path = encodePath(snapshot.id);
    return this.ajax(`/v1/${path}/config/lease`, 'POST', { data });
  },

  saveZeroAddressConfig(store, type, snapshot) {
    const path = encodePath(snapshot.id);
    const roles = store.peekAll('role-ssh').filterBy('zeroAddress').mapBy('id').join(',');
    const url = `/v1/${path}/config/zeroaddress`;
    const data = { roles };
    if (roles === '') {
      return this.ajax(url, 'DELETE');
    }
    return this.ajax(url, 'POST', { data });
  },
});
