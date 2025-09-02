/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { tabsForAuthSection } from 'vault/helpers/tabs-for-auth-section';

export default class SettingsAuthConfigureRoute extends Route {
  @service router;

  beforeModel() {
    const model = this.modelFor('vault.cluster.settings.auth.configure');
    const section = tabsForAuthSection([model])[0].routeParams.slice().pop();
    return this.router.transitionTo('vault.cluster.settings.auth.configure.section', section);
  }
}
