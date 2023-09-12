/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import ApplicationAdapter from '../application';
import { kvDataPath, kvDeletePath, kvDestroyPath, kvMetadataPath, kvUndeletePath } from 'vault/utils/kv-path';
import { assert } from '@ember/debug';
import ControlGroupError from 'vault/lib/control-group-error';

export default class KvDataAdapter extends ApplicationAdapter {
  namespace = 'v1';

  _url(fullPath) {
    return `${this.buildURL()}/${fullPath}`;
  }

  createRecord(store, type, snapshot) {
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

  fetchWrapInfo(query) {
    const { backend, path, version, wrapTTL } = query;
    const id = kvDataPath(backend, path, version);
    return this.ajax(this._url(id), 'GET', { wrapTTL }).then((resp) => resp.wrap_info);
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
        const baseResponse = { id, backend, path, version };
        const errorCode = errorOrResponse.httpStatus;
        // if it's a legitimate error - throw it!
        if (errorOrResponse instanceof ControlGroupError) {
          throw errorOrResponse;
        }

        if (errorCode === 403) {
          return {
            data: {
              ...baseResponse,
              fail_read_error_code: errorCode,
            },
          };
        }

        if (errorOrResponse.data) {
          // in the case of a deleted/destroyed secret the API returns a 404 because { data: null }
          // however, there could be a metadata block with important information like deletion_time
          // handleResponse below checks 404 status codes for metadata and updates the code to 200 if it exists.
          // we still end up in the good ol' catch() block, but instead of a 404 adapter error we've "caught"
          // the metadata that sneakily tried to hide from us
          return {
            ...errorOrResponse,
            data: {
              ...baseResponse,
              ...errorOrResponse.data, // includes the { metadata } key we want
            },
          };
        }

        // If we get here, it's probably a 404 because it doesn't exist
        throw errorOrResponse;
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
      case 'delete-version':
        return this.ajax(this._url(kvDeletePath(backend, path)), 'POST', {
          data: { versions: deleteVersions },
        });
      case 'destroy':
        return this.ajax(this._url(kvDestroyPath(backend, path)), 'PUT', {
          data: { versions: deleteVersions },
        });
      case 'undelete':
        return this.ajax(this._url(kvUndeletePath(backend, path)), 'POST', {
          data: { versions: deleteVersions },
        });
      case 'destroy-all-versions':
        return this.ajax(this._url(kvMetadataPath(backend, path)), 'DELETE');
      default:
        assert(
          'deleteType must be one of delete-latest-version, delete-version, destroy, undelete, or destroy-all-versions.'
        );
    }
  }

  handleResponse(status, headers, payload, requestData) {
    // after deleting a secret version, data is null and the API returns a 404
    // but there could be relevant metadata
    if (status === 404 && payload.data?.metadata) {
      return super.handleResponse(200, headers, payload, requestData);
    }
    return super.handleResponse(...arguments);
  }
}
