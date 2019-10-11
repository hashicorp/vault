import { assign } from '@ember/polyfills';
import ApplicationAdapter from './application';

export default ApplicationAdapter.extend({
  namespace: 'v1',
  urlForItem() {},
  optionsForQuery(/* id */) {
    return {
      data: {
        list: true,
      },
    };
  },

  fetchByQuery(store, query) {
    const { id, type } = query;
    return this.ajax(this.urlForItem(id, type), 'GET', this.optionsForQuery(id)).then(resp => {
      const data = {
        id,
        name: id,
        method: id,
      };

      return assign({}, resp, data);
    });
  },

  query(store, type, query) {
    return this.fetchByQuery(store, query);
  },

  queryRecord(store, type, query) {
    return this.fetchByQuery(store, query);
  },
});
