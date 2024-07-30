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
      resp.id = backend;
      return resp;
    });
  }

  createOrUpdate(store, type, snapshot) {
    const { data } = snapshot.adapterOptions;
    const path = encodePath(snapshot.id);
    return this.ajax(`/v1/${path}/config/root`, 'POST', { data });
  }

  createRecord() {
    return this.createOrUpdate(...arguments);
  }

  updateRecord() {
    return this.createOrUpdate(...arguments);
  }
}
