/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import queryParamString from 'vault/utils/query-param-string';
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
    const options = {
      data: this.formatQueryParams(query),
    };

    if (query?.namespace) {
      options.namespace = query.namespace;
    }

    return this.ajax(url, 'GET', options).then((resp) => {
      const response = resp || {};
      response.id = response.request_id || 'no-data';
      return response;
    });
  }

  async exportData(query) {
    const url = `${this.buildURL()}/internal/counters/activity/export${queryParamString({
      format: query?.format || 'csv',
      start_time: query?.start_time ?? undefined,
      end_time: query?.end_time ?? undefined,
    })}`;
    let errorMsg;
    try {
      const options = query?.namespace ? { namespace: query.namespace } : {};
      const resp = await this.rawRequest(url, 'GET', options);
      if (resp.status === 200) {
        return resp.blob();
      }
      // If it's an empty response (eg 204), there's no data so return an error
      errorMsg = 'no data to export in provided time range.';
    } catch (e) {
      const { errors } = await e.json();
      errorMsg = errors?.join('. ');
    }
    if (errorMsg) {
      throw new Error(errorMsg);
    }
  }

  // Only dashboard uses findRecord, the client count dashboard uses queryRecord
  findRecord(store, type, id) {
    if (id !== 'clients/activity') {
      debug(`findRecord('clients/activity') should pass 'clients/activity' as the id, you passed: '${id}'`);
    }
    const url = `${this.buildURL()}/internal/counters/activity`;
    return this.ajax(url, 'GET', { skipWarnings: true });
  }
}
