/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';

export default class ErrorRoute extends Route {
  setupController(controller) {
    super.setupController(...arguments);
    const params = this.paramsFor('vault.cluster');
    controller.set('clusterId', params.cluster_name);
    controller.set('ns', params.namespaceQueryParam);
  }
}
