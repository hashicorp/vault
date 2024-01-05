/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from './application';

export default ApplicationAdapter.extend({
  getStatusUrl(mode) {
    return this.buildURL() + `/replication/${mode}/status`;
  },

  fetchStatus(mode) {
    const url = this.getStatusUrl(mode);
    return this.ajax(url, 'GET', { unauthenticated: true }).then((resp) => {
      return resp.data;
    });
  },
});
