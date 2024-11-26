/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from '../application';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default class CubbyholeConfig extends ApplicationAdapter {
  namespace = 'v1';

  queryRecord() {
    return this.ajax(`${this.buildURL()}/cubbyhole/config/scope`, 'GET').then((resp) => {
      return {
        ...resp,
        id: 'cubbyhole',
      };
    });
  }

  updateRecord(store, type, snapshot) {
    const serializer = store.serializerFor(type.modelName);
    const data = serializer.serialize(snapshot);
    return this.ajax(`${this.buildURL()}/cubbyhole/config/scope`, 'POST', { data }).then((resp) => {
      // ember data requires an id on the response
      return {
        ...resp,
        id: 'cubbyhole',
      };
    });
  }
}
