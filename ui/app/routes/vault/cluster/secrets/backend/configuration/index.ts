/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import RouterService from '@ember/routing/router-service';
import { service } from '@ember/service';

export default class BackendConfigurationIndexRoute extends Route {
  @service declare readonly router: RouterService;

  beforeModel() {
    return this.router.replaceWith('vault.cluster.secrets.backend.configuration.general-settings');
  }
}
