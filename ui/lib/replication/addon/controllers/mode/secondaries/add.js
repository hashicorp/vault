/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ReplicationController from 'replication/controllers/application';

export default ReplicationController.extend({
  actions: {
    updateTtl: function (ttl) {
      this.set('ttl', `${ttl.seconds}s`);
    },
  },
});
