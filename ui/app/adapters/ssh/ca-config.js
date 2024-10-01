/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from '../application';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default class SshCaConfig extends ApplicationAdapter {
  namespace = 'v1';

  queryRecord(store, type, query) {
    const { backend } = query;
    return this.ajax(`${this.buildURL()}/${encodePath(backend)}/config/ca`, 'GET').then((resp) => {
      resp.id = backend;
      resp.backend = backend;
      return resp;
    });
  }

  createOrUpdate(store, type, snapshot) {
    const serializer = store.serializerFor(type.modelName);
    const data = serializer.serialize(snapshot);
    const backend = snapshot.record.backend;
    return this.ajax(`${this.buildURL()}/${backend}/config/ca`, 'POST', { data }).then((resp) => {
      // ember data requires an id on the response
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

  deleteRecord(store, type, snapshot) {
    const backend = snapshot.record.backend;
    return this.ajax(`${this.buildURL()}/${backend}/config/ca`, 'DELETE');
  }
}
