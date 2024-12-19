/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from '../application';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default class GcpConfig extends ApplicationAdapter {
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
}
