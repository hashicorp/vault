/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from './application';
import { encodePath } from 'vault/utils/path-encoding-helpers';
import { isEmpty } from '@ember/utils';

export default class TotpKeyAdapter extends ApplicationAdapter {
  namespace = 'v1';

  // TOTP keys can only be created, so no need for an update method
  createRecord(store, type, snapshot) {
    const { name, backend } = snapshot.record;
    const serializer = store.serializerFor(type.modelName);
    const data = serializer.serialize(snapshot);
    const url = this.urlForKey(backend, name);

    return this.ajax(url, 'POST', { data }).then((resp) => {
      // Ember data doesn't like 204 responses except for DELETE method
      const response = resp || { data: {} };
      response.data.id = name;
      return response;
    });
  }

  deleteRecord(store, type, snapshot) {
    const { id } = snapshot;
    return this.ajax(this.urlForKey(snapshot.record.backend, id), 'DELETE');
  }

  urlForKey(backend, id) {
    let url = `${this.buildURL()}/${encodePath(backend)}/keys`;

    if (!isEmpty(id)) {
      url = `${url}/${encodePath(id)}`;
    }

    return url;
  }

  query(store, type, query) {
    const { backend } = query;
    return this.ajax(this.urlForKey(backend), 'GET', { data: { list: true } }).then((resp) => {
      resp.backend = backend;
      return resp;
    });
  }

  queryRecord(store, type, query) {
    const { id, backend } = query;
    return this.ajax(this.urlForKey(backend, id), 'GET').then((resp) => {
      resp.id = id;
      resp.backend = backend;
      return resp;
    });
  }

  generateCode(backend, id) {
    return this.ajax(`${this.buildURL()}/${encodePath(backend)}/code/${id}`, 'GET').then((res) => {
      return res.data;
    });
  }
}
