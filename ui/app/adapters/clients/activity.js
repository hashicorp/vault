import Application from '../application';
import { formatRFC3339 } from 'date-fns';

export default Application.extend({
  // Since backend converts the timezone to UTC, sending the first (1) as start or end date can cause the month to change.
  // To mitigate this impact of timezone conversion, hard coding the dates to avoid month change.
  formatTimeParams(query) {
    let { start_time, end_time } = query;
    // check if it's an array, if it is, it's coming from an action like selecting a new startTime or new EndTime
    if (Array.isArray(start_time)) {
      let startYear = Number(start_time[0]);
      let startMonth = Number(start_time[1]);
      start_time = formatRFC3339(new Date(startYear, startMonth, 10));
    }
    if (end_time) {
      if (Array.isArray(end_time)) {
        let endYear = Number(end_time[0]);
        let endMonth = Number(end_time[1]);
        end_time = formatRFC3339(new Date(endYear, endMonth, 20));
      }

      return { start_time, end_time };
    } else {
      return { start_time };
    }
  },

  // query comes in as either: {start_time: '2021-03-17T00:00:00Z'} or
  // {start_time: Array(2), end_time: Array(2)}
  // end_time: (2) ['2022', 0]
  // start_time: (2) ['2021', 2]
  queryRecord(store, type, query) {
    let url = `${this.buildURL()}/internal/counters/activity`;
    // check if start and/or end times are in RFC3395 format, if not convert with timezone UTC/zulu.
    let queryParams = this.formatTimeParams(query);
    if (queryParams) {
      return this.ajax(url, 'GET', { data: queryParams }).then((resp) => {
        let response = resp || {};
        response.id = response.request_id || 'no-data';
        return response;
      });
    }
  },
});
