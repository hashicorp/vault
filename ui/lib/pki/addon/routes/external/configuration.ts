/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

import type RouterService from '@ember/routing/router-service';

export default class PkiExternalConfigurationRoute extends Route {
  @service('app-router') declare readonly router: RouterService;

  // There is not a dedicated "Configuration" route for this engine. Config details
  // are displayed within the dns-providers and acme-accounts routes instead.
  // This route (and redirect) exists for two reasons:
  // 1. ManageDropdown's @configRoute passes the value to an HDS <D.Interactive @route>,
  //    which only resolves engine-relative routes and currently does not support links to external routes.
  // 2. If a dedicated external config page is added in the future, this route is prepped and ready!
  redirect() {
    this.router.transitionTo('vault.cluster.secrets.backend.configuration.general-settings');
  }
}
