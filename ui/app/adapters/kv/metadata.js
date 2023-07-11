/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { encodePath } from 'vault/utils/path-encoding-helpers';
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

  // TODO: replace this with raw request for metadata request?
  query(store, type, query) {
    const { backend } = query;
    return this.ajax(this._url(`${encodePath(backend)}/metadata?list=1`), 'GET').then((resp) => {
      resp.backend = backend;
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
}
