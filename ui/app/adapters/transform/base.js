/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import ApplicationAdapter from '../application';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default ApplicationAdapter.extend({
  namespace: 'v1',

  pathForType(type) {
    return type.replace('transform/', '');
  },

  createOrUpdate(store, type, snapshot) {
    const serializer = store.serializerFor(type.modelName);
    const data = serializer.serialize(snapshot);
    const { id } = snapshot;
    const url = this.url(snapshot.record.get('backend'), type.modelName, id);

    return this.ajax(url, 'POST', { data });
  },

  createRecord() {
    return this.createOrUpdate(...arguments);
  },

  updateRecord() {
    return this.createOrUpdate(...arguments, 'update');
  },

  deleteRecord(store, type, snapshot) {
    const { id } = snapshot;
    return this.ajax(this.url(snapshot.record.get('backend'), type.modelName, id), 'DELETE');
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
      return {
        ...resp,
        backend,
      };
    });
  },

  query(store, type, query) {
    return this.fetchByQuery(query);
  },

  queryRecord(store, type, query) {
    return this.ajax(this.url(query.backend, type.modelName, query.id), 'GET').then((result) => {
      // CBS TODO: Add name to response and unmap name <> id on models
      return {
        id: query.id,
        ...result,
      };
    });
  },
});
