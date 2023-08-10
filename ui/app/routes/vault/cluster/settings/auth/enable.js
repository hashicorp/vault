/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class VaultClusterSettingsAuthEnableRoute extends Route {
  @service store;

  model() {
    const authMethod = this.store.createRecord('auth-method');
    authMethod.set('config', this.store.createRecord('mount-config'));
    return authMethod;
  }
}
