/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import ClusterRoute from 'vault/mixins/cluster-route';
import { service } from '@ember/service';

export default Route.extend(ClusterRoute, {
  store: service(),

  model() {
    return this.store.findRecord('capabilities', 'sys/leases/lookup/');
  },
});
