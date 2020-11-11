import { assign } from '@ember/polyfills';
import ApplicationAdapter from './application';
import { task } from 'ember-concurrency';

export default ApplicationAdapter.extend({
  namespace: 'v1',
  urlForItem() {},
  dynamicApiPath: '',
  getDynamicApiPath: task(function*(id) {
    let result = yield this.store.peekRecord('auth-method', id);
    this.dynamicApiPath = result.apiPath;
    return;
  }),

  fetchByQuery: task(function*(store, query, isList) {
    const { id } = query;
    let data = {};
    if (isList) {
      data.list = true;
      yield this.getDynamicApiPath.perform(id);
    }

    return this.ajax(this.urlForItem(id, isList, this.dynamicApiPath), 'GET', { data }).then(resp => {
      const data = {
        id,
        method: id,
      };
      return assign({}, resp, data);
    });
  }),

  query(store, type, query) {
    return this.fetchByQuery.perform(store, query, true);
  },

  queryRecord(store, type, query) {
    return this.fetchByQuery.perform(store, query);
  },
});
