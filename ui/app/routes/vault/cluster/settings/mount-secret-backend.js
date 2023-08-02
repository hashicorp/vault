/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class VaultClusterSettingsMountSecretBackendRoute extends Route {
  @service store;

  model() {
    const secretEngine = this.store.createRecord('secret-engine');
    secretEngine.set('config', this.store.createRecord('mount-config'));
    return secretEngine;
  }
}
