/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { encodePath } from 'vault/utils/path-encoding-helpers';
import ApplicationAdapter from '../application';

export default class PkiUrlsAdapter extends ApplicationAdapter {
  namespace = 'v1';

  _url(backend) {
    return `${this.buildURL()}/${encodePath(backend)}/config/urls`;
  }

  createOrUpdate(store, type, snapshot) {
    const data = snapshot.serialize();
    return this.ajax(this._url(snapshot.record.id), 'POST', { data });
  }

  createRecord() {
    return this.createOrUpdate(...arguments);
  }

  updateRecord() {
    return this.createOrUpdate(...arguments);
  }

  urlForFindRecord(id) {
    return this._url(id);
  }
}
