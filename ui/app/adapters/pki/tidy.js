/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { encodePath } from 'vault/utils/path-encoding-helpers';
import ApplicationAdapter from '../application';

export default class PkiUrlsAdapter extends ApplicationAdapter {
  namespace = 'v1';

  _url(backend) {
    return `${this.buildURL()}/${encodePath(backend)}/tidy`;
  }

  createRecord(store, type, snapshot) {
    const url = this._url(snapshot.record.backend);
    return this.ajax(url, 'POST', { data: this.serialize(snapshot) });
  }
}
