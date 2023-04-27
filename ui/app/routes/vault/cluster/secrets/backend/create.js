/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';

export default Route.extend({
  beforeModel() {
    const { secret, initialKey } = this.paramsFor(this.routeName);
    const qp = initialKey || secret;
    return this.transitionTo('vault.cluster.secrets.backend.create-root', {
      queryParams: { initialKey: qp },
    });
  },
});
