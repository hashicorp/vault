import Application from '../application';

export default Application.extend({
  pathForType() {
    return 'internal/counters/activity';
  },
  queryRecord(store, type, query) {
    let url = this.urlForQuery(null, type);
    // ARG query has startTime:
    // API accepts start and end as query params
    return this.ajax(url, 'GET', { data: query }).then((resp) => {
      let response = resp || {};
      // if the response is a 204 it has no request id
      response.id = response.request_id || 'no-data';
      return response;
    });
  },
});
