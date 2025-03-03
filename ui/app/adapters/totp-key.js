/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from './application';
import { encodePath } from 'vault/utils/path-encoding-helpers';
import { isEmpty } from '@ember/utils';

export default class TotpKeyAdapter extends ApplicationAdapter {
  namespace = 'v1';

  createOrUpdate(store, type, snapshot, requestType) {
    // TODO, unsure why request type is needed, but updates currently fail
    const { name, backend } = snapshot.record;
    const serializer = store.serializerFor(type.modelName);
    const data = serializer.serialize(snapshot, requestType);
    const url = this.urlForKey(backend, name);

    return this.ajax(url, 'POST', { data }).then((resp) => {
      // Ember data doesn't like 204 responses except for DELETE method
      const response = resp || { data: {} };
      response.data.id = name;
      return response;
    });
  }

  createRecord() {
    return this.createOrUpdate(...arguments);
  }

  updateRecord() {
    return this.createOrUpdate(...arguments, 'update');
  }

  deleteRecord(store, type, snapshot) {
    const { id } = snapshot;
    return this.ajax(this.urlForKey(snapshot.record.backend, id), 'DELETE');
  }

  urlForKey(backend, id) {
    let url = `${this.buildURL()}/${encodePath(backend)}/keys/`;

    if (!isEmpty(id)) {
      url = url + encodePath(id);
    }

    return url;
  }

  optionsForQuery(action) {
    const data = {};

    if (action === 'query') {
      data.list = true;
    }

    return { data };
  }

  fetchByQuery(query, action) {
    const { id, backend } = query;
    return this.ajax(this.urlForKey(backend, id), 'GET', this.optionsForQuery(action)).then((resp) => {
      resp.id = id;
      resp.backend = backend;
      return resp;
    });
  }

  query(store, type, query) {
    return this.fetchByQuery(query, 'query');
  }

  queryRecord(store, type, query) {
    return this.fetchByQuery(query, 'queryRecord');
  }

  generateCode(backend, id) {
    return this.ajax(`${this.buildURL()}/${encodePath(backend)}/code/${id}`, 'GET').then((res) => {
      return res.data;
    });
  }
}
