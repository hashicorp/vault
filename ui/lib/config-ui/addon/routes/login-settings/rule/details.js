/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class LoginSettingsRuleDetailsRoute extends Route {
  @service api;
  @service store;

  async model() {
    const { name } = this.paramsFor('login-settings.rule');

    const adapter = this.store.adapterFor('application');

    const rule = await adapter.ajax(`/v1/sys/config/ui/login/default-auth/${encodeURI(name)}`, 'GET');

    return { rule };
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    const { rule } = resolvedModel;

    controller.breadcrumbs = [{ label: 'UI Login rules', route: 'login-settings' }, { label: rule.name }];
  }
}
