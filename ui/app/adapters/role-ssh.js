/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { resolve, allSettled } from 'rsvp';
import ApplicationAdapter from './application';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default ApplicationAdapter.extend({
  namespace: 'v1',

  createOrUpdate(store, type, snapshot, requestType) {
    const { name, backend } = snapshot.record;
    const serializer = store.serializerFor(type.modelName);
    const data = serializer.serialize(snapshot, requestType);
    const url = this.urlForRole(backend, name);

    return this.ajax(url, 'POST', { data }).then((resp) => {
      // Ember data doesn't like 204 responses except for DELETE method
      const response = resp || { data: {} };
      response.data.name = name;
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
    return this.ajax(this.urlForRole(snapshot.record.backend, id), 'DELETE');
  },

  pathForType() {
    return 'roles';
  },

  urlForRole(backend, id) {
    let url = `${this.buildURL()}/${encodePath(backend)}/roles`;
    if (id) {
      url = url + '/' + encodePath(id);
    }
    return url;
  },

  optionsForQuery(id) {
    const data = {};
    if (!id) {
      data['list'] = true;
    }
    return { data };
  },

  fetchByQuery(store, query) {
    const { id, backend } = query;
    let zeroAddressAjax = resolve();
    const queryAjax = this.ajax(this.urlForRole(backend, id), 'GET', this.optionsForQuery(id));
    if (!id) {
      zeroAddressAjax = this.findAllZeroAddress(store, query);
    }

    return allSettled([queryAjax, zeroAddressAjax]).then((results) => {
      // query result 404d, so throw the adapterError
      if (!results[0].value) {
        throw results[0].reason;
      }
      const resp = {
        id,
        name: id,
        backend,
        data: {},
      };

      results.forEach((result) => {
        if (result.value) {
          if (result.value.data.roles) {
            resp.data = { ...resp.data, zero_address_roles: result.value.data.roles };
          } else {
            resp.data = { ...resp.data, ...result.value.data };
          }
        }
      });
      return resp;
    });
  },

  findAllZeroAddress(store, query) {
    const { backend } = query;
    const url = `/v1/${encodePath(backend)}/config/zeroaddress`;
    return this.ajax(url, 'GET');
  },

  query(store, type, query) {
    return this.fetchByQuery(store, query);
  },

  queryRecord(store, type, query) {
    return this.fetchByQuery(store, query);
  },
});
