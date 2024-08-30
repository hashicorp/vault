/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from '../application';
import { kvDataPath, kvDeletePath, kvDestroyPath, kvSubkeysPath, kvUndeletePath } from 'vault/utils/kv-path';
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

  fetchSubkeys(backend, path, query) {
    const url = this._url(kvSubkeysPath(backend, path, query));
    return (
      this.ajax(url, 'GET')
        .then((resp) => resp.data)
        // deleted/destroyed secret versions throw an error
        // but still have metadata that we want to return
        .catch((errorOrResponse) => {
          return this.parseErrorOrResponse(errorOrResponse, { backend, path }, true);
        })
    );
  }

  fetchWrapInfo(query) {
    const { backend, path, version, wrapTTL } = query;
    const id = kvDataPath(backend, path, version);
    return this.ajax(this._url(id), 'GET', { wrapTTL }).then((resp) => resp.wrap_info);
  }

  // patching a secret happens without retrieving the ember data model
  // so we use a custom method instead of updateRecord
  patchSecret(backend, path, patchData, version) {
    const url = this._url(kvDataPath(backend, path));
    const data = {
      options: { cas: version },
      data: patchData,
    };
    return this.ajax(url, 'PATCH', { data });
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
        return this.parseErrorOrResponse(errorOrResponse, { id, backend, path, version });
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
      default:
        assert('deleteType must be one of delete-latest-version, delete-version, destroy, or undelete.');
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

  parseErrorOrResponse(errorOrResponse, secretDataBaseResponse, isSubkeys = false) {
    // if it's a legitimate error - throw it!
    if (errorOrResponse instanceof ControlGroupError) {
      throw errorOrResponse;
    }

    const errorCode = errorOrResponse.httpStatus;
    if (errorCode === 403) {
      return {
        data: {
          ...secretDataBaseResponse,
          fail_read_error_code: errorCode,
        },
      };
    }

    // in the case of a deleted/destroyed secret the API returns a 404 because { data: null }
    // however, there could be a metadata block with important information like deletion_time
    // handleResponse below checks 404 status codes for metadata and updates the code to 200 if it exists.
    // we still end up in the good ol' catch() block, but instead of a 404 adapter error we've "caught"
    // the metadata that sneakily tried to hide from us
    if (errorOrResponse.data) {
      // subkeys response doesn't correspond to a model, no need to include base response
      if (isSubkeys) return errorOrResponse.data;

      return {
        ...errorOrResponse,
        data: {
          ...secretDataBaseResponse,
          ...errorOrResponse.data, // includes the { metadata } key we want
        },
      };
    }

    // If we get here, it's probably a 404 because it doesn't exist
    throw errorOrResponse;
  }
}
