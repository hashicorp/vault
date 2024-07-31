/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from '../application';
import { formatDateObject } from 'core/utils/client-count-utils';
import { debug } from '@ember/debug';

export default class ActivityAdapter extends ApplicationAdapter {
  formatTimeParam(dateObj, isEnd = false) {
    let formatted;
    if (dateObj) {
      try {
        const iso = dateObj.timestamp || formatDateObject(dateObj, isEnd);
        formatted = iso;
      } catch (e) {
        // carry on
      }
    }
    return formatted;
  }
  // javascript localizes new Date() objects but all activity log data is stored in UTC
  // create date object from user's input using Date.UTC() then send to backend as unix
  // time params from the backend are formatted as a zulu timestamp
  formatQueryParams(queryParams) {
    const query = {};
    const start = this.formatTimeParam(queryParams?.start_time);
    const end = this.formatTimeParam(queryParams?.end_time, true);
    if (start) {
      query.start_time = start;
    }
    if (end) {
      query.end_time = end;
    }
    return query;
  }

  queryRecord(store, type, query) {
    const url = `${this.buildURL()}/internal/counters/activity`;
    const queryParams = this.formatQueryParams(query);
    return this.ajax(url, 'GET', { data: queryParams }).then((resp) => {
      const response = resp || {};
      response.id = response.request_id || 'no-data';
      return response;
    });
  }

  urlForFindRecord(id) {
    // debug reminder so model is stored in Ember data with the same id for consistency
    if (id !== 'clients/activity') {
      debug(`findRecord('clients/activity') should pass 'clients/activity' as the id, you passed: '${id}'`);
    }
    return `${this.buildURL()}/internal/counters/activity`;
  }
}
