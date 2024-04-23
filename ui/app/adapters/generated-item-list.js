/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from './application';
import { task } from 'ember-concurrency';
import { service } from '@ember/service';

export default ApplicationAdapter.extend({
  store: service(),
  namespace: 'v1',
  urlForItem() {},
  dynamicApiPath: '',

  getDynamicApiPath: task(function* (id) {
    // TODO: remove yield at some point.
    const result = yield this.store.peekRecord('auth-method', id);
    this.dynamicApiPath = result.apiPath;
    return;
  }),

  fetchByQuery: task(function* (store, query, isList) {
    const { id } = query;
    const data = {};
    if (isList) {
      data.list = true;
      yield this.getDynamicApiPath.perform(id);
    }

    return this.ajax(this.urlForItem(id, isList, this.dynamicApiPath), 'GET', { data }).then((resp) => {
      const data = {
        id,
        method: id,
      };
      return { ...resp, ...data };
    });
  }),

  query(store, type, query) {
    return this.fetchByQuery.perform(store, query, true);
  },

  queryRecord(store, type, query) {
    return this.fetchByQuery.perform(store, query);
  },
});
