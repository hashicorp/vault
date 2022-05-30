import ApplicationAdapter from '../application';

export default class MonthlyAdapter extends ApplicationAdapter {
  queryRecord() {
    let url = `${this.buildURL()}/internal/counters/activity/monthly`;
    // Query has startTime defined. The API will return the endTime if none is provided.
    return this.ajax(url, 'GET').then((resp) => {
      let response = resp || {};
      response.id = response.request_id || 'no-data';
      return response;
    });
  }
}
