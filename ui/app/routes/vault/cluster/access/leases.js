/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import ClusterRoute from 'vault/mixins/cluster-route';
import { inject as service } from '@ember/service';

export default Route.extend(ClusterRoute, {
  store: service(),

  model() {
    return this.store.findRecord('capabilities', 'sys/leases/lookup/');
  },
});
