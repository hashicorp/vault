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
    // findAll method will return all records in store as well as response from server
    // when removing a peer via the cli, stale records would continue to appear until refresh
    // query method will only return records from response
    return this.store.query('server', {});
  },

  actions: {
    doRefresh() {
      this.refresh();
    },
  },
});
