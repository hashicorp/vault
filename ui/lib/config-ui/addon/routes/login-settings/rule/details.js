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
    const adapter = this.store.adapterFor('application');

    const rule = await adapter.ajax(
      `/v1/sys/config/ui/login/default-auth/${encodeURI('Login rule 1')}`,
      'GET'
    );

    return { rule };
  }
}
