/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import NamedPathAdapter from 'vault/adapters/named-path';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default class LdapLibraryAdapter extends NamedPathAdapter {
  getURL(backend, name) {
    const base = `${this.buildURL()}/${encodePath(backend)}/library`;
    return name ? `${base}/${name}` : base;
  }

  urlForUpdateRecord(name, modelName, snapshot) {
    return this.getURL(snapshot.attr('backend'), name);
  }
  urlForDeleteRecord(name, modelName, snapshot) {
    return this.getURL(snapshot.attr('backend'), name);
  }

  query(store, type, query) {
    const { backend } = query;
    return this.ajax(this.getURL(backend), 'GET', { data: { list: true } })
      .then((resp) => {
        return resp.data.keys.map((name) => ({ name, backend }));
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
    return this.ajax(this.getURL(backend, name), 'GET').then((resp) => ({ ...resp.data, backend, name }));
  }

  fetchStatus(backend, name) {
    const url = `${this.getURL(backend, name)}/status`;
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
    const url = `${this.getURL(backend, name)}/check-out`;
    return this.ajax(url, 'POST', { data: { ttl } }).then((resp) => {
      const { lease_id, lease_duration, renewable } = resp;
      const { service_account_name: account, password } = resp.data;
      return { account, password, lease_id, lease_duration, renewable };
    });
  }
  checkInAccount(backend, name, service_account_names) {
    const url = `${this.getURL(backend, name)}/check-in`;
    return this.ajax(url, 'POST', { data: { service_account_names } }).then((resp) => resp.data);
  }
}
