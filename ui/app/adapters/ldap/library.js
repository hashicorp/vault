/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import NamedPathAdapter from 'vault/adapters/named-path';
import { encodePath } from 'vault/utils/path-encoding-helpers';

const DIRECTORY_SEPARATOR = '/';
const API_ENDPOINTS = {
  STATUS: 'status',
  CHECK_OUT: 'check-out',
  CHECK_IN: 'check-in',
};

export default class LdapLibraryAdapter extends NamedPathAdapter {
  // path could be the library name (full path) or just part of the path i.e. west-account/
  _getURL(backend, path) {
    const base = `${this.buildURL()}/${encodePath(backend)}/library`;
    return path ? `${base}/${path}` : base;
  }

  urlForUpdateRecord(name, modelName, snapshot) {
    // For update operations, use completeLibraryName to ensure hierarchical libraries
    // (e.g., path_to_library="service-account/" + name="sc100") are correctly combined into the full path "service-account/sc100"
    const { backend, completeLibraryName } = snapshot.record;
    return this._getURL(backend, completeLibraryName);
  }

  urlForDeleteRecord(name, modelName, snapshot) {
    // For delete operations, use completeLibraryName to ensure hierarchical libraries
    // (e.g., path_to_library="service-account/" + name="sa") are correctly combined into the full path "service-account/sa"
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

  async queryRecord(store, type, query) {
    const { backend, name } = query;

    // Decode URL-encoded hierarchical paths (e.g., "service-account1%2Fsa1" -> "service-account1/sa1")
    const decodedName = decodeURIComponent(name);
    const resp = await this.ajax(this._getURL(backend, decodedName), 'GET');

    // If the decoded name contains a slash, it's hierarchical
    if (decodedName.includes(DIRECTORY_SEPARATOR)) {
      const lastSlashIndex = decodedName.lastIndexOf(DIRECTORY_SEPARATOR);
      const path_to_library = decodedName.substring(0, lastSlashIndex + 1);
      const libraryName = decodedName.substring(lastSlashIndex + 1);

      return {
        ...resp.data,
        backend,
        name: libraryName,
        path_to_library,
      };
    }

    // For non-hierarchical libraries, return as-is
    return { ...resp.data, backend, name: decodedName };
  }

  async fetchStatus(backend, completeLibraryName) {
    // The completeLibraryName parameter should be the full hierarchical path
    // (e.g., "service-account/sa") when called from the model's fetchStatus() method

    const url = `${this._getURL(backend, completeLibraryName)}/${API_ENDPOINTS.STATUS}`;
    const resp = await this.ajax(url, 'GET');

    const statuses = [];
    for (const key in resp.data) {
      const status = {
        ...resp.data[key],
        account: key,
        library: completeLibraryName,
      };
      statuses.push(status);
    }
    return statuses;
  }

  async checkOutAccount(backend, completeLibraryName, ttl) {
    // The completeLibraryName parameter should be the full hierarchical path
    // (e.g., "service-account/sa") when called from the model's checkOutAccount() method

    const url = `${this._getURL(backend, completeLibraryName)}/${API_ENDPOINTS.CHECK_OUT}`;

    return this.ajax(url, 'POST', { data: { ttl } }).then((resp) => {
      const { lease_id, lease_duration, renewable } = resp;
      const { service_account_name: account, password } = resp.data;
      return { account, password, lease_id, lease_duration, renewable };
    });
  }

  async checkInAccount(backend, completeLibraryName, service_account_names) {
    // The completeLibraryName parameter should be the full hierarchical path
    // (e.g., "service-account/sa") when called from the model's checkInAccount() method

    const url = `${this._getURL(backend, completeLibraryName)}/${API_ENDPOINTS.CHECK_IN}`;

    return this.ajax(url, 'POST', { data: { service_account_names } }).then((resp) => resp.data);
  }
}
