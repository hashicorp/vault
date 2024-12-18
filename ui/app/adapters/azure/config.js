/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from '../application';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default class AzureConfig extends ApplicationAdapter {
  namespace = 'v1';

  _url(backend) {
    return `${this.buildURL()}/${encodePath(backend)}/config`;
  }

  queryRecord(store, type, query) {
    const { backend } = query;
    return this.ajax(this._url(backend), 'GET').then((resp) => {
      return {
        ...resp,
        id: backend,
        backend,
      };
    });
  }

  createOrUpdate(store, type, snapshot) {
    const serializer = store.serializerFor(type.modelName);
    const data = serializer.serialize(snapshot);
    const backend = snapshot.record.backend;
    return this.ajax(this._url(backend), 'POST', { data }).then((resp) => {
      return {
        ...resp,
        id: backend,
      };
    });
  }

  createRecord() {
    return this.createOrUpdate(...arguments);
  }

  updateRecord() {
    return this.createOrUpdate(...arguments);
  }
}
