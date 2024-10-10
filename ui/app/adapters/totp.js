/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { resolve, allSettled } from 'rsvp';
import { isEmpty } from '@ember/utils';
import ApplicationAdapter from './application';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default ApplicationAdapter.extend({
  namespace: 'v1',

  createOrUpdate(store, type, snapshot, requestType) {
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
  },

  createRecord() {
    return this.createOrUpdate(...arguments);
  },

  updateRecord() {
    return this.createOrUpdate(...arguments, 'update');
  },

  deleteRecord(store, type, snapshot) {
    const { id } = snapshot;
    return this.ajax(this.urlForKey(snapshot.record.backend, id), 'DELETE');
  },

  pathForType() {
    return 'keys';
  },

  urlForKey(backend, id) {
    let url = `${this.buildURL()}/${encodePath(backend)}/${this.pathForType()}/`;

    if (!isEmpty(id)) {
      url = url + encodePath(id);
    }

    return url;
  },

  optionsForQuery(id, action) {
    const data = {};

    if (action === 'query') {
      data.list = true;
    }

    return { data };
  },

  fetchByQuery(query, action) {
    const { id, backend } = query;
    return this.ajax(this.urlForKey(backend, id), 'GET', this.optionsForQuery(id, action)).then((resp) => {
      resp.id = id;
      resp.backend = backend;
      return resp;
    });
  },

  query(store, type, query) {
    return this.fetchByQuery(query, 'query');
  },

  queryRecord(store, type, query) {
    return this.fetchByQuery(query, 'queryRecord');
  },

  generateCode(backend, id) {
    return this.ajax(`${this.buildURL()}/${encodePath(backend)}/code/${id}`, 'GET').then((res) => {
      return res.data;
    });
  },
});
