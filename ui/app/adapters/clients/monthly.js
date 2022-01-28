import Application from '../application';

export default Application.extend({
  queryRecord() {
    let url = `${this.buildURL()}/internal/counters/activity/monthly`;
    // Query has startTime defined. The API will return the endTime if none is provided.
    return this.ajax(url, 'GET').then((resp) => {
      let response = resp || {};
      // if the response is a 204 it has no request id (ARG TODO test that it returns a 204)
      response.id = response.request_id || 'no-data';
      return response;
    });
  },
});
