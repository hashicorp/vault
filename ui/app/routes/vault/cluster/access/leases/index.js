/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
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
