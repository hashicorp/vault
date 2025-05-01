/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class LoginSettingsRoute extends Route {
  @service api;
  @service capabilities;

  async model() {
    // const adapter = this.store.adapterFor('application');

    // const obj = await adapter.ajax('/v1/sys/config/ui/login/default-auth', 'GET', { data: { list: true } });

    return { loginRules: [{ name: 'Root level auth', namespace: 'root/' }] };
  }
}
