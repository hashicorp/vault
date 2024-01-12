/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';

export default Route.extend({
  beforeModel() {
    return this.transitionTo('vault.cluster.dashboard');
  },
});
