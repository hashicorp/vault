import Application from '../application';
import { formatRFC3339 } from 'date-fns';

export default Application.extend({
  formatTimeParams(query) {
    let { start_time, end_time } = query;
    // do not query without start_time. Otherwise returns last year data, which is not reflective of billing data.
    if (start_time) {
      if (start_time.split(',').length > 1) {
        let startYear = Number(start_time.split(',')[0]);
        let startMonth = Number(start_time.split(',')[1]);
        start_time = formatRFC3339(new Date(startYear, startMonth));
      }
      // look for end_time
      if (end_time) {
        if (end_time.split(',').length > 1) {
          let endYear = Number(end_time.split(',')[0]);
          let endMonth = Number(end_time.split(',')[1]);
          end_time = formatRFC3339(new Date(endYear, endMonth));
        }
        return { start_time, end_time };
      } else {
        return { start_time };
      }
    } else {
      // did not have a start time, return null through to component.
      return null;
    }
  },

  // ARG TODO current Month tab is hitting this endpoint. Need to amend so only hit on Monthly history (large payload)
  queryRecord(store, type, query) {
    let url = `${this.buildURL()}/internal/counters/activity`;
    // check if start and/or end times are in RFC3395 format, if not convert with timezone UTC/zulu.
    let queryParams = this.formatTimeParams(query);
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
