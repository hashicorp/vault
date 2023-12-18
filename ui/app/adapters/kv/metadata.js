/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import ApplicationAdapter from '../application';
import { kvMetadataPath } from 'vault/utils/kv-path';

export default class KvMetadataAdapter extends ApplicationAdapter {
  namespace = 'v1';

  _url(fullPath) {
    return `${this.buildURL()}/${fullPath}`;
  }

  createRecord(store, type, snapshot) {
    const { backend, path } = snapshot.record;
    const id = kvMetadataPath(backend, path);
    const url = this._url(id);
    const data = this.serialize(snapshot);
    return this.ajax(url, 'POST', { data }).then(() => {
      return {
        id,
        data,
      };
    });
  }

  updateRecord(store, type, snapshot) {
    const { backend, path } = snapshot.record;
    const id = kvMetadataPath(backend, path);
    const url = this._url(id);
    const data = this.serialize(snapshot);
    return this.ajax(url, 'POST', { data }).then(() => {
      return {
        id,
        data,
      };
    });
  }

  query(store, type, query) {
    const { backend, pathToSecret } = query;
    // example of pathToSecret: beep/boop/
    return this.ajax(this._url(kvMetadataPath(backend, pathToSecret)), 'GET', {
      data: { list: true },
    }).then((resp) => {
      resp.backend = backend;
      resp.path = pathToSecret;
      return resp;
    });
  }

  queryRecord(store, type, query) {
    const { backend, path } = query;
    // ID is the full path for the metadata
    const id = kvMetadataPath(backend, path);
    return this.ajax(this._url(id), 'GET').then((resp) => {
      return {
        id,
        ...resp,
        data: {
          backend,
          path,
          ...resp.data,
        },
      };
    });
  }

  // This method is only called when deleting from the LIST view. Otherwise, delete on kv/data
  deleteRecord(store, type, snapshot) {
    const { backend, path, fullSecretPath } = snapshot.record;
    // fullSecretPath is used when deleting from the LIST view and is defined via the serializer
    // path is used when deleting from the metadata details view.
    return this.ajax(this._url(kvMetadataPath(backend, fullSecretPath || path)), 'DELETE');
  }
}
