/* eslint-disable */
import AdapterError from '@ember-data/adapter/error';

import { isEmpty } from '@ember/utils';
import { get } from '@ember/object';
import ApplicationAdapter from './application';
import { encodePath } from 'vault/utils/path-encoding-helpers';
import { inject as service } from '@ember/service';

export default ApplicationAdapter.extend({
  store: service(),
  namespace: 'v1',

  _url(backend, id, infix = 'data') {
    let url = `${this.buildURL()}/${encodePath(backend)}/${infix}/`;
    if (!isEmpty(id)) {
      url = url + encodePath(id);
    }
    return url;
  },

  urlForFindRecord(id) {
    let [backend, path, version] = JSON.parse(id);
    let base = this._url(backend, path);
    return version ? base + `?version=${version}` : base;
  },

  urlForQueryRecord(id) {
    return this.urlForFindRecord(id);
  },

  findRecord() {
    return this._super(...arguments).catch((errorOrModel) => {
      // if the response is a real 404 or if the secret is gated by a control group this will be an error,
      // otherwise the response will be the body of a deleted / destroyed version
      if (errorOrModel instanceof AdapterError) {
        throw errorOrModel;
      }
      return errorOrModel;
    });
  },

  async getSecretDataVersion(backend, id) {
    // used in secret-edit route when you don't have current version and you need it for pulling the correct secret-v2-version record
    let url = this._url(backend, id);
    let response = await this.ajax(this._url(backend, id), 'GET');
    return response.data.metadata.version;
  },

  queryRecord(id, options) {
    return this.ajax(this.urlForQueryRecord(id), 'GET', options).then((resp) => {
      if (options.wrapTTL) {
        return resp;
      }
      resp.id = id;
      resp.backend = backend;
      return resp;
    });
  },

  querySecretDataByVersion(id) {
    return this.ajax(this.urlForQueryRecord(id), 'GET')
      .then((resp) => {
        return resp.data;
      })
      .catch((error) => {
        return error.data;
      });
  },

  urlForCreateRecord(modelName, snapshot) {
    let backend = snapshot.belongsTo('secret').belongsTo('engine').id;
    let path = snapshot.attr('path');
    return this._url(backend, path);
  },

  createRecord(store, modelName, snapshot) {
    let backend = snapshot.belongsTo('secret').belongsTo('engine').id;
    let path = snapshot.attr('path');
    return this._super(...arguments).then((resp) => {
      resp.id = JSON.stringify([backend, path, resp.version]);
      return resp;
    });
  },

  urlForUpdateRecord(id) {
    let [backend, path] = JSON.parse(id);
    return this._url(backend, path);
  },

  async deleteLatestVersion(backend, path) {
    try {
      await this.ajax(this._url(backend, path, 'data'), 'DELETE');
      let model = this.store.peekRecord('secret-v2-version', path);
      await model.reload();
      return model && model.rollbackAttributes();
    } catch (e) {
      return e;
    }
  },

  async undeleteVersion(backend, path, currentVersionForNoReadMetadata) {
    try {
      await this.ajax(this._url(backend, path, 'undelete'), 'POST', {
        data: { versions: [currentVersionForNoReadMetadata] },
      });
      let model = this.store.peekRecord('secret-v2-version', path);
      await model.reload();
      return model && model.rollbackAttributes();
    } catch (e) {
      return e;
    }
  },

  async softDelete(backend, path, version) {
    try {
      await this.ajax(this._url(backend, path, 'delete'), 'POST', {
        data: { versions: [version] },
      });
      let model = this.store.peekRecord('secret-v2-version', path);
      await model.reload();
      return model && model.rollbackAttributes();
    } catch (e) {
      return e;
    }
  },

  async deleteByDeleteType(backend, path, deleteType, version) {
    try {
      await this.ajax(this._url(backend, path, deleteType), 'POST', {
        data: { versions: [version] },
      });
      let model = this.store.peekRecord('secret-v2-version', path);
      await model.reload();
      return model && model.rollbackAttributes();
    } catch (e) {
      return e;
    }
  },

  v2DeleteOperation(store, id, deleteType = 'delete', currentVersionForNoReadMetadata) {
    let [backend, path, version] = JSON.parse(id);
    // deleteType should be 'delete', 'destroy', 'undelete', 'delete-latest-version', 'destroy-version'
    if (
      (currentVersionForNoReadMetadata && deleteType === 'delete') ||
      deleteType === 'delete-latest-version'
    ) {
      // moved to async to away model reload which is a promise
      return this.deleteLatestVersion(backend, path);
    } else if (deleteType === 'undelete' && !version) {
      // happens when no read access to metadata
      return this.undeleteVersion(backend, path, currentVersionForNoReadMetadata);
    } else if (deleteType === 'soft-delete') {
      return this.softDelete(backend, path, version);
    } else {
      version = version || currentVersionForNoReadMetadata;
      return this.deleteByDeleteType(backend, path, deleteType, version);
    }
  },

  handleResponse(status, headers, payload, requestData) {
    // the body of the 404 will have some relevant information
    if (status === 404 && get(payload, 'data.metadata')) {
      return this._super(200, headers, payload, requestData);
    }
    return this._super(...arguments);
  },
});
