import Application from '../application';

export default Application.extend({
  queryRecord(store, type, query) {
    let url = `${this.buildURL()}/internal/counters/activity`;
    // Query has startTime defined. The API will return the endTime if none is provided.
    return this.ajax(url, 'GET', { data: query }).then((resp) => {
      let response = resp || {};
      // if the response is a 204 it has no request id
      response.id = response.request_id || 'no-data';
      return response;
    });
  },
  // called from components
  queryClientActivity(start_time, end_time) {
    // do not query without start_time. Otherwise returns last year data, which is not reflective of billing data.
    if (start_time) {
      let url = `${this.buildURL()}/internal/counters/activity`;
      let queryParams = {};
      if (!end_time) {
        queryParams = { data: { start_time } };
      } else {
        queryParams = { data: { start_time, end_time } };
      }
      return this.ajax(url, 'GET', queryParams).then((resp) => {
        return resp;
      });
    }
  },
});
