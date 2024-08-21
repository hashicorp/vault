/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from '../application';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default class AwsRootConfig extends ApplicationAdapter {
  namespace = 'v1';

  queryRecord(store, type, query) {
    const { backend } = query;
    return this.ajax(`${this.buildURL()}/${encodePath(backend)}/config/root`, 'GET').then((resp) => {
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
    return this.ajax(`${this.buildURL()}/${backend}/config/root`, 'POST', { data }).then((resp) => {
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
}
