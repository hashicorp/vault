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

    return this.ajax(url, 'POST', { data: this.serialize(snapshot) }).then((resp) => {
      resp.id = id;
      return resp;
    });
  }

  findRecord(store, type, id) {
    return this.ajax(this._url(id), 'GET').then((resp) => {
      resp.id = id;
      return resp;
    });
  }

  query(store, type, query) {
    const { backend } = query;
    return super.query(store, type, query).then((resp) => {
      // this is required to properly build the model in normalizeResponse
      resp.backend = backend;
      return resp;
    });
  }
}
