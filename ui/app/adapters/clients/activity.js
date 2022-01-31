import Application from '../application';
import { zonedTimeToUtc } from 'date-fns-tz'; // https://github.com/marnusw/date-fns-tz#zonedtimetoutc

export default Application.extend({
  checkTimeType(query) {
    let { start_time, end_time } = query;
    // do not query without start_time. Otherwise returns last year data, which is not reflective of billing data.
    if (start_time) {
      // check if start_time is a RFC3339 timestamp and if not convert to one.
      if (start_time.split(',').length > 1) {
        let utcDate = this.utcDate(
          new Date(Number(start_time.split(',')[1]), Number(start_time.split(',')[0] - 1))
        );
        start_time = utcDate.toISOString();
      }
      // look for end_time. If there is one check if it is a RFC3339 timestamp otherwise convert it.
      if (end_time) {
        let utcDateEnd = this.utcDate(
          new Date(Number(end_time.split(',')[1]), Number(end_time.split(',')[0]))
        );
        // ARG TODO !!!! SUPER IMPORTANT, with endDate you need to make it the last day of the month, right now it's the first!!!
        end_time = utcDateEnd.toISOString();
        return { start_time, end_time };
      } else {
        return { start_time };
      }
    } else {
      // did not have a start time, return null through to component.
      return null;
    }
  },

  utcDate(dateObject) {
    // To remove the timezone of the local user (API returns and expects Zulu time/UTC) we need to use a method provided by date-fns-tz to return the UTC date
    let timeZone = Intl.DateTimeFormat().resolvedOptions().timeZone; // browser API method
    return zonedTimeToUtc(dateObject, timeZone);
  },

  // ARG TODO current Month tab is hitting this endpoint. Need to amend so only hit on Monthly history (large payload)
  queryRecord(store, type, query) {
    let url = `${this.buildURL()}/internal/counters/activity`;
    // check if start and/or end times are in RFC3395 format, if not convert with timezone UTC/zulu.
    let queryParams = this.checkTimeType(query);
    if (queryParams) {
      return this.ajax(url, 'GET', { data: queryParams }).then((resp) => {
        let response = resp || {};
        // if the response is a 204 it has no request id (ARG TODO test that it returns a 204)
        response.id = response.request_id || 'no-data';
        return response;
      });
    } else {
      // did not have a start time, return null through to component.
      return null;
    }
  },
});
