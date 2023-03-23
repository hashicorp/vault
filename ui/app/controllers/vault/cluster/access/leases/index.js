/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Controller from '@ember/controller';

export default Controller.extend({
  actions: {
    lookupLease(id) {
      this.transitionToRoute('vault.cluster.access.leases.show', id);
    },
  },
});
