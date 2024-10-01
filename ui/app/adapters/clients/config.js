/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Application from '../application';

export default Application.extend({
  queryRecord() {
    return this.ajax(this.urlForQuery(), 'GET').then((resp) => {
      resp.id = resp.request_id;
      return resp;
    });
  },

  urlForUpdateRecord() {
    return this.buildURL() + '/internal/counters/config';
  },

  urlForQuery() {
    return this.buildURL() + '/internal/counters/config';
  },
});
