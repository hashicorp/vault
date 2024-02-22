/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { tabsForAuthSection } from 'vault/helpers/tabs-for-auth-section';

export default Route.extend({
  beforeModel() {
    const model = this.modelFor('vault.cluster.settings.auth.configure');
    const section = tabsForAuthSection([model])[0].routeParams.lastObject;
    return this.transitionTo('vault.cluster.settings.auth.configure.section', section);
  },
});
