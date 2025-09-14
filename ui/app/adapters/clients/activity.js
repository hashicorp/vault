/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import queryParamString from 'vault/utils/query-param-string';
import ApplicationAdapter from '../application';
import { debug } from '@ember/debug';
import { parseJSON, isValid } from 'date-fns';

export default class ActivityAdapter extends ApplicationAdapter {
  formatQueryParams({ start_time, end_time }) {
    const query = {};

    if (start_time && isValid(parseJSON(start_time))) {
      query.start_time = start_time;
    }
    if (end_time && isValid(parseJSON(end_time))) {
      query.end_time = end_time;
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
    let errorMsg, httpStatus;
    try {
      const options = query?.namespace ? { namespace: query.namespace } : {};
      const resp = await this.rawRequest(url, 'GET', options);
      if (resp.status === 200) {
        return resp.blob();
      }
      // If it's an empty response (eg 204), there's no data so return an error
      errorMsg = 'No data to export in provided time range.';
      httpStatus = resp.status;
    } catch (e) {
      const { errors } = await e.json();
      errorMsg = errors?.join('. ');
      httpStatus = e.status;
    }
    // counters/activity/export returns a ReadableStream so we manually handle errors here
    // hopefully this can be improved when this file is migrated to use the api service.
    if (errorMsg) {
      const error = new Error(errorMsg);
      error.httpStatus = httpStatus;
      throw error;
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
