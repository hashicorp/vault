/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

import type { ModelFrom } from 'vault/vault/route';
import type Store from '@ember-data/store';

export type MountSecretBackendModel = ModelFrom<VaultClusterSettingsMountSecretBackendRoute>;

export default class VaultClusterSettingsMountSecretBackendRoute extends Route {
  @service declare readonly store: Store;

  model() {
    const secretEngine = this.store.createRecord('secret-engine');
    secretEngine.set('config', this.store.createRecord('mount-config'));
    return secretEngine;
  }
}
