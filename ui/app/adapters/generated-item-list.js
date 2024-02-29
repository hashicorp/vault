/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from './application';
import { task } from 'ember-concurrency';
import { service } from '@ember/service';

export default class GeneratedItemAdapter extends ApplicationAdapter {
  @service store;
  namespace = 'v1';
  urlForItem() {}
  dynamicApiPath = '';

  @task
  *getDynamicApiPath(id) {
    // TODO: remove yield at some point.
    const result = yield this.store.peekRecord('auth-method', id);
    this.dynamicApiPath = result.apiPath;
    return;
  }

  @task
  *fetchByQuery(store, query, isList) {
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
  }

  query(store, type, query) {
    return this.fetchByQuery.perform(store, query, true);
  }

  queryRecord(store, type, query) {
    return this.fetchByQuery.perform(store, query);
  }
}
