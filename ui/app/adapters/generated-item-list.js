import { assign } from '@ember/polyfills';
import ApplicationAdapter from './application';

export default ApplicationAdapter.extend({
  namespace: 'v1',
  urlForItem() {},

  fetchByQuery(store, query) {
    const { id, type } = query;
    return this.ajax(this.urlForItem(id, type), 'GET', { data: { list: true } }).then(resp => {
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
});
