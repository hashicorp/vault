/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import ClusterRoute from 'vault/mixins/cluster-route';
import { service } from '@ember/service';

export default Route.extend(ClusterRoute, {
  store: service(),
  version: service(),
  router: service(),

  beforeModel() {
    if (this.version.isCommunity) {
      this.router.transitionTo('vault.cluster');
    }
  },

  model() {
    return this.store.queryRecord('license', {});
  },
});
