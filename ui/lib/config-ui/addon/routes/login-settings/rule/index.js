/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class LoginSettingsRuleIndexRoute extends Route {
  @service('app-router') router;

  redirect() {
    this.router.transitionTo('vault.cluster.config-ui.login-settings.rule.details');
  }
}
