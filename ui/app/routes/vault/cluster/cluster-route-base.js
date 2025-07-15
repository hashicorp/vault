/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

// this is the base route for
// all of the CLUSTER_ROUTES that are states before you can use vault
//
import Route from '@ember/routing/route';
import ClusterRoute from 'vault/mixins/cluster-route';

/**
 * @type Class
 */
export default Route.extend(ClusterRoute, {
  model() {
    return this.modelFor('vault.cluster');
  },

  resetController(controller) {
    controller.reset && controller.reset();
  },
});
