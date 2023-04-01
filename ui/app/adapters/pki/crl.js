/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { encodePath } from 'vault/utils/path-encoding-helpers';
import ApplicationAdapter from '../application';

export default class PkiCrlAdapter extends ApplicationAdapter {
  namespace = 'v1';

  _url(backend) {
    return `${this.buildURL()}/${encodePath(backend)}/config/crl`;
  }

  findRecord(store, type, backend) {
    return this.ajax(this._url(backend), 'GET').then((resp) => {
      return resp.data;
    });
  }
}
