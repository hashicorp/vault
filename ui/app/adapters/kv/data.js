/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import ApplicationAdapter from '../application';
import { kvDataPath, kvDestroyPath, kvMetadataPath, kvUndeletePath } from 'vault/utils/kv-path';
import { assert } from '@ember/debug';

export default class KvDataAdapter extends ApplicationAdapter {
  namespace = 'v1';

  _url(backend, fullPath) {
    return `${this.buildURL()}/${fullPath}`;
  }

  createRecord(store, type, snapshot) {
    const { backend, path } = snapshot.record;
    const url = this._url(kvDataPath(backend, path));
    return this.ajax(url, 'POST', { data: this.serialize(snapshot) });
  }

  findRecord(store, type, id) {
    // ID is the full path for the data (including version)
    return this.ajax(this._url(id), 'GET').then((resp) => {
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
        return this.ajax(this._url(kvDataPath(backend, path)), 'DELETE');
      case 'delete-specific-version':
        return this.ajax(this._url(kvDataPath(backend, path)), 'POST', {
          data: { versions: deleteVersions },
        });
      case 'destroy-specific-version':
        return this.ajax(this._url(kvDestroyPath(backend, path)), 'PUT', {
          data: { versions: deleteVersions },
        });
      case 'destroy-everything':
        return this.ajax(this._url(kvMetadataPath(backend, path)), 'DELETE');
      case 'undelete-specific-version':
        return this.ajax(this._url(kvUndeletePath(backend, path)), 'POST', {
          data: { versions: deleteVersions },
        });
      default:
        assert(
          'deletType must be one of delete-latest-version, delete-specific-version, destroy-specific-version, destroy-everything, undelete-specific-version.'
        );
    }
  }
}
