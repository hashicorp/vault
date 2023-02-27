/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import ApplicationAdapter from 'vault/adapters/application';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default class KubernetesConfigAdapter extends ApplicationAdapter {
  namespace = 'v1';

  getURL(backend, path = 'config') {
    return `${this.buildURL()}/${encodePath(backend)}/${path}`;
  }
  urlForUpdateRecord(name, modelName, snapshot) {
    return this.getURL(snapshot.attr('backend'));
  }
  urlForDeleteRecord(backend) {
    return this.getURL(backend);
  }

  queryRecord(store, type, query) {
    const { backend } = query;
    return this.ajax(this.getURL(backend), 'GET').then((resp) => {
      resp.backend = backend;
      return resp;
    });
  }
  createRecord() {
    return this._saveRecord(...arguments);
  }
  updateRecord() {
    return this._saveRecord(...arguments);
  }
  _saveRecord(store, { modelName }, snapshot) {
    const data = store.serializerFor(modelName).serialize(snapshot);
    const url = this.getURL(snapshot.attr('backend'));
    return this.ajax(url, 'POST', { data }).then(() => data);
  }
  checkConfigVars(backend) {
    return this.ajax(`${this.getURL(backend, 'check')}`, 'GET');
  }
}
