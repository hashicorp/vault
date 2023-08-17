/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';

export default Route.extend({
  beforeModel(transition) {
    if (
      this.modelFor('vault.cluster.access.leases').get('canList') &&
      transition.targetName === this.routeName
    ) {
      return this.replaceWith('vault.cluster.access.leases.list-root');
    } else {
      return;
    }
  },
});
