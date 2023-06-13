/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import ApplicationAdapter from '../application';
import { encodePath } from 'vault/utils/path-encoding-helpers';
import { assert } from '@ember/debug';

export default class KvDataAdapter extends ApplicationAdapter {
  namespace = 'v1';

  _urlForSecret(backend, path, version) {
    const base = `${this.buildURL()}/${encodePath(backend)}/data/${encodePath(path)}`;
    return version ? base + `?version=${version}` : base;
  }

  createRecord(store, type, snapshot) {
    const { backend, path, version } = snapshot.record;
    const url = this._urlForSecret(backend, path);
    return this.ajax(url, 'POST', { data: this.serialize(snapshot) }).then((resp) => {
      resp.id = `${encodePath(backend)}/${version}/${encodePath(path)}`;
      return resp;
    });
  }

  updateRecord(store, type, snapshot) {
    const { backend, path, version } = snapshot.record;
    const data = this.serialize(snapshot);
    const url = this._urlForSecret(backend, path, version);
    return this.ajax(url, 'POST', { data });
  }

  query(store, type, query) {
    const { path, backend, version } = query;
    return this.ajax(this._urlForSecret(backend, path, version), 'GET');
  }

  queryRecord(store, type, query) {
    const { path, backend, version } = query;
    return this.ajax(this._urlForSecret(backend, path, version), 'GET').then((resp) => {
      resp.id = `${encodePath(backend)}/${version}/${encodePath(path)}`;
      return resp;
    });
  }

  /* Five types of delete operations */
  deleteRecord(store, type, snapshot) {
    const { backend, path, deleteType, deleteVersions } = snapshot;
    if (!backend || !path) {
      throw new Error('The request to delete or undelete is missing required attributes.');
    }

    switch (deleteType) {
      case 'delete-latest-version':
        return this.ajax(this._urlForSecret(backend, path), 'DELETE');
      case 'delete-specific-version':
        return this.ajax(this._urlForSecret(backend, path), 'POST', {
          data: { deleteVersions },
        });
      case 'destroy-specific-version':
        return this.ajax(`${this.buildURL()}/${encodePath(backend)}/destroy/${encodePath(path)}`, 'PUT', {
          data: { deleteVersions },
        });
      case 'destroy-everything':
        return this.ajax(this._urlForSecret(backend, path), 'POST', {
          data: { deleteVersions },
        });
      case 'undelete-specific-version':
        return this.ajax(`${this.buildURL()}/${encodePath(backend)}/undelete/${encodePath(path)}`, 'POST', {
          data: { deleteVersions },
        });
      default:
        assert(
          'deletType must be one of delete-latest-version, delete-specific-version, destroy-specific-version, destroy-everything, undelete-specific-version.'
        );
    }
  }
}
