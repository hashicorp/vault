/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { tabsForAuthSection } from 'vault/helpers/tabs-for-auth-section';

export default class AccessMethodIndexRoute extends Route {
  @service router;

  beforeModel() {
    const methodModel = this.modelFor('vault.cluster.access.method');
    const paths = methodModel.paths
      ? methodModel.paths.paths.filter((path) => path.navigation === true)
      : null;
    const activeTab = tabsForAuthSection([methodModel, 'authConfig', paths])[0];
    return this.router.transitionTo(activeTab.route, ...activeTab.routeParams);
  }
}
