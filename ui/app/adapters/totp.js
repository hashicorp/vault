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

  createRecord(store, type, snapshot) {
    const serializer = store.serializerFor(type.modelName);
    const data = serializer.serialize(snapshot);
    const { id } = snapshot;
    const name = snapshot.record.name;
    return this.ajax(this.urlForKey(snapshot.attr('backend'), name || id), 'POST', { data }).then(() => {
      data.id = name || id;
      return data;
    });
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
});
