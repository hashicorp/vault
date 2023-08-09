/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import ApplicationAdapter from './application';
const fetchUrl = '/v1/sys/storage/raft/configuration';

export default ApplicationAdapter.extend({
  urlForFindAll() {
    return fetchUrl;
  },
  urlForQuery() {
    return fetchUrl;
  },
  urlForDeleteRecord() {
    return '/v1/sys/storage/raft/remove-peer';
  },
  deleteRecord(store, type, snapshot) {
    const server_id = snapshot.attr('nodeId');
    const url = '/v1/sys/storage/raft/remove-peer';
    return this.ajax(url, 'POST', { data: { server_id } });
  },
});
