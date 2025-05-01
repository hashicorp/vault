/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class LoginSettingsRoute extends Route {
  @service api;
  @service capabilities;
  @service store;

  async model() {
    const adapter = this.store.adapterFor('application');

    const { data } = await adapter.ajax('/v1/sys/config/ui/login/default-auth', 'GET', {
      data: { list: true },
    });

    // this makes sense with data structure atm, but to be revisited
    const loginRules = [];
    data.keys.forEach((rule) => {
      loginRules.push({ name: rule, namespace: data.key_info[rule].namespace });
    });

    return { loginRules };
  }
}
