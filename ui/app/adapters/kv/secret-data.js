/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import ApplicationAdapter from '../application';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default class KvSecretDataAdapter extends ApplicationAdapter {
  // Everything on the data endpoint
  namespace = 'v1';

  getURL(backend, name) {
    const base = `${this.buildURL()}/${encodePath(backend)}/data/`;
    return name ? `${base}${encodePath(name)}` : base;
  }

  createRecord() {
    return this._saveRecord(...arguments);
  }

  _saveRecord(store, { modelName }, snapshot) {
    const data = store.serializerFor(modelName).serialize(snapshot);
    const url = this.getURL(snapshot.attr('backend'), data.path);
    // delete path and backend from the payload
    delete data.path;
    delete data.backend;
    return this.ajax(url, 'POST', { data }).then(() => data);
  }
}
