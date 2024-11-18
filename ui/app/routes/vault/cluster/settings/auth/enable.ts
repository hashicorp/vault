/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import AuthMethod from 'vault/entities/auth-method';

import type { ModelFrom } from 'vault/vault/route';

export type AuthEnableModel = ModelFrom<VaultClusterSettingsAuthEnableRoute>;

export default class VaultClusterSettingsAuthEnableRoute extends Route {
  model() {
    return new AuthMethod();
  }
}
