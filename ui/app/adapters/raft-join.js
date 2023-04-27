/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import ApplicationAdapter from './application';

export default ApplicationAdapter.extend({
  urlForCreateRecord() {
    return '/v1/sys/storage/raft/join';
  },
});
