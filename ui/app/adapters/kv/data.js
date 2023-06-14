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
      resp.id = `${backend}/data/${path}?version=${version}`;
      return resp;
    });
  }

  queryRecord(store, type, query) {
    const { path, backend, version } = query;
    return this.ajax(this._urlForSecret(backend, path, version), 'GET').then((resp) => {
      resp.id = `${backend}/data/${path}?version=${version}`;
      return resp;
    });
  }

  findRecord(store, type, id) {
    return this.ajax(`${this.buildURL()}/${id}`, 'GET').then((resp) => {
      resp.id = id;
      return resp;
    });
  }

  /* Five types of delete operations */
  deleteRecord(store, type, snapshot) {
    const { backend, path } = snapshot.record;
    const { deleteType, deleteVersions } = snapshot.adapterOptions;

    if (!backend || !path) {
      throw new Error('The request to delete or undelete is missing required attributes.');
    }

    switch (deleteType) {
      case 'delete-latest-version':
        return this.ajax(this._urlForSecret(backend, path), 'DELETE');
      case 'delete-specific-version':
        return this.ajax(this._urlForSecret(backend, path), 'POST', {
          data: { versions: deleteVersions },
        });
      case 'destroy-specific-version':
        return this.ajax(`${this.buildURL()}/${encodePath(backend)}/destroy/${encodePath(path)}`, 'PUT', {
          data: { versions: deleteVersions },
        });
      case 'destroy-everything':
        return this.ajax(`${this.buildURL()}/${encodePath(backend)}/metadata/${encodePath(path)}`, 'DELETE');
      case 'undelete-specific-version':
        return this.ajax(`${this.buildURL()}/${encodePath(backend)}/undelete/${encodePath(path)}`, 'POST', {
          data: { versions: deleteVersions },
        });
      default:
        assert(
          'deletType must be one of delete-latest-version, delete-specific-version, destroy-specific-version, destroy-everything, undelete-specific-version.'
        );
    }
  }
}
