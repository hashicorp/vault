/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import ApplicationAdapter from '../application';
import { getUnixTime } from 'date-fns';

export default class ActivityAdapter extends ApplicationAdapter {
  // javascript localizes new Date() objects but all activity log data is stored in UTC
  // create date object from user's input using Date.UTC() then send to backend as unix
  // time params from the backend are formatted as a zulu timestamp
  formatQueryParams(queryParams) {
    let { start_time, end_time } = queryParams;
    start_time = start_time.timestamp || getUnixTime(Date.UTC(start_time.year, start_time.monthIdx, 1));
    // day=0 for Date.UTC() returns the last day of the month before
    // increase monthIdx by one to get last day of queried month
    end_time = end_time.timestamp || getUnixTime(Date.UTC(end_time.year, end_time.monthIdx + 1, 0));
    return { start_time, end_time };
  }

  queryRecord(store, type, query) {
    const url = `${this.buildURL()}/internal/counters/activity`;
    const queryParams = this.formatQueryParams(query);
    if (queryParams) {
      return this.ajax(url, 'GET', { data: queryParams }).then((resp) => {
        const response = resp || {};
        response.id = response.request_id || 'no-data';
        return response;
      });
    }
  }
}
