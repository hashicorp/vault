/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class LoginSettingsRuleDetailsRoute extends Route {
  @service store;

  async model() {
    try {
      const { name } = this.paramsFor('login-settings.rule');
      if (!name) return null;

      const adapter = this.store.adapterFor('application');
      const rule = await adapter.ajax(`/v1/sys/config/ui/login/default-auth/${encodeURI(name)}`, 'GET');

      return { rule: rule.data };
    } catch (e) {
      return null;
    }
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);

    controller.breadcrumbs = [{ label: 'UI login rules', route: 'login-settings' }, { label: '' }];
  }
}
