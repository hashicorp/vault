/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class LoginSettingsRuleDetailsRoute extends Route {
  @service('app-router') router;
  @service store;

  beforeModel() {
    const { name } = this.paramsFor('login-settings.rule');
    if (!name) {
      this.router.transitionTo('vault.cluster.config-ui.login-settings.index');
    }
  }

  async model() {
    try {
      const { name } = this.paramsFor('login-settings.rule');

      const adapter = this.store.adapterFor('application');
      const rule = await adapter.ajax(`/v1/sys/config/ui/login/default-auth/${encodeURI(name)}`, 'GET');

      return { rule: rule.data };
    } catch (error) {
      return { rule: {}, error };
    }
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);

    controller.breadcrumbs = [{ label: 'UI login rules', route: 'login-settings' }, { label: '' }];
  }
}
