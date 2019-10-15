import { assign } from '@ember/polyfills';
import ApplicationAdapter from './application';

export default ApplicationAdapter.extend({
  namespace: 'v1',
  urlForItem() {},

  fetchByQuery(store, query, isList) {
    const { id, authMethodPath } = query;
    let data = {};
    if (isList) {
      data.list = true;
    }

    return this.ajax(this.urlForItem(id, isList), 'GET', { data }).then(resp => {
      const data = {
        id,
        name: id,
        method: authMethodPath,
      };

      return assign({}, resp, data);
    });
  },

  query(store, type, query) {
    return this.fetchByQuery(store, query, true);
  },

  queryRecord(store, type, query) {
    return this.fetchByQuery(store, query);
  },
});
