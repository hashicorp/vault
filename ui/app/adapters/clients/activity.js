import ApplicationAdapter from '../application';
import { formatRFC3339 } from 'date-fns';

export default class ActivityAdapter extends ApplicationAdapter {
  // format when params are objects { timestamp: null, monthIdx: 0, year: 2022 }
  formatQueryParams({ start_time, end_time }) {
    // javascript localizes new Date() objects, but we don't want a timezone attached to keep metrics data in UTC
    // hard code 10th and 20th and backend will convert to first or end of month respectively
    start_time = start_time.timestamp || formatRFC3339(new Date(start_time.year, start_time.monthIdx, 10));
    end_time = end_time.timestamp || formatRFC3339(new Date(end_time.year, end_time.monthIdx, 20));

    return { start_time, end_time };
  }

  queryRecord(store, type, query) {
    const url = `${this.buildURL()}/internal/counters/activity`;
    let queryParams =
      typeof query.start_time === 'string' && typeof query.end_time === 'string'
        ? query
        : this.formatQueryParams(query);
    if (queryParams) {
      return this.ajax(url, 'GET', { data: queryParams }).then((resp) => {
        let response = resp || {};
        response.id = response.request_id || 'no-data';
        if (response.id === 'no-data') {
          // add queryParams to return user's queried date range without data
          response = { ...response, ...queryParams };
        }
        return response;
      });
    }
  }
}
