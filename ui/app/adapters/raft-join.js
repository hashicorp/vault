/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from './application';

export default ApplicationAdapter.extend({
  urlForCreateRecord() {
    return '/v1/sys/storage/raft/join';
  },
});
