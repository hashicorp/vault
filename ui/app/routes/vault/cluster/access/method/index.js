/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { tabsForAuthSection } from 'vault/helpers/tabs-for-auth-section';
export default Route.extend({
  beforeModel() {
    let { methodType, paths } = this.modelFor('vault.cluster.access.method');
    paths = paths ? paths.paths.filter((path) => path.navigation === true) : null;
    const activeTab = tabsForAuthSection([methodType, 'authConfig', paths])[0].routeParams;
    return this.transitionTo(...activeTab);
  },
});
