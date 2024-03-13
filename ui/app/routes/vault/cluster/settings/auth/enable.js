/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class VaultClusterSettingsAuthEnableRoute extends Route {
  @service store;

  model() {
    const authMethod = this.store.createRecord('auth-method');
    authMethod.set('config', this.store.createRecord('mount-config'));
    return authMethod;
  }
}
