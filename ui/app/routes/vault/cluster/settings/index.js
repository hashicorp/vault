/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';

export default Route.extend({
  beforeModel: function (transition) {
    if (transition.targetName === this.routeName) {
      transition.abort();
      return this.replaceWith('vault.cluster.settings.mount-secret-backend');
    }
  },
});
