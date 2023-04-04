/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';

export default Route.extend({
  beforeModel() {
    return this.transitionTo('vault.cluster.policies', 'acl');
  },
});
