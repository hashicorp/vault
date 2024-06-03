/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from '../application';
import { formatDateObject } from 'core/utils/client-count-utils';

export default class ActivityAdapter extends ApplicationAdapter {
  // javascript localizes new Date() objects but all activity log data is stored in UTC
  // create date object from user's input using Date.UTC() then send to backend as unix
  // time params from the backend are formatted as a zulu timestamp
  formatQueryParams(queryParams) {
    if (queryParams?.current_billing_period) {
      // { current_billing_period: true } automatically queries the activity log
      // from the builtin license start timestamp to the current month
      return queryParams;
    }
    let { start_time, end_time } = queryParams;
    start_time = start_time.timestamp || formatDateObject(start_time);
    end_time = end_time.timestamp || formatDateObject(end_time, true);
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
