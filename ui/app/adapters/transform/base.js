/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from '../application';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default ApplicationAdapter.extend({
  namespace: 'v1',

  pathForType(type) {
    return type.replace('transform/', '');
  },

  createOrUpdate(store, type, snapshot) {
    const { backend, name } = snapshot.record;
    const serializer = store.serializerFor(type.modelName);
    const data = serializer.serialize(snapshot);
    const url = this.url(backend, type.modelName, name);
    return this.ajax(url, 'POST', { data }).then((resp) => {
      // Ember data doesn't like 204 responses except for DELETE method
      const response = resp || { data: {} };
      response.data.name = name;
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
    return this.ajax(this.url(snapshot.record.backend, type.modelName, id), 'DELETE');
  },

  url(backend, modelType, id) {
    const type = this.pathForType(modelType);
    const url = `/${this.namespace}/${encodePath(backend)}/${encodePath(type)}`;
    if (id) {
      return `${url}/${encodePath(id)}`;
    }
    return url + '?list=true';
  },

  fetchByQuery(query) {
    const { backend, modelName, id } = query;
    return this.ajax(this.url(backend, modelName, id), 'GET').then((resp) => {
      // The API response doesn't explicitly include the name/id, so add it here
      return {
        ...resp,
        backend,
        id,
        name: id,
      };
    });
  },

  query(store, type, query) {
    return this.fetchByQuery(query);
  },

  queryRecord(store, type, query) {
    return this.ajax(this.url(query.backend, type.modelName, query.id), 'GET').then((result) => {
      return {
        id: query.id,
        name: query.id,
        ...result,
      };
    });
  },
});
