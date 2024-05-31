/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { allSettled } from 'rsvp';
import ApplicationAdapter from './application';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default ApplicationAdapter.extend({
  namespace: 'v1',

  createOrUpdate(store, { modelName }, snapshot) {
    const { backend, name, type } = snapshot.record;
    const serializer = store.serializerFor(modelName);
    const data = serializer.serialize(snapshot);
    const url = this.urlForTransformations(backend, name, type);

    return this.ajax(url, 'POST', { data }).then((resp) => {
      const response = resp || {};
      response.id = name;
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
    return this.ajax(this.urlForTransformations(snapshot.record.backend, id), 'DELETE');
  },

  pathForType() {
    return 'transform';
  },

  urlForTransformations(backend, id, type) {
    const base = `${this.buildURL()}/${encodePath(backend)}`;
    // when type exists, transformations is plural
    const url = type ? `${base}/transformations/${type}` : `${base}/transformation`;
    if (id) return `${url}/${encodePath(id)}`;
    return url;
  },

  optionsForQuery(id) {
    const data = {};
    if (!id) {
      data['list'] = true;
    }
    return { data };
  },

  fetchByQuery(store, query) {
    const { id, backend } = query;
    const queryAjax = this.ajax(this.urlForTransformations(backend, id), 'GET', this.optionsForQuery(id));

    return allSettled([queryAjax]).then((results) => {
      // query result 404, so throw the adapterError
      if (!results[0].value) {
        throw results[0].reason;
      }
      const resp = {
        id,
        name: id,
        backend,
        data: {},
      };

      results.forEach((result) => {
        if (result.value) {
          let d = result.value.data;
          if (d.templates) {
            // In Transformations data goes up as "template", but comes down as "templates"
            // To keep the keys consistent we're translating here
            d = {
              ...d,
              template: d.templates,
            };
            delete d.templates;
          }
          resp.data = { ...resp.data, ...d };
        }
      });
      return resp;
    });
  },

  query(store, type, query) {
    return this.fetchByQuery(store, query);
  },

  queryRecord(store, type, query) {
    return this.fetchByQuery(store, query);
  },
});
