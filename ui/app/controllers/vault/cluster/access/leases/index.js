/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { service } from '@ember/service';

export default Controller.extend({
  router: service(),
  actions: {
    lookupLease(id) {
      this.router.transitionTo('vault.cluster.access.leases.show', id);
    },
  },
});
