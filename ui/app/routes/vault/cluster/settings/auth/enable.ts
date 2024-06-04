/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

import type { ModelFrom } from 'vault/vault/route';
import type Store from '@ember-data/store';

export type AuthEnableModel = ModelFrom<VaultClusterSettingsAuthEnableRoute>;

export default class VaultClusterSettingsAuthEnableRoute extends Route {
  @service declare readonly store: Store;

  model() {
    const authMethod = this.store.createRecord('auth-method');
    authMethod.set('config', this.store.createRecord('mount-config'));
    return authMethod;
  }
}
