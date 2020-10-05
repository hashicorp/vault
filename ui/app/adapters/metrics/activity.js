import Application from '../application';

export default Application.extend({
  pathForType() {
    return 'internal/counters/activity';
  },
  queryRecord(store, type, query) {
    const url = this.urlForQuery(null, type);
    // API accepts start and end as query params
    return this.ajax(url, 'GET', { data: query }).then(resp => {
      resp.id = resp.request_id;
      return resp;
    });
  },
});
