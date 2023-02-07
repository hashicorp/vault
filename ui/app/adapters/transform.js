/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { assign } from '@ember/polyfills';
import { allSettled } from 'rsvp';
import ApplicationAdapter from './application';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default ApplicationAdapter.extend({
  namespace: 'v1',

  createOrUpdate(store, type, snapshot) {
    const serializer = store.serializerFor(type.modelName);
    const data = serializer.serialize(snapshot);
    const { id } = snapshot;
    const url = this.urlForTransformations(snapshot.record.get('backend'), id);

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
    return this.ajax(this.urlForTransformations(snapshot.record.get('backend'), id), 'DELETE');
  },

  pathForType() {
    return 'transform';
  },

  urlForTransformations(backend, id) {
    let url = `${this.buildURL()}/${encodePath(backend)}/transformation`;
    if (id) {
      url = url + '/' + encodePath(id);
    }
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
      // query result 404d, so throw the adapterError
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
          resp.data = assign({}, resp.data, d);
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
