/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import NamedPathAdapter from 'vault/adapters/named-path';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default class LdapLibraryAdapter extends NamedPathAdapter {
  // path could be the library name (full path) or just part of the path i.e. west-account/
  _getURL(backend, path) {
    const base = `${this.buildURL()}/${encodePath(backend)}/library`;
    return path ? `${base}/${path}` : base;
  }

  urlForUpdateRecord(name, modelName, snapshot) {
    // when editing the name IS the full path so we can use "name" instead of "completeLibraryName" here
    return this._getURL(snapshot.attr('backend'), name);
  }
  urlForDeleteRecord(name, modelName, snapshot) {
    const { backend, completeLibraryName } = snapshot.record;
    return this._getURL(backend, completeLibraryName);
  }

  query(store, type, query) {
    const { backend, path_to_library } = query;
    // if we have a path_to_library then we're listing subdirectories at a hierarchical library path (i.e west-account/my-library)
    const url = this._getURL(backend, path_to_library);
    return this.ajax(url, 'GET', { data: { list: true } })
      .then((resp) => {
        return resp.data.keys.map((name) => ({ name, backend, path_to_library }));
      })
      .catch((error) => {
        if (error.httpStatus === 404) {
          return [];
        }
        throw error;
      });
  }
  queryRecord(store, type, query) {
    const { backend, name } = query;
    return this.ajax(this._getURL(backend, name), 'GET').then((resp) => ({ ...resp.data, backend, name }));
  }

  fetchStatus(backend, name) {
    const url = `${this._getURL(backend, name)}/status`;
    return this.ajax(url, 'GET').then((resp) => {
      const statuses = [];
      for (const key in resp.data) {
        const status = {
          ...resp.data[key],
          account: key,
          library: name,
        };
        statuses.push(status);
      }
      return statuses;
    });
  }
  checkOutAccount(backend, name, ttl) {
    const url = `${this._getURL(backend, name)}/check-out`;
    return this.ajax(url, 'POST', { data: { ttl } }).then((resp) => {
      const { lease_id, lease_duration, renewable } = resp;
      const { service_account_name: account, password } = resp.data;
      return { account, password, lease_id, lease_duration, renewable };
    });
  }
  checkInAccount(backend, name, service_account_names) {
    const url = `${this._getURL(backend, name)}/check-in`;
    return this.ajax(url, 'POST', { data: { service_account_names } }).then((resp) => resp.data);
  }
}
