/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import ClusterRoute from 'vault/mixins/cluster-route';
import { inject as service } from '@ember/service';

export default Route.extend(ClusterRoute, {
  store: service(),
  version: service(),

  beforeModel() {
    if (this.version.isOSS) {
      this.transitionTo('vault.cluster');
    }
  },

  model() {
    return this.store.queryRecord('license', {});
  },
});
