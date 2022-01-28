import Application from '../application';
import { zonedTimeToUtc } from 'date-fns-tz'; // https://github.com/marnusw/date-fns-tz#zonedtimetoutc

export default Application.extend({
  queryRecord(store, type, query) {
    let url = `${this.buildURL()}/internal/counters/activity`;
    // Query has startTime defined. The API will return the endTime if none is provided.
    return this.ajax(url, 'GET', { data: query }).then((resp) => {
      let response = resp || {};
      // if the response is a 204 it has no request id (ARG TODO test that it returns a 204)
      response.id = response.request_id || 'no-data';
      return response;
    });
  },

  utcDate(dateObject) {
    // To remove the timezone of the local user (API returns and expects Zulu time/UTC) we need to use a method provided by date-fns-tz to return the UTC date
    let timeZone = Intl.DateTimeFormat().resolvedOptions().timeZone; // browser API method
    return zonedTimeToUtc(dateObject, timeZone);
  },
  // called from components
  queryClientActivity(start_time, end_time) {
    // do not query without start_time. Otherwise returns last year data, which is not reflective of billing data.
    if (start_time) {
      // start and end time come in as month,year ex: "3,2021" Need to turn into RFC3395 on UTC/Zulu time.
      let startTimeQuery;
      if (start_time.split(',').length > 1) {
        // check if already an RFC33339 timestamp.
        let utcDate = this.utcDate(
          new Date(Number(start_time.split(',')[1]), Number(start_time.split(',')[0] - 1))
        );
        startTimeQuery = utcDate.toISOString();
      }
      let url = `${this.buildURL()}/internal/counters/activity`;
      let queryParams = {};
      if (!end_time) {
        queryParams = { data: { start_time: startTimeQuery } };
      } else {
        let utcDateEnd = this.utcDate(
          new Date(Number(end_time.split(',')[1]), Number(end_time.split(',')[0]))
        );
        // ARG TODO !!!! SUPER IMPORTANT, with endDate you need to make it the last day of the month, right now it's the first!!!
        end_time = utcDateEnd.toISOString();
        start_time = startTimeQuery;
        queryParams = { data: { start_time, end_time } };
      }
      return this.ajax(url, 'GET', queryParams).then((resp) => {
        return resp;
      });
    }
  },
});
