/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import ApplicationAdapter from '../application';
import { encodePath } from 'vault/utils/path-encoding-helpers';
import { kvId } from 'vault/utils/kv-id';

export default class KvMetadataAdapter extends ApplicationAdapter {
  namespace = 'v1';

  _urlForMetadata(backend, path) {
    return `${this.buildURL()}/${encodePath(backend)}/metadata/${encodePath(path)}`;
  }

  createRecord(store, type, snapshot) {
    const { backend, path } = snapshot.record;
    const url = this._urlForMetadata(backend, path);

    return this.ajax(url, 'POST', { data: this.serialize(snapshot) }).then((resp) => {
      resp.id = kvId(backend, path, 'metadata');
      return resp;
    });
  }

  queryRecord(store, type, query) {
    const { path, backend } = query;
    return this.ajax(this._urlForMetadata(backend, path), 'GET').then((resp) => {
      resp.id = kvId(backend, path, 'metadata');
      return resp;
    });
  }

  findRecord(store, type, id) {
    return this.ajax(`${this.buildURL()}/${id}`, 'GET').then((resp) => {
      resp.id = id;
      return resp;
    });
  }
}
