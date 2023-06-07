/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import ApplicationAdapter from '../application';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default class KvMetadataAdapter extends ApplicationAdapter {
  namespace = 'v1';

  _urlForMetadata(backend, path) {
    // path is "kv-test/2/my-secret"
    return `${this.buildURL()}/${encodePath(backend)}/metadata/${encodePath(path)}`;
  }

  createRecord(store, type, snapshot) {
    const backend = snapshot.record.backend;
    const path = snapshot.attr('path');
    const url = this._urlForMetadata(backend, path);

    return this.ajax(url, 'POST', { data: this.serialize(snapshot) }).then((resp) => {
      resp.id = `${backend}/${path}`;
      return resp;
    });
  }

  updateRecord(store, type, snapshot) {
    const { backend, path } = snapshot.record;
    const data = this.serialize(snapshot);
    const url = this._urlForMetadata(backend, path);
    return this.ajax(url, 'POST', { data });
  }

  query(store, type, query) {
    const { path, backend } = query;
    return this.ajax(this._urlForMetadata(backend, path), 'GET');
  }

  queryRecord(store, type, query) {
    const { path, backend } = query;
    return this.ajax(this._urlForMetadata(backend, path), 'GET');
  }
}
