/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from './application';
import { pluralize } from 'ember-inflector';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default ApplicationAdapter.extend({
  namespace: 'v1',

  createOrUpdate(store, type, snapshot, requestType) {
    const serializer = store.serializerFor(type.modelName);
    const data = serializer.serialize(snapshot, requestType);
    const name = snapshot.attr('name');
    let url = this.urlForSecret(snapshot.record.backend, name);
    if (requestType === 'update') {
      url = url + '/config';
    }

    return this.ajax(url, 'POST', { data }).then((resp) => {
      const response = resp || {};
      response.id = name;
      return response;
    });
  },

  createRecord() {
    return this.createOrUpdate(...arguments);
  },

  updateRecord() {
    return this.createOrUpdate(...arguments, 'update');
  },

  deleteRecord(store, type, snapshot) {
    const { id } = snapshot;
    return this.ajax(this.urlForSecret(snapshot.record.backend, id), 'DELETE');
  },

  pathForType(type) {
    let path;
    switch (type) {
      case 'cluster':
        path = 'clusters';
        break;
      case 'secret-engine':
        path = 'secrets';
        break;
      default:
        path = pluralize(type);
        break;
    }
    return path;
  },

  urlForSecret(backend, id) {
    let url = `${this.buildURL()}/${encodePath(backend)}/keys/`;
    if (id) {
      url += encodePath(id);
    }
    return url;
  },

  urlForAction(action, backend, id, param) {
    const urlBase = `${this.buildURL()}/${encodePath(backend)}/${action}`;
    // these aren't key-specific
    if (action === 'hash' || action === 'random') {
      return urlBase;
    }
    if (action === 'datakey' && param) {
      // datakey action has `wrapped` or `plaintext` as part of the url
      return `${urlBase}/${param}/${encodePath(id)}`;
    }
    if (action === 'export' && param) {
      const [type, version] = param;
      const exportBase = `${urlBase}/${type}-key/${encodePath(id)}`;
      return version ? `${exportBase}/${version}` : exportBase;
    }
    return `${urlBase}/${encodePath(id)}`;
  },

  optionsForQuery(id) {
    const data = {};
    if (!id) {
      data['list'] = true;
    }
    return { data };
  },

  fetchByQuery(query) {
    const { id, backend } = query;
    return this.ajax(this.urlForSecret(backend, id), 'GET', this.optionsForQuery(id)).then((resp) => {
      resp.id = id;
      resp.backend = backend;
      return resp;
    });
  },

  query(store, type, query) {
    return this.fetchByQuery(query);
  },

  queryRecord(store, type, query) {
    return this.fetchByQuery(query);
  },

  // rotate, encrypt, decrypt, sign, verify, hmac, rewrap, datakey
  keyAction(action, { backend, id, payload }, options = {}) {
    const verb = action === 'export' ? 'GET' : 'POST';
    const { wrapTTL } = options;
    if (action === 'rotate') {
      return this.ajax(this.urlForSecret(backend, id) + '/rotate', verb);
    }
    const { param } = payload;

    delete payload.param;
    return this.ajax(this.urlForAction(action, backend, id, param), verb, {
      data: payload,
      wrapTTL,
    });
  },
});
