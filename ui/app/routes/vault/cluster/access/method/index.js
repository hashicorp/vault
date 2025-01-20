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
    let { methodType, paths } = this.modelFor('vault.cluster.access.method');
    paths = paths ? paths.paths.filter((path) => path.navigation === true) : null;
    const activeTab = tabsForAuthSection([methodType, 'authConfig', paths])[0].routeParams;
    return this.router.transitionTo(...activeTab);
  }
}
