/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';

export default Route.extend({
  setupController(controller) {
    this._super(...arguments);
    const targetRoute = location.pathname || '';
    controller.set('isCallback', targetRoute.includes('oidc/callback'));
  },
});
