/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import AuthMethodForm from 'vault/forms/auth/method';

import type { ModelFrom } from 'vault/vault/route';

export type AuthEnableModel = ModelFrom<VaultClusterSettingsAuthEnableRoute>;

export default class VaultClusterSettingsAuthEnableRoute extends Route {
  model() {
    const defaults = {
      config: { listing_visibility: false },
      user_lockout_config: {},
    };
    return new AuthMethodForm(defaults, { isNew: true });
  }
}
