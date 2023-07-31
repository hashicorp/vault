/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import ApplicationAdapter from '../application';
import AdapterError from '@ember-data/adapter/error';
import { kvDataPath, kvDestroyPath, kvMetadataPath, kvUndeletePath } from 'vault/utils/kv-path';
import { assert } from '@ember/debug';

export default class KvDataAdapter extends ApplicationAdapter {
  namespace = 'v1';

  _url(fullPath) {
    return `${this.buildURL()}/${fullPath}`;
  }

  _createOrUpdate(snapshot) {
    const { backend, path } = snapshot.record;
    const url = this._url(kvDataPath(backend, path));

    return this.ajax(url, 'POST', { data: this.serialize(snapshot) }).then((res) => {
      return {
        data: {
          id: kvDataPath(backend, path, res.data.version),
          backend,
          path,
          ...res.data,
        },
      };
    });
  }

  createRecord(store, type, snapshot) {
    return this._createOrUpdate(snapshot);
  }

  queryRecord(store, type, query) {
    const { backend, path, version } = query;
    // ID is the full path for the data (including version)
    let id = kvDataPath(backend, path, version);
    return this.ajax(this._url(id), 'GET')
      .then((resp) => {
        // if no version is queried, add version from response to ID
        // otherwise duplicate ember data models will exist in store
        // (one with an ID that includes the version and one without)
        if (!version) {
          id = kvDataPath(backend, path, resp.data.metadata.version);
        }
        return {
          ...resp,
          data: {
            id,
            backend,
            path,
            ...resp.data,
          },
        };
      })
      .catch((errorOrResponse) => {
        // if it's a legitimate error - throw it!
        if (errorOrResponse instanceof AdapterError) {
          throw errorOrResponse;
        }
        // in the case of a deleted/destroyed secret the API returns a 404 because { data: null }
        // however, there could be a metadata block with important information like deletion_time
        // handleResponse below checks 404 status codes for metadata and updates the code to 200 if it exists.
        // we still end up in the good ol' catch() block, but instead of a 404 adapter error we've "caught"
        // the metadata that sneakily tried to hide from us
        return {
          ...errorOrResponse,
          data: {
            ...errorOrResponse.data, // includes the { metadata } key we want
            id,
            backend,
            path,
          },
        };
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

  handleResponse(status, headers, payload, requestData) {
    // after deleting a secret version, data is null and the API returns a 404
    // but there could be relevant metadata
    if (status === 404 && payload.data.metadata) {
      return super.handleResponse(200, headers, payload, requestData);
    }
    return super.handleResponse(...arguments);
  }
}
