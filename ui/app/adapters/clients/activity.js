import ApplicationAdapter from '../application';
import { getUnixTime } from 'date-fns';

export default class ActivityAdapter extends ApplicationAdapter {
  // javascript localizes new Date() objects so use Date.UTC() then convert to unix
  // from the selected month index and year
  formatQueryParams({ start_time, end_time }) {
    // time params from the backend are formatted as a zulu timestamp
    start_time = start_time.timestamp || getUnixTime(Date.UTC(start_time.year, start_time.monthIdx, 1));
    // day=0 for Date.UTC() returns the last day of the month before
    // increase monthIdx by one to get last day of queried month
    end_time = end_time.timestamp || getUnixTime(Date.UTC(end_time.year, end_time.monthIdx + 1, 0));
    return { start_time, end_time };
  }

  queryRecord(store, type, query) {
    const url = `${this.buildURL()}/internal/counters/activity`;
    const queryParams =
      typeof query.start_time === 'string' && typeof query.end_time === 'string'
        ? query
        : this.formatQueryParams(query);
    if (queryParams) {
      return this.ajax(url, 'GET', { data: queryParams }).then((resp) => {
        const response = resp || {};
        response.id = response.request_id || 'no-data';
        return response;
      });
    }
  }
}
