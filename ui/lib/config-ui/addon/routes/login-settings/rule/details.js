/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class LoginSettingsRuleDetailsRoute extends Route {
  @service('app-router') router;
  @service api;

  beforeModel() {
    const { name } = this.paramsFor('login-settings.rule');
    if (!name) {
      this.router.transitionTo('vault.cluster.config-ui.login-settings.index');
    }
  }

  async model() {
    const { name } = this.paramsFor('login-settings.rule');

    const rule = await this.api.sys.uiLoginDefaultAuthReadConfiguration(name);

    return { rule: { name, ...rule.data } };
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);

    controller.breadcrumbs = [
      { label: 'Vault', route: 'vault', icon: 'vault', linkExternal: true },
      { label: 'UI login settings', route: 'login-settings' },
      { label: resolvedModel.rule.name },
    ];
  }
}
