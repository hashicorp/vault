/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Actions from 'core/components/replication-actions-single';

export default Actions.extend({
  tagName: '',

  actions: {
    onSubmit(replicationMode, clusterMode, evt) {
      // No data is submitted for disable request
      return this.onSubmit(replicationMode, clusterMode, null, evt);
    },
  },
});
