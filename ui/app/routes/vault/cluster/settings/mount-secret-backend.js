/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class VaultClusterSettingsMountSecretBackendRoute extends Route {
  @service store;

  model() {
    const secretEngine = this.store.createRecord('secret-engine');
    secretEngine.set('config', this.store.createRecord('mount-config'));
    return secretEngine;
  }
}
